package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An already-cancelled request must not reach auth or HTTP, even if a transport
// could return a response immediately (as the httpmock transport does).
func TestRequestCancellationBeforeExecution(t *testing.T) {
	operations := []struct {
		name string
		run  func(*RequestBuilder) error
	}{
		{"Get", func(b *RequestBuilder) error { _, err := b.Get("/v1/test"); return err }},
		{"Post", func(b *RequestBuilder) error { _, err := b.Post("/v1/test"); return err }},
		{"Put", func(b *RequestBuilder) error { _, err := b.Put("/v1/test"); return err }},
		{"Patch", func(b *RequestBuilder) error { _, err := b.Patch("/v1/test"); return err }},
		{"Delete", func(b *RequestBuilder) error { _, err := b.Delete("/v1/test"); return err }},
		{"GetBytes", func(b *RequestBuilder) error { _, _, err := b.GetBytes("/v1/test"); return err }},
		{"GetPaginated", func(b *RequestBuilder) error {
			_, err := b.GetPaginated("/v1/test", func([]byte) error { t.Error("cancelled request must not merge a page"); return nil })
			return err
		}},
	}
	for _, operation := range operations {
		for _, expired := range []bool{false, true} {
			name := "Cancelled"
			if expired {
				name = "Expired"
			}
			t.Run(operation.name+"/"+name, func(t *testing.T) {
				tr := setupTestTransport(t)
				httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test", httpmock.NewJsonResponderOrPanic(200, map[string]any{"data": []any{}}))
				ctx, cancel := context.WithCancel(context.Background())
				wantErr := context.Canceled
				if expired {
					cancel()
					ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
					wantErr = context.DeadlineExceeded
				} else {
					cancel()
				}
				defer cancel()
				err := operation.run(tr.NewRequest(ctx))
				assert.ErrorIs(t, err, wantErr)
				assert.Zero(t, httpmock.GetTotalCallCount(), "cancelled request reached HTTP transport")
			})
		}
	}
}

func TestRequestCancellationBetweenPages(t *testing.T) {
	tr := setupTestTransport(t)
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test", func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("cursor") == "next" {
			return httpmock.NewJsonResponse(200, map[string]any{"data": []any{}})
		}
		return httpmock.NewJsonResponse(200, map[string]any{"data": []any{}, "links": map[string]string{"next": "https://api-business.apple.com/v1/test?cursor=next"}})
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	merged := 0
	_, err := tr.NewRequest(ctx).GetPaginated("/v1/test", func([]byte) error { merged++; cancel(); return nil })
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, merged)
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "pagination fetched another page after cancellation")
}
