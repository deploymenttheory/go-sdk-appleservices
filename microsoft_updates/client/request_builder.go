package client

import (
	"context"
	"io"

	"resty.dev/v3"
)

// requestExecutor is the execution backend for a RequestBuilder.
// Transport implements it directly; tests supply a mock via NewMockRequestBuilder.
type requestExecutor interface {
	execute(req *resty.Request, path string, result any) (*resty.Response, error)
	executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error)
	executeHead(req *resty.Request, path string) (*resty.Response, error)
	executeDownload(req *resty.Request, path string, w io.Writer) (*resty.Response, int64, error)
}

// RequestBuilder constructs a single API request. The service layer owns the
// full request shape — headers, query params, result target — before handing
// the completed request to the executor (transport) which handles retry and
// logging.
//
// Usage:
//
//	var result StandaloneResponse
//	resp, err := s.client.NewRequest(ctx).
//	    SetHeader("Accept", constants.ApplicationXML).
//	    SetResult(&result).
//	    Get(constants.StandaloneCDNWordURL)
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

// SetResult sets the target for JSON unmarshaling of a successful response.
func (b *RequestBuilder) SetResult(result any) *RequestBuilder {
	b.result = result
	return b
}

// Get executes the request as GET against path.
func (b *RequestBuilder) Get(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, path, b.result)
}

// GetBytes executes a GET request and returns raw response bytes without JSON
// unmarshaling. Use for binary or XML responses that require custom parsing.
func (b *RequestBuilder) GetBytes(path string) (*resty.Response, []byte, error) {
	return b.executor.executeGetBytes(b.req, path)
}

// Head executes a HEAD request against path.
// Use for querying file metadata or resolving redirect URLs without downloading.
func (b *RequestBuilder) Head(path string) (*resty.Response, error) {
	return b.executor.executeHead(b.req, path)
}

// Download streams the response body of a GET request into w. The response
// body is never buffered in memory — suitable for large file downloads.
// Returns the resty response (for header/status inspection), the number of
// bytes written, and any error.
func (b *RequestBuilder) Download(path string, w io.Writer) (*resty.Response, int64, error) {
	return b.executor.executeDownload(b.req, path, w)
}

// mockRequestExecutor backs a RequestBuilder in tests, routing execution
// through a caller-supplied dispatch function instead of a real Transport.
type mockRequestExecutor struct {
	fn              func(path string, result any) (*resty.Response, error)
	headFn          func(path string) (*resty.Response, error)
	downloadFn      func(path string, w io.Writer) (*resty.Response, int64, error)
	queryParamStore *map[string]string
}

func (m *mockRequestExecutor) execute(req *resty.Request, path string, result any) (*resty.Response, error) {
	m.captureQueryParams(req)
	return m.fn(path, result)
}

func (m *mockRequestExecutor) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	m.captureQueryParams(req)
	resp, err := m.fn(path, nil)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}

func (m *mockRequestExecutor) executeHead(req *resty.Request, path string) (*resty.Response, error) {
	m.captureQueryParams(req)
	if m.headFn != nil {
		return m.headFn(path)
	}
	return m.fn(path, nil)
}

func (m *mockRequestExecutor) executeDownload(req *resty.Request, path string, w io.Writer) (*resty.Response, int64, error) {
	m.captureQueryParams(req)
	if m.downloadFn != nil {
		return m.downloadFn(path, w)
	}
	return nil, 0, io.ErrUnexpectedEOF
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
// The fn callback receives the path and result pointer and returns a
// pre-programmed response.
func NewMockRequestBuilder(ctx context.Context, fn func(path string, result any) (*resty.Response, error)) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R().SetContext(ctx),
		executor: &mockRequestExecutor{fn: fn},
	}
}

// NewMockRequestBuilderWithQueryCapture returns a RequestBuilder for unit tests
// that also captures query parameters into the provided map pointer.
func NewMockRequestBuilderWithQueryCapture(ctx context.Context, fn func(path string, result any) (*resty.Response, error), queryStore *map[string]string) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R().SetContext(ctx),
		executor: &mockRequestExecutor{fn: fn, queryParamStore: queryStore},
	}
}

// NewMockRequestBuilderWithHead returns a RequestBuilder for unit tests that
// supports both GET and HEAD dispatch via separate callbacks.
func NewMockRequestBuilderWithHead(ctx context.Context, fn func(path string, result any) (*resty.Response, error), headFn func(path string) (*resty.Response, error)) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R().SetContext(ctx),
		executor: &mockRequestExecutor{fn: fn, headFn: headFn},
	}
}
