package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport represents the Microsoft Updates HTTP transport layer.
// This service spans multiple external hosts (officecdnmac.microsoft.com,
// edgeupdates.microsoft.com, g.live.com, itunes.apple.com, learn.microsoft.com),
// so all endpoint constants are full absolute URLs and no base URL is set on the
// underlying HTTP client.
type Transport struct {
	httpClient   *resty.Client
	logger       *zap.Logger
	errorHandler *ErrorHandler
}

// Ensure Transport implements Client interface.
var _ Client = (*Transport)(nil)

// NewTransport creates a new HTTP transport for the Microsoft Updates SDK.
// This is an internal function — users should use microsoft_updates.NewClient() instead.
func NewTransport(options ...ClientOption) (*Transport, error) {
	logger := zap.NewNop()

	httpClient := resty.New()
	httpClient.
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetHeader("User-Agent", DefaultUserAgent)

	errorHandler := NewErrorHandler(logger)

	transport := &Transport{
		httpClient:   httpClient,
		logger:       logger,
		errorHandler: errorHandler,
	}

	for _, option := range options {
		if err := option(transport); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		transport.logger.Info("Microsoft Updates API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)
		return nil
	})

	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		transport.logger.Info("Microsoft Updates API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)
		return nil
	})

	transport.logger.Info("Microsoft Updates SDK client created")

	return transport, nil
}

// NewRequest returns a new RequestBuilder for constructing API requests.
func (t *Transport) NewRequest(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		req:      t.httpClient.R().SetContext(ctx),
		executor: t,
	}
}

// QueryBuilder returns a new query builder instance.
func (t *Transport) QueryBuilder() *QueryBuilder {
	return NewQueryBuilder()
}

// GetLogger returns the configured logger.
func (t *Transport) GetLogger() *zap.Logger {
	return t.logger
}

// GetHTTPClient returns the underlying HTTP client for testing purposes.
func (t *Transport) GetHTTPClient() *resty.Client {
	return t.httpClient
}

// Close closes the HTTP client and cleans up resources.
func (t *Transport) Close() error {
	if t.httpClient != nil {
		t.httpClient.Close()
	}
	return nil
}

// execute implements requestExecutor — handles GET requests and error processing.
func (t *Transport) execute(req *resty.Request, path string, result any) (*resty.Response, error) {
	resp, err := req.Get(path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return resp, t.errorHandler.HandleError(resp)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Bytes(), result); err != nil {
			return resp, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return resp, nil
}

// executeGetBytes implements requestExecutor — returns raw response bytes.
func (t *Transport) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	resp, err := req.Get(path)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return resp, nil, t.errorHandler.HandleError(resp)
	}

	return resp, resp.Bytes(), nil
}

// executeHead implements requestExecutor — issues a HEAD request and returns
// the response (headers only, no body). Used for resolving redirect URLs.
func (t *Transport) executeHead(req *resty.Request, path string) (*resty.Response, error) {
	resp, err := req.Head(path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return resp, t.errorHandler.HandleError(resp)
	}

	return resp, nil
}

// executeDownload implements requestExecutor — streams a GET response body
// into w without buffering it in memory.
func (t *Transport) executeDownload(req *resty.Request, path string, w io.Writer) (*resty.Response, int64, error) {
	resp, err := req.SetDoNotParseResponse(true).Get(path)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		resp.Body.Close()
		return resp, 0, t.errorHandler.HandleError(resp)
	}

	defer resp.Body.Close()

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		return resp, n, fmt.Errorf("failed to stream response body: %w", err)
	}

	return resp, n, nil
}
