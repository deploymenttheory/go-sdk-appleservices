package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":"Resource Not Found"}`))
}

// loadMockResponse loads JSON response from the mocks folder
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

type MSAppsMock struct{}

// RegisterMocks registers all the HTTP mock responders for Microsoft Mac Apps API
func (m *MSAppsMock) RegisterMocks() {
	// GET /latest - Get latest application versions
	httpmock.RegisterResponder("GET", "https://appledevicepolicy.tools/api/latest", func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_latest_apps.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":"Failed to load mock data"}`), nil
		}

		// Parse and return as proper JSON
		var jsonData map[string]any
		if err := json.Unmarshal(mockData, &jsonData); err != nil {
			return httpmock.NewStringResponse(500, `{"error":"Failed to parse mock data"}`), nil
		}

		return httpmock.NewJsonResponse(200, jsonData)
	})
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *MSAppsMock) RegisterErrorMocks() {
	// GET /latest - Return error
	httpmock.RegisterResponder("GET", "https://appledevicepolicy.tools/api/latest", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"error":"Internal Server Error"}`), nil
	})
}

// CleanupMockState clears all mock state data
func (m *MSAppsMock) CleanupMockState() {
	// No state to cleanup for this API
}
