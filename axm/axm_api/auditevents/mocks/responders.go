package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	auditEvents []map[string]any
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// AuditEventsMock provides httpmock responders for audit event endpoints.
type AuditEventsMock struct{}

// RegisterMocks registers all HTTP mock responders for audit events.
func (m *AuditEventsMock) RegisterMocks() {
	mockState.Lock()
	mockState.auditEvents = nil
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/auditEvents", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_audit_events.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

// RegisterErrorMocks registers mock responders that return error responses.
func (m *AuditEventsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.auditEvents = nil
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/auditEvents", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *AuditEventsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.auditEvents = nil
}
