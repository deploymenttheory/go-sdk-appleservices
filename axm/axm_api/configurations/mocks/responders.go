package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	configurations map[string]map[string]any
}

func init() {
	mockState.configurations = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// ConfigurationsMock provides httpmock responders for configuration endpoints.
type ConfigurationsMock struct{}

// RegisterMocks registers all HTTP mock responders for configurations.
func (m *ConfigurationsMock) RegisterMocks() {
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestConfiguration()

	// GET /configurations — list configurations
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/configurations", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_configurations.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /configurations — create configuration
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/configurations", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_create_configuration_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /configurations/{id} — update configuration
	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		configID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.configurations[configID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_update_configuration_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// DELETE /configurations/{id} — delete configuration
	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		configID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.configurations[configID]
		if exists {
			delete(mockState.configurations, configID)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// GET /configurations/{id} — get configuration by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		configID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.configurations[configID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_configuration_information.json")
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
func (m *ConfigurationsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/configurations", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
	})

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/configurations", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/configurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Configuration Not Found","detail":"The requested configuration was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *ConfigurationsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.configurations {
		delete(mockState.configurations, id)
	}
}

func (m *ConfigurationsMock) seedTestConfiguration() {
	testConfig := map[string]any{
		"type": "configurations",
		"id":   "config-12345",
	}
	mockState.Lock()
	mockState.configurations["config-12345"] = testConfig
	mockState.Unlock()
}
