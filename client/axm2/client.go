package axm2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"resty.dev/v3"
)

// Client is the main AXM API client following the direct client pattern
type Client struct {
	httpClient    *resty.Client
	tokenProvider *TokenProvider
	logger        *zap.Logger
	config        Config
}

// APIError represents structured API error responses following Resty v3 patterns
type APIError struct {
	StatusCode int             `json:"status_code,omitempty"`
	Message    string          `json:"message,omitempty"`
	Code       string          `json:"code,omitempty"`
	RequestID  string          `json:"request_id,omitempty"`
	Errors     []AppleAPIError `json:"errors,omitempty"` // Apple's error format
}

// AppleAPIError represents Apple's specific error format
type AppleAPIError struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source *struct {
		Parameter string `json:"parameter,omitempty"`
		Pointer   string `json:"pointer,omitempty"`
	} `json:"source,omitempty"`
}

func (e *APIError) Error() string {
	// If we have Apple-specific errors, format them nicely
	if len(e.Errors) > 0 {
		var messages []string
		for _, appleErr := range e.Errors {
			if appleErr.Detail != "" {
				messages = append(messages, fmt.Sprintf("%s: %s", appleErr.Code, appleErr.Detail))
			} else if appleErr.Title != "" {
				messages = append(messages, fmt.Sprintf("%s: %s", appleErr.Code, appleErr.Title))
			}
		}
		if len(messages) > 0 {
			return fmt.Sprintf("Apple API errors: %s", strings.Join(messages, "; "))
		}
	}

	// Fallback to standard error format
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}

// Config holds configuration for the AXM client
type Config struct {
	BaseURL        string        // Apple School/Business Manager API base URL (optional, will be set based on APIType)
	APIType        string        // API type: "abm" for Apple Business Manager, "asm" for Apple School Manager
	ClientID       string        // Client ID from Apple
	KeyID          string        // Key ID from Apple
	PrivateKey     string        // Private key content (PEM format)
	Scope          string        // OAuth scope ("business.api" or "school.api")
	Timeout        time.Duration // HTTP timeout (default: 30s)
	RetryCount     int           // Number of retries (default: 3)
	RetryMinWait   time.Duration // Minimum wait time between retries (default: 1s)
	RetryMaxWait   time.Duration // Maximum wait time between retries (default: 10s)
	UserAgent      string        // User agent string
	Debug          bool          // Enable debug logging
	EnableRetryLog bool          // Enable detailed retry logging (default: true)
}

