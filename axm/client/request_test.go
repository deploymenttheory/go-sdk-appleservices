package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	req.SetHeader("Authorization", "Bearer mock-token")
	return nil
}

func setupMockClient(t *testing.T) *Client {
	mockAuth := &MockAuthProvider{}

	client := &Client{
		httpClient:   resty.New(),
		logger:       zap.NewNop(),
		auth:         mockAuth,
		errorHandler: NewErrorHandler(zap.NewNop()),
		baseURL:      "https://api-business.apple.com",
	}

	client.httpClient.SetBaseURL(client.baseURL)

	// Setup auth middleware
	client.httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		return client.auth.ApplyAuth(req)
	})

	httpmock.ActivateNonDefault(client.httpClient.Client())

	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return client
}

func TestClient_Get_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	err := client.Get(context.Background(), "/v1/test", nil, nil, &result)

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("result['status'] = %v, want 'ok'", result["status"])
	}
}

func TestClient_Get_WithQueryParams(t *testing.T) {
	client := setupMockClient(t)

	// Use a more flexible responder that matches any query string
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/test\?`,
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	queryParams := map[string]string{
		"limit":  "10",
		"cursor": "abc",
	}

	var result map[string]string
	err := client.Get(context.Background(), "/v1/test", queryParams, nil, &result)

	if err != nil {
		t.Fatalf("Get with query params failed: %v", err)
	}
}

func TestClient_Get_WithHeaders(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	headers := map[string]string{
		"X-Custom-Header": "custom-value",
		"Accept":          "application/json",
	}

	var result map[string]string
	err := client.Get(context.Background(), "/v1/test", nil, headers, &result)

	if err != nil {
		t.Fatalf("Get with headers failed: %v", err)
	}
}

func TestClient_Post_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(201, map[string]string{"id": "12345"}))

	requestBody := map[string]string{"name": "test"}
	var result map[string]string
	err := client.Post(context.Background(), "/v1/test", requestBody, nil, &result)

	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}

	if result["id"] != "12345" {
		t.Errorf("result['id'] = %v, want '12345'", result["id"])
	}
}

func TestClient_Post_NilBody(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	err := client.Post(context.Background(), "/v1/test", nil, nil, &result)

	if err != nil {
		t.Fatalf("Post with nil body failed: %v", err)
	}
}

func TestClient_PostWithQuery_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/test\?`,
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	queryParams := map[string]string{"mode": "sync"}
	requestBody := map[string]string{"data": "value"}
	var result map[string]string

	err := client.PostWithQuery(context.Background(), "/v1/test", queryParams, requestBody, nil, &result)

	if err != nil {
		t.Fatalf("PostWithQuery failed: %v", err)
	}
}

func TestClient_Put_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("PUT", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"updated": "true"}))

	requestBody := map[string]string{"name": "updated"}
	var result map[string]string
	err := client.Put(context.Background(), "/v1/test/123", requestBody, nil, &result)

	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	if result["updated"] != "true" {
		t.Errorf("result['updated'] = %v, want 'true'", result["updated"])
	}
}

func TestClient_Patch_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("PATCH", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"patched": "true"}))

	requestBody := map[string]string{"status": "active"}
	var result map[string]string
	err := client.Patch(context.Background(), "/v1/test/123", requestBody, nil, &result)

	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
}

