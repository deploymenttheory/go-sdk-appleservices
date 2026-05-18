package client

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// Common error types
var (
	ErrNoNextPage      = fmt.Errorf("no next page available")
	ErrInvalidCursor   = fmt.Errorf("invalid pagination cursor")
	ErrAuthFailed      = fmt.Errorf("authentication failed")
	ErrRateLimited     = fmt.Errorf("rate limit exceeded")
	ErrInvalidResponse = fmt.Errorf("invalid response format")
)

// APIError represents an error from the Apple Notary API.
// The Notary API returns errors as { "description": "...", "labels": [...], "name": "..." }.
type APIError struct {
	Description string
	Labels      []string
	Name        string
	StatusCode  int
}

func (e *APIError) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Name, e.Description)
	}
	if e.Description != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Description)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// ErrorResponse represents the error response structure returned by the Notary API
type ErrorResponse struct {
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Name        string   `json:"name"`
}

// ErrorHandler centralizes error handling for all API requests
type ErrorHandler struct {
	logger *zap.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError processes Notary API error responses and returns structured errors
func (eh *ErrorHandler) HandleError(resp *resty.Response, errorResp *ErrorResponse) error {
	statusCode := resp.StatusCode()

	if errorResp != nil && (errorResp.Name != "" || errorResp.Description != "") {
		if eh.logger != nil {
			eh.logger.Error("API request failed",
				zap.Int("status_code", statusCode),
				zap.String("name", errorResp.Name),
				zap.String("description", errorResp.Description),
				zap.Strings("labels", errorResp.Labels),
				zap.String("url", resp.Request.URL),
				zap.String("method", resp.Request.Method),
			)
		}

		return &APIError{
			Description: errorResp.Description,
			Labels:      errorResp.Labels,
			Name:        errorResp.Name,
			StatusCode:  statusCode,
		}
	}

	if eh.logger != nil {
		eh.logger.Error("API request failed (no structured error)",
			zap.Int("status_code", statusCode),
			zap.String("url", resp.Request.URL),
			zap.String("method", resp.Request.Method),
			zap.String("response_body", resp.String()),
		)
	}

	return &APIError{
		Name:        fmt.Sprintf("HTTP_%d", statusCode),
		Description: fmt.Sprintf("HTTP %d: %s", statusCode, http.StatusText(statusCode)),
		StatusCode:  statusCode,
	}
}
