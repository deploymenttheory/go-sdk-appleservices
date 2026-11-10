package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const maxNotificationBytes = 1 << 20

type notification struct {
	Data struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
}

// receiver is acceptance infrastructure, never a public SDK server.
type receiver struct {
	secret  []byte
	mu      sync.Mutex
	events  []notification
	changed chan struct{}
}

func newReceiver(secret string) *receiver {
	return &receiver{secret: []byte(secret), changed: make(chan struct{})}
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/healthz" && req.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}
	if req.URL.Path != "/apple" {
		http.NotFound(w, req)
		return
	}
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, req.Body, maxNotificationBytes))
	if err != nil {
		var limitErr *http.MaxBytesError
		if errors.As(err, &limitErr) {
			http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid body", http.StatusBadRequest)
		}
		return
	}
	signature, ok := strings.CutPrefix(req.Header.Get("x-apple-signature"), "hmacsha256=")
	got, err := hex.DecodeString(signature)
	mac := hmac.New(sha256.New, r.secret)
	_, _ = mac.Write(body)
	if !ok || err != nil || !hmac.Equal(got, mac.Sum(nil)) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	var event notification
	if json.Unmarshal(body, &event) != nil || event.Data.ID == "" || event.Data.Type == "" {
		http.Error(w, "invalid event", http.StatusBadRequest)
		return
	}
	r.mu.Lock()
	r.events = append(r.events, event)
	close(r.changed)
	r.changed = make(chan struct{})
	r.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func (r *receiver) count() int { r.mu.Lock(); defer r.mu.Unlock(); return len(r.events) }

func (r *receiver) wait(ctx context.Context, after int, matches func(notification) bool) (notification, error) {
	for {
		if err := ctx.Err(); err != nil {
			return notification{}, err
		}
		r.mu.Lock()
		for _, event := range r.events[after:] {
			if matches(event) {
				r.mu.Unlock()
				return event, nil
			}
		}
		changed := r.changed
		r.mu.Unlock()
		select {
		case <-ctx.Done():
			return notification{}, ctx.Err()
		case <-changed:
		}
	}
}

func startReceiver(t *testing.T, r *receiver) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := &http.Server{Handler: r, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 10 * time.Second}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	t.Cleanup(func() {
		http.DefaultClient.CloseIdleConnections()
		server.SetKeepAlivesEnabled(false)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
			t.Errorf("receiver shutdown: %v", err)
		}
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("receiver stopped: %v", err)
		}
	})
	return "http://" + listener.Addr().String()
}

func signedRequest(t *testing.T, address, secret, body string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, address+"/apple", strings.NewReader(body))
	require.NoError(t, err)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	req.Header.Set("x-apple-signature", "hmacsha256="+hex.EncodeToString(mac.Sum(nil)))
	return req
}

func TestReceiverAuthentication(t *testing.T) {
	const body = `{"data":{"id":"ping-1","type":"webhookPingCreated"}}`
	for _, tc := range []struct {
		name, secret, header, body, method string
		status                             int
	}{
		{name: "valid", secret: "secret", body: body, status: 200},
		{name: "wrong secret", secret: "wrong", body: body, status: 401},
		{name: "missing", header: "missing", body: body, status: 401},
		{name: "malformed", header: "hmacsha256=invalid", body: body, status: 401},
		{name: "tampered", secret: "secret", body: body + " ", status: 401},
		{name: "invalid JSON", secret: "secret", body: "invalid", status: 400},
		{name: "method", method: "GET", status: 405},
		{name: "oversize", body: strings.Repeat("x", maxNotificationBytes+1), status: 413},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newReceiver("secret")
			addr := startReceiver(t, r)
			req := signedRequest(t, addr, tc.secret, tc.body)
			if tc.name == "tampered" {
				original := signedRequest(t, addr, tc.secret, body)
				req.Header = original.Header
			}
			if tc.header == "missing" {
				req.Header.Del("x-apple-signature")
			} else if tc.header != "" {
				req.Header.Set("x-apple-signature", tc.header)
			}
			if tc.method != "" {
				req.Method = tc.method
			}
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tc.status, resp.StatusCode)
			if tc.status == 200 {
				assert.Equal(t, 1, r.count())
			} else {
				assert.Zero(t, r.count())
			}
		})
	}
}

func TestReceiverConcurrentReceiptAndWait(t *testing.T) {
	r := newReceiver("secret")
	addr := startReceiver(t, r)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := r.wait(ctx, 0, func(n notification) bool { return n.Data.ID == "ping-1" })
		result <- err
	}()
	var wg sync.WaitGroup
	for range 10 {
		req := signedRequest(t, addr, "secret", `{"data":{"id":"ping-1","type":"webhookPingCreated"}}`)
		wg.Go(func() {
			resp, err := http.DefaultClient.Do(req)
			if assert.NoError(t, err) {
				_ = resp.Body.Close()
				assert.Equal(t, 200, resp.StatusCode)
			}
		})
	}
	wg.Wait()
	require.NoError(t, <-result)
	assert.Equal(t, 10, r.count())
	cancel()
	_, err := r.wait(ctx, 10, func(notification) bool { return true })
	assert.ErrorIs(t, err, context.Canceled)
	deadline, stop := context.WithTimeout(context.Background(), time.Millisecond)
	defer stop()
	_, err = r.wait(deadline, 10, func(notification) bool { return true })
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
