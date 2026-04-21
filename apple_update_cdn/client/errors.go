package client

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// ErrorHandler centralizes error handling for all Apple Update CDN API requests.
type ErrorHandler struct {
	logger *zap.Logger
}

// NewErrorHandler creates a new error handler.
func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// HandleError processes a failed HTTP response and returns a descriptive error.
func (eh *ErrorHandler) HandleError(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	eh.logger.Error("Apple Update CDN API request failed",
		zap.Int("status_code", statusCode),
		zap.String("url", resp.Request.URL),
		zap.String("method", resp.Request.Method),
		zap.String("response_body", resp.String()),
	)

	return fmt.Errorf("HTTP %d: %s", statusCode, http.StatusText(statusCode))
}
