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
	apps map[string]map[string]any
}

func init() {
	mockState.apps = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// AppsMock provides httpmock responders for app endpoints.
type AppsMock struct{}

// RegisterMocks registers all HTTP mock responders for apps.
func (m *AppsMock) RegisterMocks() {
	mockState.Lock()
	mockState.apps = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestApp()

	// GET /apps — list apps
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/apps", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_apps.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /apps/{id} — get app by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/apps/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		appID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.apps[appID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"App Not Found","detail":"The requested app was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_app_information.json")
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
func (m *AppsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.apps = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/apps", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/apps/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"App Not Found","detail":"The requested app was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *AppsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.apps {
		delete(mockState.apps, id)
	}
}

func (m *AppsMock) seedTestApp() {
	testApp := map[string]any{
		"type": "apps",
		"id":   "361309726",
	}
	mockState.Lock()
	mockState.apps["361309726"] = testApp
	mockState.Unlock()
}
