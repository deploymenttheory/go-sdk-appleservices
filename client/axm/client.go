package axm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// Client represents the main Apple Business Manager API client
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

// NewClient creates a new Apple Business Manager API client
func NewClient(config Config) (*Client, error) {
	// Set defaults
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
		config.UserAgent = "go-api-sdk-apple/1.0.0"
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	// Validate required fields
	if config.Auth == nil {
		return nil, fmt.Errorf("auth provider is required")
	}

	// Create Resty client with v3 best practices
	httpClient := resty.New()

	// Configure client settings
	httpClient.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWait).
		SetRetryMaxWaitTime(config.RetryWait*10).
		SetHeader("User-Agent", config.UserAgent)

	// Enable debug mode if requested
	if config.Debug {
		httpClient.SetDebug(true)
	}

	// Add request middleware for logging and auth
	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		// Apply authentication
		if err := config.Auth.ApplyAuth(req); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		// Log request
		config.Logger.Info("API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)

		return nil
	})

	// Add response middleware for logging and OAuth token refresh
	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		config.Logger.Info("API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)

		// Handle 401 responses for OAuth2Auth by forcing token refresh
		if resp.StatusCode() == 401 {
			if oauth2Auth, ok := config.Auth.(*OAuth2Auth); ok {
				config.Logger.Info("Received 401 response, forcing OAuth token refresh")
				oauth2Auth.ForceRefresh()
			}
		}

		return nil
	})

	client := &Client{
		httpClient:   httpClient,
		logger:       config.Logger,
		auth:         config.Auth,
		errorHandler: NewErrorHandler(config.Logger),
		baseURL:      config.BaseURL,
	}

	return client, nil
}

// Close closes the HTTP client and cleans up resources
func (c *Client) Close() error {
	if c.httpClient != nil {
		c.httpClient.Close()
	}
	return nil
}

// QueryBuilder returns a new query builder instance
func (c *Client) QueryBuilder() *QueryBuilder {
	return NewQueryBuilder()
}

// GetHTTPClient returns the underlying HTTP client for testing purposes
func (c *Client) GetHTTPClient() *resty.Client {
	return c.httpClient
}
