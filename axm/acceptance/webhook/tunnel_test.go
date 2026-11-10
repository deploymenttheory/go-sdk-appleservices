package webhook

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

func startTunnel(t *testing.T, ctx context.Context, origin string) string {
	t.Helper()
	binary := os.Getenv("CLOUDFLARED_BIN")
	require.NotEmpty(t, binary, "use ./axm/acceptance/webhook/run.sh to prepare cloudflared")
	tunnelCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(tunnelCtx, binary, "tunnel", "--no-autoupdate", "--protocol", "http2", "--edge-ip-version", "4", "--url", origin)
	logs, err := cmd.StderrPipe()
	require.NoError(t, err)
	cmd.Stdout = io.Discard
	require.NoError(t, cmd.Start())
	done := make(chan struct{})
	urls := make(chan string, 1)
	registered := make(chan struct{})
	scanned := make(chan struct{})
	var logMu sync.Mutex
	var problems []string
	go func() {
		defer close(scanned)
		connected := false
		pattern := regexp.MustCompile(`https://[a-z0-9-]+\.trycloudflare\.com`)
		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			line := scanner.Text()
			if !connected && strings.Contains(line, "Registered tunnel connection") {
				close(registered)
				connected = true
			}
			if strings.Contains(line, "ERR") || strings.Contains(line, "WRN") {
				logMu.Lock()
				problems = append(problems, line)
				if len(problems) > 5 {
					problems = problems[1:]
				}
				logMu.Unlock()
			}
			if address := pattern.FindString(scanner.Text()); address != "" {
				select {
				case urls <- address:
				default:
				}
			}
		}
		if err := scanner.Err(); err != nil {
			logMu.Lock()
			problems = append(problems, "tunnel log reader: "+err.Error())
			logMu.Unlock()
		}
	}()
	go func() { <-scanned; _ = cmd.Wait(); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	startup, cancelStartup := context.WithTimeout(ctx, 60*time.Second)
	defer cancelStartup()
	address, err := waitForTunnel(startup, urls, registered, done)
	if err != nil {
		logMu.Lock()
		t.Log(strings.Join(problems, "\n"))
		logMu.Unlock()
		t.Fatalf("tunnel startup: %v", err)
	}
	probe := &http.Client{Timeout: 5 * time.Second}
	defer probe.CloseIdleConnections()
	// DNS can propagate after connection registration. Check the provider first
	// so an early system lookup cannot cache NXDOMAIN for the test's lifetime.
	hostname := strings.TrimPrefix(address, "https://")
	err = waitForTunnelDNS(startup, probe, hostname)
	require.NoError(t, err, "Cloudflare did not publish the temporary tunnel DNS record")
	lastProbe := "no response"
	err = poll(startup, func() (bool, error) {
		req, err := http.NewRequestWithContext(startup, http.MethodGet, address+"/healthz", nil)
		if err != nil {
			return false, err
		}
		resp, err := probe.Do(req)
		if err != nil {
			lastProbe = err.Error()
			return false, nil
		}
		defer resp.Body.Close()
		lastProbe = resp.Status
		return resp.StatusCode == 200, nil
	})
	if err != nil {
		t.Logf("Last HTTPS readiness probe: %s", lastProbe)
		logMu.Lock()
		t.Log(strings.Join(problems, "\n"))
		logMu.Unlock()
	}
	require.NoError(t, err, "temporary HTTPS tunnel did not become ready")
	return address
}

func waitForTunnelDNS(ctx context.Context, client *http.Client, hostname string) error {
	// Registration and DNS publication are separate operations. Even a public
	// recursive resolver can cache an NXDOMAIN requested immediately at registration.
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	return poll(ctx, func() (bool, error) { return tunnelDNSReady(ctx, client, hostname) })
}

