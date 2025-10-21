package axm

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// Client represents the main Apple Business Manager API client with embedded services
type Client struct {
	httpClient   *resty.Client
	logger       *zap.Logger
	auth         AuthProvider
	errorHandler *ErrorHandler
	baseURL      string

}

// Config holds configuration for the client
type Config struct {
	BaseURL    string
	Auth       AuthProvider
	Logger     *zap.Logger
	Timeout    time.Duration
	RetryCount int
	RetryWait  time.Duration
	UserAgent  string
	Debug      bool
}

// APIResponse represents the standard API response structure
type APIResponse[T any] struct {
	Data  []T   `json:"data"`
	Meta  Meta  `json:"meta"`
	Links Links `json:"links"`
}

// NewClient creates a new Apple Business Manager API client with embedded services
func NewClient(config Config) (*Client, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api-business.apple.com/v1"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryWait == 0 {
		config.RetryWait = 1 * time.Second
	}
	if config.UserAgent == "" {
		config.UserAgent = "go-api-sdk-apple/3.0.0"
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	if config.Auth == nil {
		return nil, fmt.Errorf("auth provider is required")
	}

	httpClient := resty.New()

	httpClient.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWait).
		SetRetryMaxWaitTime(config.RetryWait*10).
		SetHeader("User-Agent", config.UserAgent)

	if config.Debug {
		httpClient.SetDebug(true)
	}

	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		if err := config.Auth.ApplyAuth(req); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		config.Logger.Info("API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)

		return nil
	})

	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		config.Logger.Info("API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)

		if resp.StatusCode() == 401 {
			if oauth2Auth, ok := config.Auth.(*OAuth2Auth); ok {
				config.Logger.Info("Received 401 response, forcing OAuth token refresh")
				oauth2Auth.ForceRefresh()
			}
		}

		return nil
	})

	errorHandler := NewErrorHandler(config.Logger)

	client := &Client{
		httpClient:   httpClient,
		logger:       config.Logger,
		auth:         config.Auth,
		errorHandler: errorHandler,
		baseURL:      config.BaseURL,
	}

	return client, nil
}

// HTTPClient interface that services will use
type HTTPClient interface {
	Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	Post(ctx context.Context, path string, body any, headers map[string]string, result any) error
	PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) error
	Put(ctx context.Context, path string, body any, headers map[string]string, result any) error
	Patch(ctx context.Context, path string, body any, headers map[string]string, result any) error
	Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) error
	PostMultipart(ctx context.Context, path string, files map[string]string, fields map[string]string, result any) error
	GetPaginated(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	GetNextPage(ctx context.Context, nextURL string, headers map[string]string, result any) error
	GetAllPages(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, processPage func([]byte) error) error
	QueryBuilder() ServiceQueryBuilder
}

// ServiceQueryBuilder defines the query builder contract for services
type ServiceQueryBuilder interface {
	AddString(key, value string) *QueryBuilder
	AddInt(key string, value int) *QueryBuilder
	AddInt64(key string, value int64) *QueryBuilder
	AddBool(key string, value bool) *QueryBuilder
	AddTime(key string, value time.Time) *QueryBuilder
	AddStringSlice(key string, values []string) *QueryBuilder
	AddIntSlice(key string, values []int) *QueryBuilder
	AddCustom(key, value string) *QueryBuilder
	AddIfNotEmpty(key, value string) *QueryBuilder
	AddIfTrue(condition bool, key, value string) *QueryBuilder
	Merge(other map[string]string) *QueryBuilder
	Remove(key string) *QueryBuilder
	Has(key string) bool
	Get(key string) string
	Build() map[string]string
	BuildString() string
	Clear() *QueryBuilder
	Count() int
	IsEmpty() bool
}

// Ensure Client implements HTTPClient interface
var _ HTTPClient = (*Client)(nil)

// QueryBuilder returns a new query builder instance
func (c *Client) QueryBuilder() ServiceQueryBuilder {
	return NewQueryBuilder()
}

// GetHTTPClient returns the underlying HTTP client for testing purposes
func (c *Client) GetHTTPClient() *resty.Client {
	return c.httpClient
}

// Close closes the HTTP client and cleans up resources
func (c *Client) Close() error {
	if c.httpClient != nil {
		c.httpClient.Close()
	}
	return nil
}
