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
}

func TestAPIError_Error_WithName(t *testing.T) {
	apiError := &APIError{
		StatusCode:  403,
		Name:        "FORBIDDEN",
		Description: "Authentication failure",
	}

	got := apiError.Error()
	want := "API error 403: FORBIDDEN - Authentication failure"

	if got != want {
		t.Errorf("APIError.Error() = %v, want %v", got, want)
	}
}

func TestAPIError_Error_WithoutName(t *testing.T) {
	apiError := &APIError{
		StatusCode:  404,
		Description: "The specified identifier can't be found",
	}

	got := apiError.Error()
	want := "API error 404: The specified identifier can't be found"

	if got != want {
		t.Errorf("APIError.Error() = %v, want %v", got, want)
	}
}

func TestAPIError_Error_Empty(t *testing.T) {
	apiError := &APIError{
		StatusCode: 500,
	}

	got := apiError.Error()
	if got == "" {
		t.Error("APIError.Error() returned empty string")
	}
}

func TestAPIError_Fields(t *testing.T) {
	apiError := &APIError{
		StatusCode:  403,
		Name:        "FORBIDDEN",
		Description: "Authentication failure",
		Labels:      []string{"auth", "forbidden"},
	}

	if apiError.StatusCode != 403 {
		t.Errorf("StatusCode = %v, want 403", apiError.StatusCode)
	}
	if apiError.Name != "FORBIDDEN" {
		t.Errorf("Name = %v, want FORBIDDEN", apiError.Name)
	}
	if apiError.Description != "Authentication failure" {
		t.Errorf("Description = %v, want 'Authentication failure'", apiError.Description)
	}
	if len(apiError.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(apiError.Labels))
	}
}

func TestErrorResponse_Fields(t *testing.T) {
	errorResp := &ErrorResponse{
		Name:        "FORBIDDEN",
		Description: "Authentication failure",
		Labels:      []string{"auth"},
	}

	if errorResp.Name != "FORBIDDEN" {
		t.Errorf("Name = %v, want FORBIDDEN", errorResp.Name)
	}
	if errorResp.Description != "Authentication failure" {
		t.Errorf("Description = %v, want 'Authentication failure'", errorResp.Description)
	}
	if len(errorResp.Labels) != 1 {
		t.Errorf("len(Labels) = %d, want 1", len(errorResp.Labels))
	}
}

func TestCommonErrors(t *testing.T) {
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

	messages := make(map[string]bool)
	for _, err := range errors {
		msg := err.Error()
		if messages[msg] {
			t.Errorf("Duplicate error message: %s", msg)
		}
		messages[msg] = true
	}
}
