package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func setupTestTransport(t *testing.T) *Transport {
	mockAuth := &testAuthProvider{}

	transport := &Transport{
		httpClient:   resty.New(),
		logger:       zap.NewNop(),
		auth:         mockAuth,
		errorHandler: NewErrorHandler(zap.NewNop()),
		baseURL:      "https://api-business.apple.com",
	}

	transport.httpClient.SetBaseURL(transport.baseURL)

	// Setup auth middleware
	transport.httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		return transport.auth.ApplyAuth(req)
	})

	httpmock.ActivateNonDefault(transport.httpClient.Client())

	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return transport
}

// testAuthProvider is a local auth provider for request tests to avoid collision with transport_test.go's MockAuthProvider.
type testAuthProvider struct{}

func (m *testAuthProvider) ApplyAuth(req *resty.Request) error {
	req.SetHeader("Authorization", "Bearer mock-token")
	return nil
}

func TestTransport_Get_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background())
	_, err := transport.execute(req, "GET", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute GET failed: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("result['status'] = %v, want 'ok'", result["status"])
	}
}

func TestTransport_Get_WithQueryParams(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/test\?`,
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetQueryParam("limit", "10").
		SetQueryParam("cursor", "abc")
	_, err := transport.execute(req, "GET", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute GET with query params failed: %v", err)
	}
}

func TestTransport_Get_WithHeaders(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetHeader("X-Custom-Header", "custom-value").
		SetHeader("Accept", "application/json")
	_, err := transport.execute(req, "GET", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute GET with headers failed: %v", err)
	}
}

func TestTransport_Post_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(201, map[string]string{"id": "12345"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetBody(map[string]string{"name": "test"})
	_, err := transport.execute(req, "POST", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute POST failed: %v", err)
	}

	if result["id"] != "12345" {
		t.Errorf("result['id'] = %v, want '12345'", result["id"])
	}
}

func TestTransport_Post_NilBody(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background())
	_, err := transport.execute(req, "POST", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute POST with nil body failed: %v", err)
	}
}

func TestTransport_Put_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("PUT", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"updated": "true"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetBody(map[string]string{"name": "updated"})
	_, err := transport.execute(req, "PUT", "/v1/test/123", &result)

	if err != nil {
		t.Fatalf("execute PUT failed: %v", err)
	}

	if result["updated"] != "true" {
		t.Errorf("result['updated'] = %v, want 'true'", result["updated"])
	}
}

func TestTransport_Patch_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("PATCH", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"patched": "true"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetBody(map[string]string{"status": "active"})
	_, err := transport.execute(req, "PATCH", "/v1/test/123", &result)

	if err != nil {
		t.Fatalf("execute PATCH failed: %v", err)
	}
}

func TestTransport_Delete_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("DELETE", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(204, nil))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background())
	_, err := transport.execute(req, "DELETE", "/v1/test/123", &result)

	if err != nil {
		t.Fatalf("execute DELETE failed: %v", err)
	}
}

func TestTransport_Delete_WithQueryParams(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/test\?`,
		httpmock.NewJsonResponderOrPanic(204, nil))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetQueryParam("force", "true")
	_, err := transport.execute(req, "DELETE", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute DELETE with query params failed: %v", err)
	}
}

func TestTransport_DeleteWithBody_Success(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("DELETE", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"deleted": "2"}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background()).
		SetBody(map[string][]string{"ids": {"id1", "id2"}})
	_, err := transport.execute(req, "DELETE", "/v1/test", &result)

	if err != nil {
		t.Fatalf("execute DELETE with body failed: %v", err)
	}
}

func TestTransport_HTTPError(t *testing.T) {
	transport := setupTestTransport(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(404, map[string]any{
			"errors": []map[string]string{
				{
					"status": "404",
					"code":   "NOT_FOUND",
					"detail": "Resource not found",
				},
			},
		}))

	var result map[string]string
	req := transport.httpClient.R().SetContext(context.Background())
	_, err := transport.execute(req, "GET", "/v1/test", &result)

	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	transport := setupTestTransport(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]string
	req := transport.httpClient.R().SetContext(ctx)
	_, err := transport.execute(req, "GET", "/v1/test", &result)

	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
}

func TestTransport_GetHTTPClient(t *testing.T) {
	transport := setupTestTransport(t)

	httpClient := transport.GetHTTPClient()

	if httpClient == nil {
		t.Error("GetHTTPClient returned nil")
	}
}

func TestTransport_Close(t *testing.T) {
	transport := setupTestTransport(t)

	err := transport.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestTransport_QueryBuilder(t *testing.T) {
	transport := setupTestTransport(t)

	qb := transport.QueryBuilder()

	if qb == nil {
		t.Error("QueryBuilder returned nil")
	}

	// Test it returns a functional query builder
	qb.AddString("test", "value")
	if !qb.Has("test") {
		t.Error("QueryBuilder not functional")
	}
}

func TestExecute_UnsupportedMethod(t *testing.T) {
	transport := setupTestTransport(t)

	req := transport.httpClient.R()
	_, err := transport.execute(req, "INVALID", "/test", nil)

	if err == nil {
		t.Error("Expected error for unsupported HTTP method, got nil")
	}
}

func TestTransport_AllHTTPMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"GET", "GET"},
		{"POST", "POST"},
		{"PUT", "PUT"},
		{"PATCH", "PATCH"},
		{"DELETE", "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := setupTestTransport(t)

			httpmock.RegisterResponder(tt.method, "https://api-business.apple.com/v1/test",
				httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

			var result map[string]string
			req := transport.httpClient.R().SetContext(context.Background())
			_, err := transport.execute(req, tt.method, "/v1/test", &result)
			if err != nil {
				t.Errorf("%s request failed: %v", tt.method, err)
			}
		})
	}
}

func TestTransport_HTTPErrorStatuses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"BadRequest", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"NotFound", http.StatusNotFound},
		{"TooManyRequests", http.StatusTooManyRequests},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := setupTestTransport(t)

			errorResponse := map[string]any{
				"errors": []map[string]string{
					{
						"status": "error",
						"code":   "ERROR",
						"detail": "Test error",
					},
				},
			}

			httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
				httpmock.NewJsonResponderOrPanic(tt.statusCode, errorResponse))

			var result map[string]string
			req := transport.httpClient.R().SetContext(context.Background())
			_, err := transport.execute(req, "GET", "/v1/test", &result)

			if err == nil {
				t.Errorf("Expected error for status %d, got nil", tt.statusCode)
			}
		})
	}
}
