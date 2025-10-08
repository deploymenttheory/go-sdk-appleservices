package axm2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/client/shared"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/mdm_servers"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/org_device_activities"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/org_devices"
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

	// Service instances for hierarchical API access
	orgDevicesService          *org_devices.Service
	mdmServersService          *mdm_servers.Service
	orgDeviceActivitiesService *org_device_activities.Service
}

// Ensure Client implements AXMClient interface at compile time
var _ AXMClient = (*Client)(nil)

// APIError represents structured API error responses following Resty v3 patterns
type APIError struct {
	StatusCode int             `json:"status_code,omitempty"`
	Message    string          `json:"message,omitempty"`
	Code       string          `json:"code,omitempty"`
	RequestID  string          `json:"request_id,omitempty"`
	Errors     []AppleAPIError `json:"errors,omitempty"` // Apple's error format
}

// AppleAPIError represents Apple's specific error format based on official API spec
type AppleAPIError struct {
	ID     string       `json:"id,omitempty"`     // The unique ID of a specific instance of an error
	Status string       `json:"status"`           // The HTTP status code of the error
	Code   string       `json:"code"`             // A machine-readable code indicating the type of error
	Title  string       `json:"title"`            // A summary of the error
	Detail string       `json:"detail"`           // A detailed explanation of the error
	Source *ErrorSource `json:"source,omitempty"` // Source of the error (parameter or JSON pointer)
	Links  *ErrorLinks  `json:"links,omitempty"`  // Links related to the error
	Meta   any          `json:"meta,omitempty"`   // Additional metadata
}

// ErrorSource represents the source of an error (parameter or JSON pointer)
type ErrorSource struct {
	Parameter string `json:"parameter,omitempty"` // Query parameter that produces the error
	Pointer   string `json:"pointer,omitempty"`   // JSON pointer indicating the location of the error
}

// ErrorLinks represents links related to error responses
type ErrorLinks struct {
	Self string `json:"self,omitempty"`
}

// ErrorResponse represents the complete error response structure from Apple API
type ErrorResponse struct {
	Errors []AppleAPIError `json:"errors"` // An array of one or more errors
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

// Enhanced Config holds configuration for the AXM client with Resty v3 enhancements
type Config struct {
	BaseURL                 string        // Apple School/Business Manager API base URL (optional, will be set based on APIType)
	APIType                 string        // API type: "abm" for Apple Business Manager, "asm" for Apple School Manager
	ClientID                string        // Client ID from Apple
	KeyID                   string        // Key ID from Apple
	PrivateKey              string        // Private key content (PEM format)
	Scope                   string        // OAuth scope ("business.api" or "school.api")
	Timeout                 time.Duration // HTTP timeout (default: 30s)
	RetryCount              int           // Number of retries (default: 3)
	RetryMinWait            time.Duration // Minimum wait time between retries (default: 1s)
	RetryMaxWait            time.Duration // Maximum wait time between retries (default: 10s)
	UserAgent               string        // User agent string
	Debug                   bool          // Enable debug logging
	EnableRetryLog          bool          // Enable detailed retry logging (default: true)
	DebugLogFormat          string        // Debug log format: "human" or "json" (default: "human")
	DebugLogBodyLimit       int64         // Debug log body size limit in bytes (default: 1MB)
	CircuitBreakerEnabled   bool          // Enable circuit breaker (default: false)
	CircuitBreakerThreshold int           // Circuit breaker failure threshold (default: 5)
	BackupEndpoints         []string      // Backup API endpoints for load balancing
}

// NewClient creates a new AXM API client
func NewClient(config Config) (AXMClient, error) {
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

	// Add response middleware for Apple-specific processing
	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		// Log Apple-specific response headers for debugging
		if appleTraceID := resp.Header().Get("X-Apple-Trace-ID"); appleTraceID != "" {
			logger.Debug("Apple trace ID", zap.String("trace_id", appleTraceID))
		}

		// Handle Apple-specific rate limiting headers
		if retryAfter := resp.Header().Get("Retry-After"); retryAfter != "" {
			logger.Info("Apple API rate limit hit", zap.String("retry_after", retryAfter))
		}

		// Log request tracking information
		if requestID := resp.Header().Get("X-Request-ID"); requestID != "" {
			logger.Debug("Apple request ID", zap.String("request_id", requestID))
		}

		return nil
	})

	if config.Debug {
		httpClient.SetDebug(true)
		// Note: Debug log formatting will be available in future Resty v3 releases
		// For now, we use the default debug format
		logger.Debug("Debug logging enabled", zap.String("format", config.DebugLogFormat))
	}

	// Add custom content-type encoders and decoders for Apple AXM API
	setupContentTypeHandlers(httpClient, logger)

	client.httpClient = httpClient

	// Initialize service instances
	client.orgDevicesService = org_devices.NewService(client, logger)
	client.mdmServersService = mdm_servers.NewService(client, logger)
	client.orgDeviceActivitiesService = org_device_activities.NewService(client, logger)

	return client, nil
}

// HealthCheck performs a simple connectivity and authentication check
func (c *Client) HealthCheck(ctx context.Context) error {
	c.logger.Debug("Performing health check")

	// Simple HEAD request to base URL to verify connectivity
	// Apple doesn't provide a specific health endpoint, so we'll test auth
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetTimeout(10 * time.Second).
		Head(c.config.BaseURL)

	if err != nil {
		return fmt.Errorf("health check failed - network error: %w", err)
	}

	// Check for common error status codes
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("health check failed - API returned %d", resp.StatusCode())
	}

	c.logger.Debug("Health check passed", zap.Int("status_code", resp.StatusCode()))
	return nil
}

