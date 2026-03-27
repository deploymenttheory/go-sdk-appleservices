package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport represents the iTunes Search API HTTP transport layer.
type Transport struct {
	httpClient   *resty.Client
	logger       *zap.Logger
	errorHandler *ErrorHandler
	baseURL      string
}

// Ensure Transport implements Client interface.
var _ Client = (*Transport)(nil)

// NewTransport creates a new HTTP transport for the iTunes Search API.
// This is an internal function — users should use itunes.NewClient() instead.
func NewTransport(options ...ClientOption) (*Transport, error) {
	logger := zap.NewNop()

	httpClient := resty.New()
	httpClient.
		SetBaseURL(DefaultBaseURL).
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
		baseURL:      DefaultBaseURL,
	}

	for _, option := range options {
		if err := option(transport); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		transport.logger.Info("iTunes API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)
		return nil
	})

	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		transport.logger.Info("iTunes API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)
		return nil
	})

	transport.logger.Info("iTunes Search API client created",
		zap.String("base_url", transport.baseURL))

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
//
// The iTunes Search API returns Content-Type: text/javascript rather than
// application/json, so resty's automatic JSON unmarshaling is bypassed.
// We unmarshal the raw body directly instead.
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
	resp, err := t.execute(req, path, nil)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}