func TestClient_Delete_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("DELETE", "https://api-business.apple.com/v1/test/123",
		httpmock.NewJsonResponderOrPanic(204, nil))

	var result map[string]string
	err := client.Delete(context.Background(), "/v1/test/123", nil, nil, &result)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestClient_Delete_WithQueryParams(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/test\?`,
		httpmock.NewJsonResponderOrPanic(204, nil))

	queryParams := map[string]string{"force": "true"}
	var result map[string]string
	err := client.Delete(context.Background(), "/v1/test", queryParams, nil, &result)

	if err != nil {
		t.Fatalf("Delete with query params failed: %v", err)
	}
}

func TestClient_DeleteWithBody_Success(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("DELETE", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"deleted": "2"}))

	requestBody := map[string][]string{"ids": {"id1", "id2"}}
	var result map[string]string
	err := client.DeleteWithBody(context.Background(), "/v1/test", requestBody, nil, &result)

	if err != nil {
		t.Fatalf("DeleteWithBody failed: %v", err)
	}
}

func TestClient_HTTPError(t *testing.T) {
	client := setupMockClient(t)

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
	err := client.Get(context.Background(), "/v1/test", nil, nil, &result)

	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	client := setupMockClient(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]string
	err := client.Get(ctx, "/v1/test", nil, nil, &result)

	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
}

func TestClient_EmptyQueryParams(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	// Query params with empty values should be skipped
	queryParams := map[string]string{
		"key1": "value1",
		"key2": "",
		"key3": "value3",
	}

	var result map[string]string
	err := client.Get(context.Background(), "/v1/test", queryParams, nil, &result)

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestClient_EmptyHeaders(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	// Headers with empty values should be skipped
	headers := map[string]string{
		"Header1": "value1",
		"Header2": "",
		"Header3": "value3",
	}

	var result map[string]string
	err := client.Get(context.Background(), "/v1/test", nil, headers, &result)

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestClient_GetHTTPClient(t *testing.T) {
	client := setupMockClient(t)

	httpClient := client.GetHTTPClient()

	if httpClient == nil {
		t.Error("GetHTTPClient returned nil")
	}
}

func TestClient_Close(t *testing.T) {
	client := setupMockClient(t)

	err := client.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestClient_QueryBuilder(t *testing.T) {
	client := setupMockClient(t)

	qb := client.QueryBuilder()

	if qb == nil {
		t.Error("QueryBuilder returned nil")
	}

	// Test it returns a functional query builder
	qb.AddString("test", "value")
	if !qb.Has("test") {
		t.Error("QueryBuilder not functional")
	}
}

func TestExecuteRequest_UnsupportedMethod(t *testing.T) {
	client := setupMockClient(t)

	req := client.httpClient.R()
	err := client.executeRequest(req, "INVALID", "/test")

	if err == nil {
		t.Error("Expected error for unsupported HTTP method, got nil")
	}
}

func TestClient_AllHTTPMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
		fn     func(*Client, context.Context, string) error
	}{
		{
			name:   "GET",
			method: "GET",
			fn: func(c *Client, ctx context.Context, path string) error {
				var result map[string]string
				return c.Get(ctx, path, nil, nil, &result)
			},
		},
		{
			name:   "POST",
			method: "POST",
			fn: func(c *Client, ctx context.Context, path string) error {
				var result map[string]string
				return c.Post(ctx, path, nil, nil, &result)
			},
		},
		{
			name:   "PUT",
			method: "PUT",
			fn: func(c *Client, ctx context.Context, path string) error {
				var result map[string]string
				return c.Put(ctx, path, nil, nil, &result)
			},
		},
		{
			name:   "PATCH",
			method: "PATCH",
			fn: func(c *Client, ctx context.Context, path string) error {
				var result map[string]string
				return c.Patch(ctx, path, nil, nil, &result)
			},
		},
		{
			name:   "DELETE",
			method: "DELETE",
			fn: func(c *Client, ctx context.Context, path string) error {
				var result map[string]string
				return c.Delete(ctx, path, nil, nil, &result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := setupMockClient(t)

			httpmock.RegisterResponder(tt.method, "https://api-business.apple.com/v1/test",
				httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

			err := tt.fn(client, context.Background(), "/v1/test")
			if err != nil {
				t.Errorf("%s request failed: %v", tt.method, err)
			}
		})
	}
}

func TestClient_HTTPErrorStatuses(t *testing.T) {
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
			client := setupMockClient(t)

			errorResponse := map[string]any{
				"errors": []map[string]string{
					{
						"status": string(rune(tt.statusCode)),
						"code":   "ERROR",
						"detail": "Test error",
					},
				},
			}

			httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
				httpmock.NewJsonResponderOrPanic(tt.statusCode, errorResponse))

			var result map[string]string
			err := client.Get(context.Background(), "/v1/test", nil, nil, &result)

			if err == nil {
				t.Errorf("Expected error for status %d, got nil", tt.statusCode)
			}
		})
	}
}

func TestClient_PostWithQuery_EmptyQuery(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	err := client.PostWithQuery(context.Background(), "/v1/test", nil, map[string]string{"data": "test"}, nil, &result)

	if err != nil {
		t.Fatalf("PostWithQuery with nil query failed: %v", err)
	}
}

func TestClient_DeleteWithBody_NilBody(t *testing.T) {
	client := setupMockClient(t)

	httpmock.RegisterResponder("DELETE", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(204, nil))

	var result map[string]string
	err := client.DeleteWithBody(context.Background(), "/v1/test", nil, nil, &result)

	if err != nil {
		t.Fatalf("DeleteWithBody with nil body failed: %v", err)
	}
}
