package client

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewErrorHandler(t *testing.T) {
	logger := zap.NewNop()
	handler := NewErrorHandler(logger)

	if handler == nil {
		t.Fatal("NewErrorHandler returned nil")
	}

	if handler.logger == nil {
		t.Error("ErrorHandler logger is nil")
	}
}

func TestNewErrorHandler_NilLogger(t *testing.T) {
	handler := NewErrorHandler(nil)

	if handler == nil {
		t.Fatal("NewErrorHandler(nil) returned nil")
	}

	// Should not panic with nil logger
	if handler.logger != nil {
		t.Error("Expected nil logger")
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiError *APIError
		wantText string
	}{
		{
			name: "Full error with code",
			apiError: &APIError{
				Status: "400",
				Code:   "INVALID_REQUEST",
				Title:  "Invalid Request",
				Detail: "The request parameters are invalid",
			},
			wantText: "API error 400: INVALID_REQUEST - The request parameters are invalid",
		},
		{
			name: "Error without code",
			apiError: &APIError{
				Status: "404",
				Title:  "Not Found",
				Detail: "Resource not found",
			},
			wantText: "API error 404: Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiError.Error()
			if got != tt.wantText {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.wantText)
			}
		})
	}
}

func TestAPIError_ErrorFields(t *testing.T) {
	apiError := &APIError{
		ID:     "error-123",
		Status: "400",
		Code:   "INVALID_PARAM",
		Title:  "Invalid Parameter",
		Detail: "The parameter 'name' is invalid",
		Source: &APIErrorSource{
			Parameter: &Parameter{
				Parameter: "name",
			},
		},
	}

	if apiError.ID != "error-123" {
		t.Errorf("ID = %v, want 'error-123'", apiError.ID)
	}

	if apiError.Status != "400" {
		t.Errorf("Status = %v, want '400'", apiError.Status)
	}

	if apiError.Source == nil {
		t.Fatal("Source is nil")
	}

	if apiError.Source.Parameter == nil {
		t.Fatal("Source.Parameter is nil")
	}

	if apiError.Source.Parameter.Parameter != "name" {
		t.Errorf("Source.Parameter.Parameter = %v, want 'name'", apiError.Source.Parameter.Parameter)
	}
}

func TestAPIErrorWithJsonPointer(t *testing.T) {
	apiError := &APIError{
		Status: "400",
		Code:   "INVALID_FIELD",
		Detail: "Invalid JSON field",
		Source: &APIErrorSource{
			JsonPointer: &JsonPointer{
				Pointer: "/data/attributes/name",
			},
		},
	}

	if apiError.Source == nil || apiError.Source.JsonPointer == nil {
		t.Fatal("JsonPointer source is nil")
	}

	if apiError.Source.JsonPointer.Pointer != "/data/attributes/name" {
		t.Errorf("JsonPointer.Pointer = %v, want '/data/attributes/name'", apiError.Source.JsonPointer.Pointer)
	}
}

func TestAPIErrorWithLinks(t *testing.T) {
	apiError := &APIError{
		Status: "400",
		Detail: "Error with links",
		Links: &ErrorLinks{
			About: "https://docs.apple.com/error/400",
			Associated: &ErrorLinksAssociated{
				Href: "https://api.apple.com/resource",
			},
		},
	}

	if apiError.Links == nil {
		t.Fatal("Links is nil")
	}

	if apiError.Links.About != "https://docs.apple.com/error/400" {
		t.Errorf("Links.About = %v", apiError.Links.About)
	}

	if apiError.Links.Associated == nil {
		t.Fatal("Links.Associated is nil")
	}

	if apiError.Links.Associated.Href != "https://api.apple.com/resource" {
		t.Errorf("Links.Associated.Href = %v", apiError.Links.Associated.Href)
	}
}

func TestErrorResponse(t *testing.T) {
	errorResp := &ErrorResponse{
		Errors: []APIError{
			{
				Status: "400",
				Code:   "ERROR_1",
				Detail: "First error",
			},
			{
				Status: "400",
				Code:   "ERROR_2",
				Detail: "Second error",
			},
		},
	}

	if len(errorResp.Errors) != 2 {
		t.Errorf("len(Errors) = %d, want 2", len(errorResp.Errors))
	}

	if errorResp.Errors[0].Code != "ERROR_1" {
		t.Errorf("Errors[0].Code = %v, want 'ERROR_1'", errorResp.Errors[0].Code)
	}

	if errorResp.Errors[1].Code != "ERROR_2" {
		t.Errorf("Errors[1].Code = %v, want 'ERROR_2'", errorResp.Errors[1].Code)
	}
}

func TestAPIErrorResponse_LegacyFormat(t *testing.T) {
	legacyError := &APIErrorResponse{
		ErrorCode: "LEGACY_ERROR",
		Message:   "This is a legacy error",
		Details: map[string]any{
			"field": "name",
			"value": "invalid",
		},
	}

	if legacyError.ErrorCode != "LEGACY_ERROR" {
		t.Errorf("ErrorCode = %v, want 'LEGACY_ERROR'", legacyError.ErrorCode)
	}

	if legacyError.Message != "This is a legacy error" {
		t.Errorf("Message = %v", legacyError.Message)
	}

	if len(legacyError.Details) != 2 {
		t.Errorf("len(Details) = %d, want 2", len(legacyError.Details))
	}
}

func TestCommonErrors(t *testing.T) {
	// Test common error variables
	if ErrNoNextPage == nil {
		t.Error("ErrNoNextPage is nil")
	}

	if ErrInvalidCursor == nil {
		t.Error("ErrInvalidCursor is nil")
	}

	if ErrAuthFailed == nil {
		t.Error("ErrAuthFailed is nil")
	}

	if ErrRateLimited == nil {
		t.Error("ErrRateLimited is nil")
	}

	if ErrInvalidResponse == nil {
		t.Error("ErrInvalidResponse is nil")
	}
}

func TestErrorConstants_UniqueMessages(t *testing.T) {
	errors := []error{
		ErrNoNextPage,
		ErrInvalidCursor,
		ErrAuthFailed,
		ErrRateLimited,
		ErrInvalidResponse,
	}

	// Check that all error messages are unique
	messages := make(map[string]bool)
	for _, err := range errors {
		msg := err.Error()
		if messages[msg] {
			t.Errorf("Duplicate error message: %s", msg)
		}
		messages[msg] = true
	}
}

func TestAPIErrorSource_BothTypes(t *testing.T) {
	// Test source with both JsonPointer and Parameter (edge case)
	source := &APIErrorSource{
		JsonPointer: &JsonPointer{
			Pointer: "/data/attributes/name",
		},
		Parameter: &Parameter{
			Parameter: "name",
		},
	}

	if source.JsonPointer == nil {
		t.Error("JsonPointer is nil")
	}

	if source.Parameter == nil {
		t.Error("Parameter is nil")
	}
}