// NewClient creates a new AXM API client
func NewClient(config Config) (*Client, error) {
	// Configure logger
	logger, err := configureLogger(config.Debug)
	if err != nil {
		return nil, fmt.Errorf("failed to configure logger: %w", err)
	}

	// Set base URL based on API type
	if config.BaseURL == "" {
		switch config.APIType {
		case APITypeABM:
			config.BaseURL = AppleBusinessManagerBaseURL
		case APITypeASM:
			config.BaseURL = AppleSchoolManagerBaseURL
		default:
			config.BaseURL = AppleBusinessManagerBaseURL
		}
	}

	// Set scope based on API type if not specified
	if config.Scope == "" {
		switch config.APIType {
		case APITypeABM:
			config.Scope = BusinessScope
		case APITypeASM:
			config.Scope = SchoolScope
		default:
			config.Scope = BusinessScope
		}
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryMinWait == 0 {
		config.RetryMinWait = 1 * time.Second
	}
	if config.RetryMaxWait == 0 {
		config.RetryMaxWait = 10 * time.Second
	}
	if config.UserAgent == "" {
		config.UserAgent = "go-api-sdk-apple/2.0.0"
	}
	// Enable retry logging by default
	if !config.Debug {
		config.EnableRetryLog = true
	}

	// Create token provider
	tokenProvider, err := NewTokenProvider(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create token provider: %w", err)
	}

	client := &Client{
		tokenProvider: tokenProvider,
		logger:        logger,
		config:        config,
	}

	// Create HTTP client with automatic token injection (Resty v3 API)
	httpClient := resty.New().
		SetTimeout(config.Timeout).
		SetHeader("User-Agent", config.UserAgent).
		SetBaseURL(config.BaseURL)

	// Add request middleware for token injection (Resty v3 pattern)
	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		token, err := client.tokenProvider.GetToken(req.Context())
		if err != nil {
			logger.Error("Failed to get access token", zap.Error(err))
			return fmt.Errorf("failed to get access token: %w", err)
		}
		req.SetAuthToken(token) // Use Resty's built-in auth token method
		return nil
	})

	// Configure retry mechanism (Resty v3 pattern)
	httpClient.
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryMinWait).
		SetRetryMaxWaitTime(config.RetryMaxWait)

	// Add smart retry conditions for Apple AXM API
	httpClient.AddRetryConditions(
		func(resp *resty.Response, err error) bool {
			if err != nil {
				logger.Debug("Retry condition: network error occurred", zap.Error(err))
				return false // Don't retry on network errors
			}

			// Retry on 401 Unauthorized (token expired)
			if resp.StatusCode() == http.StatusUnauthorized {
				logger.Debug("Retry condition: 401 Unauthorized - will refresh token")
				return true
			}

			// Retry on 429 Too Many Requests (rate limiting)
			if resp.StatusCode() == http.StatusTooManyRequests {
				logger.Debug("Retry condition: 429 Too Many Requests - rate limited")
				return true
			}

			// Retry on 5xx server errors (except 501 Not Implemented)
			if resp.StatusCode() >= 500 && resp.StatusCode() != http.StatusNotImplemented {
				logger.Debug("Retry condition: 5xx server error", zap.Int("status_code", resp.StatusCode()))
				return true
			}

			// Don't retry on other status codes
			return false
		},
	)

	// Add retry hooks for token refresh and logging
	httpClient.AddRetryHooks(
		func(resp *resty.Response, err error) {
			// Handle 401 Unauthorized - refresh token before retry
			if resp != nil && resp.StatusCode() == http.StatusUnauthorized {
				logger.Debug("Retry hook: Refreshing token due to 401 Unauthorized")
				if refreshErr := client.tokenProvider.ForceRefresh(context.Background()); refreshErr != nil {
					logger.Error("Failed to refresh token in retry hook", zap.Error(refreshErr))
				} else {
					logger.Debug("Successfully refreshed token in retry hook")
				}
			}

			// Log retry attempt (conditional based on config)
			if config.EnableRetryLog && resp != nil {
				retryReason := "unknown"
				switch resp.StatusCode() {
				case http.StatusUnauthorized:
					retryReason = "token_expired"
				case http.StatusTooManyRequests:
					retryReason = "rate_limited"
				default:
					if resp.StatusCode() >= 500 {
						retryReason = "server_error"
					}
				}

				logger.Info("Retrying request",
					zap.String("method", resp.Request.Method),
					zap.String("url", resp.Request.URL),
					zap.Int("status_code", resp.StatusCode()),
					zap.Int("attempt", resp.Request.Attempt),
					zap.String("retry_reason", retryReason),
					zap.Duration("wait_time", config.RetryMinWait))
			}
		},
	)

	if config.Debug {
		httpClient.SetDebug(true)
	}

	// Add custom content-type encoders and decoders for Apple AXM API
	setupContentTypeHandlers(httpClient, logger)

	client.httpClient = httpClient

	return client, nil
}

// Close cleans up resources following Resty v3 patterns
func (c *Client) Close() {
	if c.httpClient != nil {
		c.httpClient.Close() // Resty v3 requires explicit client close
	}
	if c.logger != nil {
		c.logger.Sync()
	}
}

// GetClientID returns the configured client ID
func (c *Client) GetClientID() string {
	return c.config.ClientID
}

// IsAuthenticated returns true if we have a valid access token
func (c *Client) IsAuthenticated() bool {
	return c.tokenProvider.IsValid()
}

// ForceReauthenticate forces a new authentication cycle
func (c *Client) ForceReauthenticate() error {
	return c.tokenProvider.ForceRefresh(context.Background())
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.config.BaseURL
}

// GetAPIType returns the API type
func (c *Client) GetAPIType() string {
	return c.config.APIType
}

// Helper function to configure logger
func configureLogger(debug bool) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if debug {
		// Development config with colors and console encoder
		developmentConfig := zap.NewDevelopmentConfig()
		developmentConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		developmentConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		developmentConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		logger, err = developmentConfig.Build()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return zap.NewNop(), fmt.Errorf("failed to create logger: %w", err)
	}

	return logger, nil
}

