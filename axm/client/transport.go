package client

import (
	"fmt"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/interfaces"
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

// APIResponse represents the standard API response structure
type APIResponse[T any] struct {
	Data  []T   `json:"data"`
	Meta  Meta  `json:"meta"`
	Links Links `json:"links"`
}

// NewTransport creates a new HTTP transport for Apple Business Manager API.
// This is an internal function - users should use axm.NewClient() instead.
func NewTransport(keyID, issuerID string, privateKey any, options ...ClientOption) (*Client, error) {
	// Validate required parameters
	if keyID == "" {
		return nil, fmt.Errorf("keyID is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("issuerID is required")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("privateKey is required")
	}

	logger := zap.NewNop()

	// Create JWT authentication provider
	auth := NewJWTAuth(JWTAuthConfig{
		KeyID:      keyID,
		IssuerID:   issuerID,
		PrivateKey: privateKey,
		Audience:   DefaultJWTAudience,
		Scope:      ScopeBusinessAPI,
	})

	// Create resty HTTP client with defaults
	httpClient := resty.New()
	httpClient.
		SetBaseURL(DefaultBaseURL).
		SetTimeout(30*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetHeader("User-Agent", DefaultUserAgent)

	errorHandler := NewErrorHandler(logger)

	// Create client instance
	client := &Client{
		httpClient:   httpClient,
		logger:       logger,
		auth:         auth,
		errorHandler: errorHandler,
		baseURL:      DefaultBaseURL,
	}

	// Apply any additional options
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	// Setup authentication middleware (after options to use configured logger)
	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		if err := client.auth.ApplyAuth(req); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		client.logger.Info("API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)

		return nil
	})

	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		client.logger.Info("API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)

		if resp.StatusCode() == 401 {
			if jwtAuth, ok := client.auth.(*JWTAuth); ok {
				client.logger.Info("Received 401 response, forcing JWT token refresh")
				jwtAuth.ForceRefresh()
			}
		}

		return nil
	})

	client.logger.Info("Apple Business Manager API client created",
		zap.String("issuer_id", issuerID),
		zap.String("base_url", client.baseURL))

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

// NewTransportFromEnv creates a transport using environment variables
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
func NewTransportFromEnv(options ...ClientOption) (*Client, error) {
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

	return NewTransport(keyID, issuerID, privateKey, options...)
}

// NewTransportFromFile creates a transport using credentials from files
func NewTransportFromFile(keyID, issuerID, privateKeyPath string, options ...ClientOption) (*Client, error) {
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

	return NewTransport(keyID, issuerID, privateKey, options...)
}
