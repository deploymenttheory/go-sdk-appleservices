package client

import (
	"encoding/json"
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

// APIError represents a single error from the Apple Business Manager API
type APIError struct {
	ID     string          `json:"id,omitempty"`
	Status string          `json:"status"`
	Code   string          `json:"code"`
	Title  string          `json:"title"`
	Detail string          `json:"detail"`
	Source *APIErrorSource `json:"source,omitempty"`
	Links  *ErrorLinks     `json:"links,omitempty"`
	Meta   *APIErrorMeta   `json:"meta,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error %s: %s - %s", e.Status, e.Code, e.Detail)
	}
	return fmt.Sprintf("API error %s: %s", e.Status, e.Detail)
}

// APIErrorSource represents the source of an error (JsonPointer or Parameter)
type APIErrorSource struct {
	JsonPointer *JsonPointer `json:"jsonPointer,omitempty"`
	Parameter   *Parameter   `json:"parameter,omitempty"`
}

// JsonPointer represents a JSON pointer source
type JsonPointer struct {
	Pointer string `json:"pointer"`
}

// Parameter represents a query parameter source
type Parameter struct {
	Parameter string `json:"parameter"`
}

// ErrorLinks contains error-related links
type ErrorLinks struct {
	About      string                `json:"about,omitempty"`
	Associated *ErrorLinksAssociated `json:"associated,omitempty"`
}

// ErrorLinksAssociated represents associated error links
type ErrorLinksAssociated struct {
	Href string                    `json:"href"`
	Meta *ErrorLinksAssociatedMeta `json:"meta,omitempty"`
}

// ErrorLinksAssociatedMeta contains metadata for associated error links
type ErrorLinksAssociatedMeta struct {
	// Can contain any key-value pairs as specified in the API
	AdditionalProperties map[string]any `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling for ErrorLinksAssociatedMeta
func (m *ErrorLinksAssociatedMeta) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.AdditionalProperties)
}

// MarshalJSON implements custom marshaling for ErrorLinksAssociatedMeta
func (m *ErrorLinksAssociatedMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.AdditionalProperties)
}

// APIErrorMeta contains additional error metadata
type APIErrorMeta struct {
	// Can contain any key-value pairs as specified in the API documentation
	AdditionalProperties map[string]any `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling for APIErrorMeta
func (m *APIErrorMeta) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.AdditionalProperties)
}

// MarshalJSON implements custom marshaling for APIErrorMeta
func (m *APIErrorMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.AdditionalProperties)
}

// ErrorResponse represents the complete error response structure returned by the API
type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

// APIErrorResponse represents the legacy error structure (keeping for backward compatibility)
type APIErrorResponse struct {
	ErrorCode string         `json:"error_code"`
	Message   string         `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
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

// HandleError processes API error responses and returns structured errors
func (eh *ErrorHandler) HandleError(resp *resty.Response, errorResp *ErrorResponse) error {
	statusCode := resp.StatusCode()

	if len(errorResp.Errors) > 0 {
		for i, apiError := range errorResp.Errors {
			logFields := []zap.Field{
				zap.Int("error_index", i),
				zap.String("error_id", apiError.ID),
				zap.String("status", apiError.Status),
				zap.String("code", apiError.Code),
				zap.String("title", apiError.Title),
				zap.String("detail", apiError.Detail),
				zap.String("url", resp.Request.URL),
				zap.String("method", resp.Request.Method),
			}

			if apiError.Source != nil {
				if apiError.Source.JsonPointer != nil {
					logFields = append(logFields, zap.String("source_json_pointer", apiError.Source.JsonPointer.Pointer))
				}
				if apiError.Source.Parameter != nil {
					logFields = append(logFields, zap.String("source_parameter", apiError.Source.Parameter.Parameter))
				}
			}

			if apiError.Links != nil {
				if apiError.Links.About != "" {
					logFields = append(logFields, zap.String("links_about", apiError.Links.About))
				}
				if apiError.Links.Associated != nil {
					logFields = append(logFields, zap.String("links_associated_href", apiError.Links.Associated.Href))
					if apiError.Links.Associated.Meta != nil && apiError.Links.Associated.Meta.AdditionalProperties != nil {
						logFields = append(logFields, zap.Any("links_associated_meta", apiError.Links.Associated.Meta.AdditionalProperties))
					}
				}
			}

			if apiError.Meta != nil && apiError.Meta.AdditionalProperties != nil {
				logFields = append(logFields, zap.Any("error_meta", apiError.Meta.AdditionalProperties))
			}

			eh.logger.Error("API request failed", logFields...)
		}

		firstError := errorResp.Errors[0]
		return &firstError
	}

	eh.logger.Error("API request failed (no structured error)",
		zap.Int("status_code", statusCode),
		zap.String("url", resp.Request.URL),
		zap.String("method", resp.Request.Method),
		zap.String("response_body", resp.String()),
	)

	return &APIError{
		Status: fmt.Sprintf("%d", statusCode),
		Code:   fmt.Sprintf("HTTP_%d", statusCode),
		Title:  http.StatusText(statusCode),
		Detail: fmt.Sprintf("HTTP %d: %s", statusCode, http.StatusText(statusCode)),
	}
}
