package client

import (
	"fmt"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/interfaces"
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

// Ensure Client implements HTTPClient interface
var _ interfaces.HTTPClient = (*Client)(nil)

// QueryBuilder returns a new query builder instance
func (c *Client) QueryBuilder() interfaces.ServiceQueryBuilder {
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

// NewClientFromEnv creates a client using environment variables
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
func NewClientFromEnv() (*Client, error) {
	keyID := os.Getenv("APPLE_KEY_ID")
	issuerID := os.Getenv("APPLE_ISSUER_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	if keyID == "" {
		return nil, fmt.Errorf("APPLE_KEY_ID environment variable is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("APPLE_ISSUER_ID environment variable is required")
	}
	if privateKeyPath == "" {
		return nil, fmt.Errorf("APPLE_PRIVATE_KEY_PATH environment variable is required")
	}

	privateKey, err := LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	config := Config{
		BaseURL: "https://api-business.apple.com/v1",
		Auth: NewJWTAuth(JWTAuthConfig{
			KeyID:      keyID,
			IssuerID:   issuerID,
			PrivateKey: privateKey,
			Audience:   "appstoreconnect-v1",
		}),
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryWait:  1 * time.Second,
		UserAgent:  "go-api-sdk-apple/3.0.0",
	}

	return NewClient(config)
}

// NewClientFromFile creates a client using credentials from files
func NewClientFromFile(keyID, issuerID, privateKeyPath string) (*Client, error) {
	if keyID == "" {
		return nil, fmt.Errorf("keyID is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("issuerID is required")
	}
	if privateKeyPath == "" {
		return nil, fmt.Errorf("privateKeyPath is required")
	}

	privateKey, err := LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	config := Config{
		BaseURL: "https://api-business.apple.com/v1",
		Auth: NewJWTAuth(JWTAuthConfig{
			KeyID:      keyID,
			IssuerID:   issuerID,
			PrivateKey: privateKey,
			Audience:   "appstoreconnect-v1",
		}),
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryWait:  1 * time.Second,
		UserAgent:  "go-api-sdk-apple/3.0.0",
	}

	return NewClient(config)
}