// DoRequest performs a generic HTTP request using Resty v3 patterns
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) (*resty.Response, error) {
	var apiError APIError

	request := c.httpClient.R().
		SetContext(ctx).
		SetError(&apiError)

	// Set result if provided
	if result != nil {
		request.SetResult(result)
	}

	// Set request body for POST/PUT requests
	if body != nil {
		request.SetBody(body)
	}

	var response *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		response, err = request.Get(endpoint)
	case http.MethodPost:
		response, err = request.Post(endpoint)
	case http.MethodPut:
		response, err = request.Put(endpoint)
	case http.MethodDelete:
		response, err = request.Delete(endpoint)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}

	if response.IsError() {
		c.logger.Error("API request failed",
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", apiError))
		return response, fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	return response, nil
}

// SetRetryConfig allows overriding retry settings for specific use cases
func (c *Client) SetRetryConfig(retryCount int, minWait, maxWait time.Duration) *Client {
	c.httpClient.
		SetRetryCount(retryCount).
		SetRetryWaitTime(minWait).
		SetRetryMaxWaitTime(maxWait)
	return c
}

// EnableRetryDefaultConditions enables Resty's default retry conditions
// in addition to the Apple AXM-specific conditions
func (c *Client) EnableRetryDefaultConditions() *Client {
	c.httpClient.EnableRetryDefaultConditions()
	return c
}

// setupContentTypeHandlers configures custom encoders and decoders for Apple AXM API
func setupContentTypeHandlers(client *resty.Client, logger *zap.Logger) {
	// Add custom JSON encoder for Apple AXM API requests
	client.AddContentTypeEncoder("application/json", func(w io.Writer, v any) error {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false) // Apple APIs don't require HTML escaping
		encoder.SetIndent("", "")    // Compact JSON for efficiency

		if err := encoder.Encode(v); err != nil {
			logger.Error("Failed to encode JSON for Apple AXM API", zap.Error(err))
			return fmt.Errorf("JSON encoding error: %w", err)
		}
		return nil
	})

	// Add custom JSON decoder for Apple AXM API responses
	client.AddContentTypeDecoder("application/json", func(r io.Reader, v any) error {
		decoder := json.NewDecoder(r)
		// Don't use DisallowUnknownFields for Apple API - they return extra fields
		// decoder.DisallowUnknownFields() // Too strict for Apple's API responses

		if err := decoder.Decode(v); err != nil {
			logger.Error("Failed to decode JSON from Apple AXM API", zap.Error(err))
			return fmt.Errorf("JSON decoding error: %w", err)
		}
		return nil
	})

	// Add decoder for Apple's error responses (often text/plain or application/problem+json)
	client.AddContentTypeDecoder("text/plain", func(r io.Reader, v any) error {
		body, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read text response: %w", err)
		}

		// Try to parse as JSON first (Apple sometimes returns JSON with text/plain)
		if strings.HasPrefix(strings.TrimSpace(string(body)), "{") {
			if err := json.Unmarshal(body, v); err == nil {
				return nil
			}
		}

		// If not JSON, handle as plain text error
		if apiErr, ok := v.(*APIError); ok {
			apiErr.Message = string(body)
			return nil
		}

		return fmt.Errorf("unexpected text/plain response: %s", string(body))
	})

	// Add decoder for RFC 7807 Problem Details (application/problem+json)
	client.AddContentTypeDecoder("application/problem+json", func(r io.Reader, v any) error {
		// Apple sometimes uses RFC 7807 Problem Details format
		var problemDetails struct {
			Type     string `json:"type,omitempty"`
			Title    string `json:"title,omitempty"`
			Status   int    `json:"status,omitempty"`
			Detail   string `json:"detail,omitempty"`
			Instance string `json:"instance,omitempty"`
		}

		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&problemDetails); err != nil {
			return fmt.Errorf("failed to decode problem+json: %w", err)
		}

		// Map to our APIError structure
		if apiErr, ok := v.(*APIError); ok {
			apiErr.StatusCode = problemDetails.Status
			apiErr.Message = problemDetails.Detail
			if apiErr.Message == "" {
				apiErr.Message = problemDetails.Title
			}
			apiErr.Code = problemDetails.Type
			return nil
		}

		return json.NewDecoder(strings.NewReader("")).Decode(v) // Fallback to standard JSON
	})

	logger.Debug("Configured custom content-type handlers for Apple AXM API")
}

// AddContentTypeEncoder adds a custom encoder for a specific content-type
// following Resty v3 patterns
func (c *Client) AddContentTypeEncoder(contentType string, encoder func(io.Writer, any) error) *Client {
	c.httpClient.AddContentTypeEncoder(contentType, encoder)
	c.logger.Debug("Added custom content-type encoder", zap.String("content_type", contentType))
	return c
}

// AddContentTypeDecoder adds a custom decoder for a specific content-type
// following Resty v3 patterns
func (c *Client) AddContentTypeDecoder(contentType string, decoder func(io.Reader, any) error) *Client {
	c.httpClient.AddContentTypeDecoder(contentType, decoder)
	c.logger.Debug("Added custom content-type decoder", zap.String("content_type", contentType))
	return c
}
