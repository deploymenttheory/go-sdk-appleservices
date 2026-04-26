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
	blueprints map[string]map[string]any
}

func init() {
	mockState.blueprints = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// BlueprintsMock provides httpmock responders for blueprint endpoints.
type BlueprintsMock struct{}

// RegisterMocks registers all HTTP mock responders for blueprints.
func (m *BlueprintsMock) RegisterMocks() {
	mockState.Lock()
	mockState.blueprints = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestBlueprint()

	// POST /blueprints — create blueprint
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/blueprints", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_create_blueprint_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Relationship GET endpoints — must be registered before the /{id}$ pattern to avoid overlap.

	// GET /blueprints/{id}/relationships/apps
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/apps$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_app_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /blueprints/{id}/relationships/configurations
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/configurations$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_configuration_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /blueprints/{id}/relationships/packages
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/packages$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_package_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /blueprints/{id}/relationships/orgDevices
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/orgDevices$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_device_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /blueprints/{id}/relationships/users
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_user_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /blueprints/{id}/relationships/userGroups
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/userGroups$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_user_group_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// PATCH /blueprints/{id} — update blueprint
	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-1]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_update_blueprint_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}
		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// DELETE /blueprints/{id} — delete blueprint
	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-1]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		if exists {
			delete(mockState.blueprints, blueprintID)
		}
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Relationship POST endpoints (add) — all return 204
	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/apps$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/configurations$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/packages$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/orgDevices$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/userGroups$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Relationship DELETE endpoints (remove) — all return 204
	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/apps$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/configurations$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/packages$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/orgDevices$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/userGroups$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-3]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// GET /blueprints/{id} — get blueprint by ID (registered last, least specific)
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintID := parts[len(parts)-1]
		mockState.Lock()
		_, exists := mockState.blueprints[blueprintID]
		mockState.Unlock()
		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
		}
		mockData, err := loadMockResponse("validate_get_blueprint_information.json")
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
func (m *BlueprintsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.blueprints = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/blueprints", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+/relationships/`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
	})

	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/blueprints/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Blueprint Not Found","detail":"The requested blueprint was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *BlueprintsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.blueprints {
		delete(mockState.blueprints, id)
	}
}

func (m *BlueprintsMock) seedTestBlueprint() {
	testBlueprint := map[string]any{
		"type": "blueprints",
		"id":   "blueprint-12345",
	}
	mockState.Lock()
	mockState.blueprints["blueprint-12345"] = testBlueprint
	mockState.Unlock()
}
