package axm2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"resty.dev/v3"
)

// ErrorHandler provides advanced error handling for Apple AXM API responses
type ErrorHandler struct {
	client *Client
}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler(client *Client) *ErrorHandler {
	return &ErrorHandler{client: client}
}

// HandleAPIError processes API errors with enhanced context and recovery suggestions
func (eh *ErrorHandler) HandleAPIError(resp *resty.Response, err error) error {
	if err != nil {
		return eh.handleNetworkError(err)
	}

	if resp.IsError() {
		return eh.handleHTTPError(resp)
	}

	return nil
}

// handleNetworkError processes network-level errors
func (eh *ErrorHandler) handleNetworkError(err error) error {
	return &EnhancedAPIError{
		Type:        "network_error",
		StatusCode:  0,
		Message:     fmt.Sprintf("Network error occurred: %v", err),
		OriginalErr: err,
		Suggestions: []string{
			"Check your internet connection",
			"Verify the Apple API endpoint is accessible",
			"Check if there are any firewall restrictions",
		},
	}
}

// handleHTTPError processes HTTP status code errors
func (eh *ErrorHandler) handleHTTPError(resp *resty.Response) error {
	var apiErr APIError
	
	// Try to unmarshal the error response
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil || json.Unmarshal(bodyBytes, &apiErr) != nil {
		// If unmarshaling fails, create a generic error
		return eh.createGenericError(resp)
	}

	// Enhance the error with context and suggestions
	return eh.enhanceAPIError(&apiErr, resp)
}

// enhanceAPIError adds context and suggestions to API errors
func (eh *ErrorHandler) enhanceAPIError(apiErr *APIError, resp *resty.Response) error {
	enhanced := &EnhancedAPIError{
		Type:        eh.categorizeError(resp.StatusCode()),
		StatusCode:  resp.StatusCode(),
		Message:     apiErr.Message,
		Code:        apiErr.Code,
		RequestID:   apiErr.RequestID,
		Errors:      apiErr.Errors,
		OriginalErr: apiErr,
		Suggestions: eh.getSuggestions(resp.StatusCode(), apiErr),
		Headers:     eh.extractRelevantHeaders(resp),
	}

	return enhanced
}

// createGenericError creates a generic error when response parsing fails
func (eh *ErrorHandler) createGenericError(resp *resty.Response) error {
	bodyBytes, _ := io.ReadAll(resp.Body)
	return &EnhancedAPIError{
		Type:       "unknown_error",
		StatusCode: resp.StatusCode(),
		Message:    fmt.Sprintf("API returned %d: %s", resp.StatusCode(), string(bodyBytes)),
		Suggestions: eh.getSuggestions(resp.StatusCode(), nil),
		Headers:    eh.extractRelevantHeaders(resp),
	}
}

// categorizeError categorizes errors by HTTP status code
func (eh *ErrorHandler) categorizeError(statusCode int) string {
	switch {
	case statusCode == http.StatusUnauthorized:
		return "authentication_error"
	case statusCode == http.StatusForbidden:
		return "authorization_error"
	case statusCode == http.StatusNotFound:
		return "resource_not_found"
	case statusCode == http.StatusTooManyRequests:
		return "rate_limit_error"
	case statusCode >= 400 && statusCode < 500:
		return "client_error"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown_error"
	}
}

// getSuggestions provides actionable suggestions based on error type
func (eh *ErrorHandler) getSuggestions(statusCode int, apiErr *APIError) []string {
	switch statusCode {
	case http.StatusUnauthorized:
		return []string{
			"Check if your JWT token has expired",
			"Verify your client credentials (client_id, key_id, private_key)",
			"Ensure your private key is in the correct ECDSA format",
			"Try calling ForceReauthenticate() to refresh tokens",
		}
	case http.StatusForbidden:
		return []string{
			"Verify your API scope permissions (business.api or school.api)",
			"Check if your Apple Developer account has the required entitlements",
			"Ensure your client_id is authorized for this API type",
		}
	case http.StatusNotFound:
		return []string{
			"Verify the resource ID is correct",
			"Check if the resource exists in your organization",
			"Ensure you're using the correct API endpoint",
		}
	case http.StatusTooManyRequests:
		return []string{
			"Implement exponential backoff retry logic",
			"Reduce the frequency of API calls",
			"Check the Retry-After header for recommended wait time",
		}
	case http.StatusBadRequest:
		suggestions := []string{
			"Verify your request parameters are valid",
			"Check the API documentation for required fields",
		}
		
		// Add specific suggestions based on Apple API error details
		if apiErr != nil && len(apiErr.Errors) > 0 {
			for _, err := range apiErr.Errors {
				if err.Source != nil && err.Source.Parameter != "" {
					suggestions = append(suggestions, fmt.Sprintf("Check parameter: %s", err.Source.Parameter))
				}
			}
		}
		return suggestions
	case http.StatusInternalServerError:
		return []string{
			"This is likely a temporary server issue",
			"Try again after a short delay",
			"Contact Apple Support if the issue persists",
		}
	default:
		return []string{
			"Check the Apple API status page",
			"Verify your request format matches the API documentation",
			"Contact Apple Support if needed",
		}
	}
}

// extractRelevantHeaders extracts useful headers from the response
func (eh *ErrorHandler) extractRelevantHeaders(resp *resty.Response) map[string]string {
	headers := make(map[string]string)
	
	relevantHeaders := []string{
		"X-Apple-Trace-ID",
		"X-Request-ID", 
		"Retry-After",
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
	}
	
	for _, headerName := range relevantHeaders {
		if value := resp.Header().Get(headerName); value != "" {
			headers[headerName] = value
		}
	}
	
	return headers
}

// EnhancedAPIError provides detailed error information with suggestions
type EnhancedAPIError struct {
	Type        string            `json:"type"`
	StatusCode  int               `json:"status_code"`
	Message     string            `json:"message"`
	Code        string            `json:"code,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	Errors      []AppleAPIError   `json:"errors,omitempty"`
	Suggestions []string          `json:"suggestions,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	OriginalErr error             `json:"-"`
}

// Error implements the error interface
func (e *EnhancedAPIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%s] %s (Status: %d)", e.Type, e.Message, e.StatusCode)
	}
	return fmt.Sprintf("[%s] API error (Status: %d)", e.Type, e.StatusCode)
}

// Unwrap returns the original error for error unwrapping
func (e *EnhancedAPIError) Unwrap() error {
	return e.OriginalErr
}

// IsRetryable indicates if the error suggests a retry might succeed
func (e *EnhancedAPIError) IsRetryable() bool {
	switch e.StatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, 
		 http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	case http.StatusUnauthorized:
		// Retryable if it's a token expiration issue
		return true
	default:
		return false
	}
}

// GetRetryAfter extracts the retry-after duration from headers
func (e *EnhancedAPIError) GetRetryAfter() string {
	if e.Headers != nil {
		return e.Headers["Retry-After"]
	}
	return ""
}