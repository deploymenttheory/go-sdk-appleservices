package axm

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// ClientBuilder provides a fluent interface for building AXM clients
type ClientBuilder struct {
	keyID      string
	issuerID   string
	privateKey any // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	audience   string
	baseURL    string
	timeout    time.Duration
	retryCount int
	retryWait  time.Duration
	userAgent  string
	debug      bool
	logger     *zap.Logger
}

// NewClientBuilder creates a new client builder with default values
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		audience:   "appstoreconnect-v1",
		baseURL:    "https://api-business.apple.com/v1",
		timeout:    30 * time.Second,
		retryCount: 3,
		retryWait:  1 * time.Second,
		userAgent:  "go-api-sdk-apple/1.0.0",
		debug:      false,
	}
}

// WithJWTAuth configures JWT authentication with provided credentials
func (cb *ClientBuilder) WithJWTAuth(keyID, issuerID string, privateKey any) *ClientBuilder {
	cb.keyID = keyID
	cb.issuerID = issuerID
	cb.privateKey = privateKey
	return cb
}

// WithJWTAuthFromFile configures JWT authentication from a private key file
func (cb *ClientBuilder) WithJWTAuthFromFile(keyID, issuerID, privateKeyPath string) *ClientBuilder {
	privateKey, err := LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		// Store error for later handling in Build()
		cb.privateKey = nil
		return cb
	}
	return cb.WithJWTAuth(keyID, issuerID, privateKey)
}

// WithJWTAuthFromEnv configures JWT authentication from environment variables
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
func (cb *ClientBuilder) WithJWTAuthFromEnv() *ClientBuilder {
	keyID := os.Getenv("APPLE_KEY_ID")
	issuerID := os.Getenv("APPLE_ISSUER_ID")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	if keyID == "" || issuerID == "" || privateKeyPath == "" {
		// Store error state for later handling in Build()
		cb.keyID = ""
		cb.issuerID = ""
		cb.privateKey = nil
		return cb
	}

	return cb.WithJWTAuthFromFile(keyID, issuerID, privateKeyPath)
}

// WithBaseURL sets the base URL for the API (default: https://api-business.apple.com/v1)
func (cb *ClientBuilder) WithBaseURL(baseURL string) *ClientBuilder {
	cb.baseURL = baseURL
	return cb
}

// WithTimeout sets the request timeout (default: 30 seconds)
func (cb *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
	cb.timeout = timeout
	return cb
}

// WithRetry configures retry settings (default: 3 retries, 1 second wait)
func (cb *ClientBuilder) WithRetry(count int, wait time.Duration) *ClientBuilder {
	cb.retryCount = count
	cb.retryWait = wait
	return cb
}

// WithUserAgent sets the user agent string (default: go-api-sdk-apple/1.0.0)
func (cb *ClientBuilder) WithUserAgent(userAgent string) *ClientBuilder {
	cb.userAgent = userAgent
	return cb
}

// WithDebug enables or disables debug mode (default: false)
func (cb *ClientBuilder) WithDebug(debug bool) *ClientBuilder {
	cb.debug = debug
	return cb
}

// WithLogger sets a custom logger (default: no-op logger, or development logger if debug is enabled)
func (cb *ClientBuilder) WithLogger(logger *zap.Logger) *ClientBuilder {
	cb.logger = logger
	return cb
}

// WithAudience sets the JWT audience (default: appstoreconnect-v1)
func (cb *ClientBuilder) WithAudience(audience string) *ClientBuilder {
	cb.audience = audience
	return cb
}

// Build creates and returns the configured AXM client
func (cb *ClientBuilder) Build() (*Client, error) {
	// Validate authentication credentials
	if cb.keyID == "" || cb.issuerID == "" || cb.privateKey == nil {
		return nil, fmt.Errorf("JWT authentication credentials are required (keyID, issuerID, and privateKey)")
	}

	// Validate private key
	if err := ValidatePrivateKey(cb.privateKey); err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Create JWT auth provider
	jwtAuth := NewJWTAuth(JWTAuthConfig{
		KeyID:      cb.keyID,
		IssuerID:   cb.issuerID,
		PrivateKey: cb.privateKey,
		Audience:   cb.audience,
	})

	// Set default logger if none provided
	if cb.logger == nil {
		if cb.debug {
			cb.logger, _ = zap.NewDevelopment()
		} else {
			cb.logger = zap.NewNop()
		}
	}

	// Validate configuration
	if cb.baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if cb.timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	if cb.retryCount < 0 {
		return nil, fmt.Errorf("retry count cannot be negative")
	}

	// Create client config
	config := Config{
		BaseURL:    cb.baseURL,
		Auth:       jwtAuth,
		Logger:     cb.logger,
		Timeout:    cb.timeout,
		RetryCount: cb.retryCount,
		RetryWait:  cb.retryWait,
		UserAgent:  cb.userAgent,
		Debug:      cb.debug,
	}

	return NewClient(config)
}

// MustBuild creates the AXM client and panics if there's an error
// Use this only when you're certain the configuration is valid
func (cb *ClientBuilder) MustBuild() *Client {
	client, err := cb.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build AXM client: %v", err))
	}
	return client
}

// Validate checks the current configuration without building the client
func (cb *ClientBuilder) Validate() error {
	if cb.keyID == "" {
		return fmt.Errorf("key ID is required")
	}
	if cb.issuerID == "" {
		return fmt.Errorf("issuer ID is required")
	}
	if cb.privateKey == nil {
		return fmt.Errorf("private key is required")
	}
	if err := ValidatePrivateKey(cb.privateKey); err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	if cb.baseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	if cb.timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if cb.retryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}
	return nil
}

// Clone creates a copy of the client builder
func (cb *ClientBuilder) Clone() *ClientBuilder {
	return &ClientBuilder{
		keyID:      cb.keyID,
		issuerID:   cb.issuerID,
		privateKey: cb.privateKey,
		audience:   cb.audience,
		baseURL:    cb.baseURL,
		timeout:    cb.timeout,
		retryCount: cb.retryCount,
		retryWait:  cb.retryWait,
		userAgent:  cb.userAgent,
		debug:      cb.debug,
		logger:     cb.logger,
	}
}

// Public convenience functions for common client creation patterns

// NewClientFromEnv creates an AXM client using environment variables
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
func NewClientFromEnv() (*Client, error) {
	return NewClientBuilder().
		WithJWTAuthFromEnv().
		Build()
}

// NewClientFromFile creates an AXM client using credentials from files
func NewClientFromFile(keyID, issuerID, privateKeyPath string) (*Client, error) {
	return NewClientBuilder().
		WithJWTAuthFromFile(keyID, issuerID, privateKeyPath).
		Build()
}

// NewClientFromEnvWithOptions creates an AXM client from environment variables with custom options
func NewClientFromEnvWithOptions(debug bool, timeout time.Duration, userAgent string) (*Client, error) {
	builder := NewClientBuilder().
		WithJWTAuthFromEnv().
		WithDebug(debug).
		WithUserAgent(userAgent)

	if timeout > 0 {
		builder = builder.WithTimeout(timeout)
	}

	return builder.Build()
}

// NewClientFromFileWithOptions creates an AXM client from files with custom options
func NewClientFromFileWithOptions(keyID, issuerID, privateKeyPath string, debug bool, timeout time.Duration, userAgent string) (*Client, error) {
	builder := NewClientBuilder().
		WithJWTAuthFromFile(keyID, issuerID, privateKeyPath).
		WithDebug(debug).
		WithUserAgent(userAgent)

	if timeout > 0 {
		builder = builder.WithTimeout(timeout)
	}

	return builder.Build()
}