func TestTunnelDNSPropagationGrace(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var queried atomic.Bool
		client := &http.Client{Transport: tunnelRoundTripper(func(*http.Request) (*http.Response, error) {
			queried.Store(true)
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"Status":0,"Answer":[{"name":"example.trycloudflare.com","type":1,"data":"104.16.230.132"}]}`)), Header: make(http.Header)}, nil
		})}
		done := make(chan error, 1)
		go func() { done <- waitForTunnelDNS(t.Context(), client, "example.trycloudflare.com") }()
		synctest.Wait()
		time.Sleep(4 * time.Second)
		synctest.Wait()
		require.False(t, queried.Load(), "DNS must not be queried before the propagation grace period")
		time.Sleep(time.Second)
		synctest.Wait()
		require.NoError(t, <-done)
		require.True(t, queried.Load())
	})
}

func tunnelDNSReady(ctx context.Context, client *http.Client, hostname string) (bool, error) {
	query := url.Values{"name": {hostname}, "type": {"A"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://cloudflare-dns.com/dns-query?"+query.Encode(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", "application/dns-json")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	var result struct {
		Status *int
		Answer []struct {
			Name string
			Type int
			Data string
		}
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&result); err != nil {
		return false, fmt.Errorf("decode Cloudflare DNS response: %w", err)
	}
	if result.Status == nil || *result.Status != 0 {
		return false, nil
	}
	for _, answer := range result.Answer {
		if strings.EqualFold(strings.TrimSuffix(answer.Name, "."), hostname) && answer.Type == 1 && net.ParseIP(answer.Data).To4() != nil {
			return true, nil
		}
	}
	return false, nil
}

type tunnelRoundTripper func(*http.Request) (*http.Response, error)

func (f tunnelRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestTunnelDNSPublication(t *testing.T) {
	for _, tc := range []struct {
		name, body     string
		status         int
		ready, invalid bool
	}{
		{name: "published", body: `{"Status":0,"Answer":[{"name":"example.trycloudflare.com.","type":1,"data":"104.16.230.132"}]}`, status: 200, ready: true},
		{name: "not published", body: `{"Status":3}`, status: 200},
		{name: "missing record", body: `{"Status":0}`, status: 200},
		{name: "unrelated record", body: `{"Status":0,"Answer":[{"name":"other.trycloudflare.com","type":1,"data":"104.16.230.132"}]}`, status: 200},
		{name: "malformed address", body: `{"Status":0,"Answer":[{"name":"example.trycloudflare.com","type":1,"data":"invalid"}]}`, status: 200},
		{name: "unavailable", status: 503},
		{name: "invalid JSON", body: `invalid`, status: 200, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: tunnelRoundTripper(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, "application/dns-json", req.Header.Get("Accept"))
				require.Equal(t, "example.trycloudflare.com", req.URL.Query().Get("name"))
				require.Equal(t, "A", req.URL.Query().Get("type"))
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			ready, err := tunnelDNSReady(t.Context(), client, "example.trycloudflare.com")
			if tc.invalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.ready, ready)
		})
	}
}

// Cloudflare publishes the URL before registering its DNS record. Probing at
// publication can cache NXDOMAIN in the system resolver for the whole test run.
func waitForTunnel(ctx context.Context, urls <-chan string, registered <-chan struct{}, exited <-chan struct{}) (string, error) {
	var address string
	connected := false
	for address == "" || !connected {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		select {
		case address = <-urls:
		case <-registered:
			connected = true
			registered = nil
		case <-exited:
			return "", fmt.Errorf("cloudflared exited before tunnel registration")
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return address, nil
}

func TestTunnelWaitFailureAndOrdering(t *testing.T) {
	t.Run("registration precedes URL", func(t *testing.T) {
		registered := make(chan struct{})
		close(registered)
		urls := make(chan string, 1)
		urls <- "https://example.trycloudflare.com"
		address, err := waitForTunnel(t.Context(), urls, registered, make(chan struct{}))
		require.NoError(t, err)
		require.Equal(t, "https://example.trycloudflare.com", address)
	})
	t.Run("early process exit", func(t *testing.T) {
		exited := make(chan struct{})
		close(exited)
		_, err := waitForTunnel(t.Context(), make(chan string), make(chan struct{}), exited)
		require.ErrorContains(t, err, "exited before tunnel registration")
	})
	t.Run("cancelled startup", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		_, err := waitForTunnel(ctx, make(chan string), make(chan struct{}), make(chan struct{}))
		require.ErrorIs(t, err, context.Canceled)
	})
}

func TestTunnelWaitsForRegistration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		urls := make(chan string, 1)
		registered := make(chan struct{})
		exited := make(chan struct{})
		urls <- "https://example.trycloudflare.com"
		result := make(chan string, 1)
		go func() {
			address, err := waitForTunnel(t.Context(), urls, registered, exited)
			if err != nil {
				result <- err.Error()
				return
			}
			result <- address
		}()
		synctest.Wait()
		select {
		case <-result:
			t.Fatal("URL publication must not permit DNS probing before tunnel registration")
		default:
		}
		close(registered)
		require.Equal(t, "https://example.trycloudflare.com", <-result)
	})
}

// TestTunnelHTTPS checks infrastructure without Apple credentials or API mutations.
func TestTunnelHTTPS(t *testing.T) {
	if os.Getenv("AXM_WEBHOOK_TUNNEL_CHECK") != "true" {
		t.Skip("set AXM_WEBHOOK_TUNNEL_CHECK=true to check temporary HTTPS connectivity")
	}
	ctx, cancel := context.WithTimeout(liveContext, 90*time.Second)
	defer cancel()
	receiver := newReceiver("connectivity-check")
	origin := startReceiver(t, receiver)
	address := startTunnel(t, ctx, origin)
	req := signedRequest(t, address, "connectivity-check", `{"data":{"id":"synthetic-ping","type":"webhookPingCreated"}}`).WithContext(ctx)
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_, err = receiver.wait(ctx, 0, func(n notification) bool { return n.Data.ID == "synthetic-ping" })
	require.NoError(t, err)
	t.Log("Temporary HTTPS receiver accepted a signed synthetic notification")
}

type rateLimitDelay struct {
	delay time.Duration
	cause error
}

func (r *rateLimitDelay) Error() string { return r.cause.Error() }
func (r *rateLimitDelay) Unwrap() error { return r.cause }

func deliveryPollingError(resp *resty.Response, err error) error {
	if resp == nil || resp.StatusCode() != http.StatusTooManyRequests || err == nil {
		return err
	}
	delay := time.Minute
	value := resp.Header().Get("Retry-After")
	if seconds, e := strconv.Atoi(value); e == nil && seconds >= 0 {
		delay = max(time.Duration(seconds)*time.Second, 2*time.Second)
	} else if date, e := http.ParseTime(value); e == nil {
		delay = max(time.Until(date), 2*time.Second)
	}
	return &rateLimitDelay{delay: delay, cause: err}
}

func TestDeliveryPollingRetryAfter(t *testing.T) {
	for _, tc := range []struct {
		name, header string
		want         time.Duration
	}{
		{"seconds", "60", time.Minute},
		{"missing", "", time.Minute},
		{"invalid", "invalid", time.Minute},
		{"minimum delay", "0", 2 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &resty.Response{RawResponse: &http.Response{StatusCode: 429, Header: make(http.Header)}}
			resp.Header().Set("Retry-After", tc.header)
			cause := errors.New("limited")
			err := deliveryPollingError(resp, cause)
			var limited *rateLimitDelay
			require.ErrorAs(t, err, &limited)
			require.Equal(t, tc.want, limited.delay)
			require.ErrorIs(t, err, cause)
		})
	}
	t.Run("HTTP date", func(t *testing.T) {
		deadline := time.Now().UTC().Add(time.Minute).Truncate(time.Second)
		resp := &resty.Response{RawResponse: &http.Response{StatusCode: 429, Header: make(http.Header)}}
		resp.Header().Set("Retry-After", deadline.Format(http.TimeFormat))
		var limited *rateLimitDelay
		require.ErrorAs(t, deliveryPollingError(resp, errors.New("limited")), &limited)
		require.InDelta(t, time.Until(deadline).Seconds(), limited.delay.Seconds(), 1)
	})
	t.Run("other failures are not retried", func(t *testing.T) {
		cause := errors.New("forbidden")
		resp := &resty.Response{RawResponse: &http.Response{StatusCode: 403, Header: make(http.Header)}}
		require.Same(t, cause, deliveryPollingError(resp, cause))
	})
}

func poll(ctx context.Context, fn func() (bool, error)) error {
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		ok, err := fn()
		delay := 2 * time.Second
		if err != nil {
			var rateLimit *rateLimitDelay
			if !errors.As(err, &rateLimit) {
				return err
			}
			delay = rateLimit.delay
		}
		if ok && err == nil {
			return nil
		}
		timer.Reset(delay)
		select {
		case <-ctx.Done():
			return fmt.Errorf("poll deadline: %w", ctx.Err())
		case <-timer.C:
		}
	}
}

func TestPollHonorsRateLimitDelay(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var calls atomic.Int64
		done := make(chan error, 1)
		go func() {
			done <- poll(t.Context(), func() (bool, error) {
				if calls.Add(1) == 1 {
					return false, &rateLimitDelay{delay: time.Minute, cause: errors.New("limited")}
				}
				return true, nil
			})
		}()
		synctest.Wait()
		time.Sleep(59 * time.Second)
		synctest.Wait()
		require.Equal(t, int64(1), calls.Load(), "must not poll while Retry-After is active")
		time.Sleep(time.Second)
		synctest.Wait()
		require.NoError(t, <-done)
		require.Equal(t, int64(2), calls.Load())
	})
}

func TestPollDeadlineDuringRateLimit(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()
		calls := 0
		err := poll(ctx, func() (bool, error) {
			calls++
			return false, &rateLimitDelay{delay: time.Minute, cause: errors.New("limited")}
		})
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Equal(t, 1, calls)
	})
}
