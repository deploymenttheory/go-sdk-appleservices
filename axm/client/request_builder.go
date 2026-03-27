package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"resty.dev/v3"
)

// requestExecutor is the execution backend for a RequestBuilder.
// Transport implements it directly; tests supply a mock via NewMockRequestBuilder.
type requestExecutor interface {
	execute(req *resty.Request, method, path string, result any) (*resty.Response, error)
	executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error)
	executePaginated(req *resty.Request, path string, mergePage func([]byte) error) (*resty.Response, error)
}

// RequestBuilder constructs a single API request. The service layer owns the
// full request shape — headers, body, query params, result target — before
// handing the completed request to the executor (transport) which handles
// auth, retry, and throttling.
//
// Usage:
//
//	resp, err := s.client.NewRequest(ctx).
//	    SetHeader("Accept", constants.ApplicationJSON).
//	    SetHeader("Content-Type", constants.ApplicationJSON).
//	    SetBody(payload).
//	    SetResult(&result).
//	    Post(constants.EndpointOrgDeviceActivities)
type RequestBuilder struct {
	req      *resty.Request
	executor requestExecutor
	result   any
}

// SetHeader sets a request-level header. Empty values are ignored.
func (b *RequestBuilder) SetHeader(key, value string) *RequestBuilder {
	if value != "" {
		b.req.SetHeader(key, value)
	}
	return b
}

// SetQueryParam adds a URL query parameter. Empty values are ignored.
func (b *RequestBuilder) SetQueryParam(key, value string) *RequestBuilder {
	if value != "" {
		b.req.SetQueryParam(key, value)
	}
	return b
}

// SetQueryParams adds multiple URL query parameters in bulk. Empty values are ignored.
func (b *RequestBuilder) SetQueryParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		if v != "" {
			b.req.SetQueryParam(k, v)
		}
	}
	return b
}

// SetBody sets the request body. Nil is ignored.
func (b *RequestBuilder) SetBody(body any) *RequestBuilder {
	if body != nil {
		b.req.SetBody(body)
	}
	return b
}

// SetResult sets the target for JSON unmarshaling of a successful response.
func (b *RequestBuilder) SetResult(result any) *RequestBuilder {
	b.result = result
	b.req.SetResult(result)
	return b
}

// SetMultipartFile configures the request for a multipart file upload.
// Content-Type is managed automatically by resty.
func (b *RequestBuilder) SetMultipartFile(fileField, fileName string, fileReader io.Reader, fileSize int64) *RequestBuilder {
	if fileReader != nil && fileName != "" && fileField != "" {
		field := &resty.MultipartField{
			Name:        fileField,
			FileName:    fileName,
			ContentType: "application/octet-stream",
			Reader:      fileReader,
			FileSize:    fileSize,
		}
		b.req.SetMultipartFields(field)
	}
	return b
}

// SetMultipartFormData adds additional form fields to a multipart request.
func (b *RequestBuilder) SetMultipartFormData(formFields map[string]string) *RequestBuilder {
	if len(formFields) > 0 {
		b.req.SetMultipartFormData(formFields)
	}
	return b
}

// Get executes the request as GET against path.
func (b *RequestBuilder) Get(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "GET", path, b.result)
}

// Post executes the request as POST against path.
func (b *RequestBuilder) Post(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "POST", path, b.result)
}

// Put executes the request as PUT against path.
func (b *RequestBuilder) Put(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "PUT", path, b.result)
}

// Patch executes the request as PATCH against path.
func (b *RequestBuilder) Patch(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "PATCH", path, b.result)
}

// Delete executes the request as DELETE against path.
func (b *RequestBuilder) Delete(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "DELETE", path, b.result)
}

// GetBytes executes a GET request and returns raw response bytes without JSON
// unmarshaling. Use for binary responses such as files or exports.
func (b *RequestBuilder) GetBytes(path string) (*resty.Response, []byte, error) {
	return b.executor.executeGetBytes(b.req, path)
}

// GetPaginated transparently fetches all pages of a cursor-based paginated
// Apple Business Manager API endpoint, calling mergePage with each page's raw
// JSON response. Query parameters already set on the builder (via SetQueryParam
// / SetQueryParams) are forwarded as the base params; cursor management is
// handled internally by the transport.
func (b *RequestBuilder) GetPaginated(path string, mergePage func([]byte) error) (*resty.Response, error) {
	return b.executor.executePaginated(b.req, path, mergePage)
}

// mockRequestExecutor backs a RequestBuilder in tests, routing execution
// through a caller-supplied dispatch function instead of a real Transport.
type mockRequestExecutor struct {
	fn              func(method, path string, result any) (*resty.Response, error)
	queryParamStore *map[string]string
}

func (m *mockRequestExecutor) execute(req *resty.Request, method, path string, result any) (*resty.Response, error) {
	m.captureQueryParams(req)
	return m.fn(method, path, result)
}

func (m *mockRequestExecutor) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	m.captureQueryParams(req)
	resp, err := m.fn("GET", path, nil)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}

func (m *mockRequestExecutor) executePaginated(req *resty.Request, path string, mergePage func([]byte) error) (*resty.Response, error) {
	m.captureQueryParams(req)
	resp, err := m.fn("GET", path, nil)
	if err != nil {
		return resp, err
	}
	body := resp.Bytes()
	if mergePage != nil && len(body) > 0 {
		var pageResp struct {
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(body, &pageResp) == nil && len(pageResp.Data) > 0 {
			if err := mergePage(body); err != nil {
				return resp, fmt.Errorf("mergePage failed: %w", err)
			}
		} else {
			if err := mergePage(body); err != nil {
				return resp, fmt.Errorf("mergePage failed: %w", err)
			}
		}
	}
	return resp, nil
}

func (m *mockRequestExecutor) captureQueryParams(req *resty.Request) {
	if m.queryParamStore != nil && req != nil {
		params := make(map[string]string)
		for k, v := range req.QueryParams {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		if len(params) > 0 {
			*m.queryParamStore = params
		}
	}
}

// NewMockRequestBuilder returns a RequestBuilder suitable for unit tests.
// The fn callback receives the HTTP method, path, and result pointer and
// returns a pre-programmed response.
func NewMockRequestBuilder(ctx context.Context, fn func(method, path string, result any) (*resty.Response, error)) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R().SetContext(ctx),
		executor: &mockRequestExecutor{fn: fn, queryParamStore: nil},
	}
}

// NewMockRequestBuilderWithQueryCapture returns a RequestBuilder suitable for
// unit tests that also captures query parameters into the provided map pointer.
func NewMockRequestBuilderWithQueryCapture(ctx context.Context, fn func(method, path string, result any) (*resty.Response, error), queryStore *map[string]string) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R().SetContext(ctx),
		executor: &mockRequestExecutor{fn: fn, queryParamStore: queryStore},
	}
}