// GetDiagnostics returns diagnostic information about the client
func (c *Client) GetDiagnostics() map[string]any {
	return map[string]any{
		"base_url":         c.config.BaseURL,
		"api_type":         c.config.APIType,
		"retry_count":      c.config.RetryCount,
		"timeout":          c.config.Timeout.String(),
		"is_authenticated": c.IsAuthenticated(),
		"client_id":        c.config.ClientID,
		"scope":            c.config.Scope,
		"debug_enabled":    c.config.Debug,
		"user_agent":       c.config.UserAgent,
	}
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

// Get performs a GET request with automatic error handling
func (c *Client) Get(ctx context.Context, endpoint string, result any, opts ...any) error {
	var apiError APIError
	
	request := c.httpClient.R().
		SetContext(ctx).
		SetResult(result).
		SetError(&apiError)
	
	// Apply options
	c.ApplyRequestOptions(&RequestWrapper{req: request, parent: c}, opts...)
	
	response, err := request.Get(endpoint)
	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}
	
	if response.IsError() {
		return c.handleAPIError(response, &apiError)
	}
	
	return nil
}

// Post performs a POST request with automatic error handling
func (c *Client) Post(ctx context.Context, endpoint string, body, result any, opts ...any) error {
	var apiError APIError
	
	request := c.httpClient.R().
		SetContext(ctx).
		SetError(&apiError)
	
	if result != nil {
		request.SetResult(result)
	}
	
	if body != nil {
		request.SetBody(body)
	}
	
	// Apply options
	c.ApplyRequestOptions(&RequestWrapper{req: request, parent: c}, opts...)
	
	response, err := request.Post(endpoint)
	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}
	
	if response.IsError() {
		return c.handleAPIError(response, &apiError)
	}
	
	return nil
}

// GetWithPagination performs a paginated GET request
func (c *Client) GetWithPagination(ctx context.Context, endpoint string, newResponseFunc func() shared.PaginatedResponse, opts ...any) (any, error) {
	return c.DoRequestWithPagination(ctx, endpoint, func() shared.PaginatedResponse {
		return newResponseFunc()
	}, opts...)
}

// handleAPIError processes API errors consistently
func (c *Client) handleAPIError(response *resty.Response, apiError *APIError) error {
	c.logger.Error("API request failed",
		zap.String("method", response.Request.Method),
		zap.String("url", response.Request.URL),
		zap.Int("status_code", response.StatusCode()),
		zap.Any("error", apiError))
	return fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
}


// DoRequestWithPagination performs a paginated GET request using Resty v3 patterns
func (c *Client) DoRequestWithPagination(ctx context.Context, endpoint string, newResponseFunc func() shared.PaginatedResponse, opts ...any) (any, error) {
	var allData any
	nextURL := endpoint
	pageCount := 0

	c.logger.Debug("Starting paginated request", zap.String("endpoint", endpoint))

	for nextURL != "" {
		pageCount++
		c.logger.Debug("Fetching page", zap.Int("page", pageCount), zap.String("url", nextURL))

		// Create new response instance for this page
		pageResponse := newResponseFunc()
		var apiError APIError

		request := &RequestWrapper{req: c.httpClient.R(), parent: c}
		request.SetContext(ctx).
			SetResult(pageResponse).
			SetError(&apiError)

		// Apply RequestOption parameters for first page only
		if pageCount == 1 && len(opts) > 0 {
			c.ApplyRequestOptions(request, opts...)
		}

		response, err := request.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to execute paginated GET request (page %d): %w", pageCount, err)
		}

		if response.IsError() {
			c.logger.Error("API error in paginated request",
				zap.Int("page", pageCount),
				zap.Int("status_code", response.StatusCode()),
				zap.Any("error", apiError))
			return nil, fmt.Errorf("API error (page %d): %d %s", pageCount, response.StatusCode(), response.String())
		}

		// For first page, initialize allData with a copy
		if pageCount == 1 {
			allData = pageResponse.GetData()
		} else {
			// Append current page data to accumulated data
			if currentPageData := pageResponse.GetData(); currentPageData != nil {
				switch existingData := allData.(type) {
				case []OrgDevice:
					if newData, ok := currentPageData.([]OrgDevice); ok {
						allData = append(existingData, newData...)
					}
				case []MdmServer:
					if newData, ok := currentPageData.([]MdmServer); ok {
						allData = append(existingData, newData...)
					}
				case []string: // For device IDs
					if newData, ok := currentPageData.([]string); ok {
						allData = append(existingData, newData...)
					}
				}
			}
		}

		// Get next URL for pagination
		nextURL = pageResponse.GetNextURL()

		c.logger.Debug("Page fetched successfully",
			zap.Int("page", pageCount),
			zap.String("next_url", nextURL))
	}

	c.logger.Info("Successfully completed paginated request",
		zap.String("endpoint", endpoint),
		zap.Int("total_pages", pageCount))

	return allData, nil
}

// DoRequest performs a generic HTTP request using Resty v3 patterns
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body any, result any) (*resty.Response, error) {
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

// Service accessor methods for hierarchical API access

// OrgDevices returns the organization devices service
func (c *Client) OrgDevices() OrgDevicesService {
	return c.orgDevicesService
}

// MdmServers returns the MDM servers service
func (c *Client) MdmServers() MdmServersService {
	return c.mdmServersService
}

// OrgDeviceActivities returns the organization device activities service
func (c *Client) OrgDeviceActivities() OrgDeviceActivitiesService {
	return c.orgDeviceActivitiesService
}

// ApplyRequestOptions applies client RequestOptions to a service request
func (c *Client) ApplyRequestOptions(req shared.RequestInterface, opts ...any) {
	for _, opt := range opts {
		if requestOption, ok := opt.(RequestOption); ok {
			rb := &RequestBuilder{req: req}
			requestOption(rb)
		}
	}
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
