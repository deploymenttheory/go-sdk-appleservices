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
	userGroups map[string]map[string]any
}

func init() {
	mockState.userGroups = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// UserGroupsMock provides httpmock responders for user group endpoints.
type UserGroupsMock struct{}

// RegisterMocks registers all HTTP mock responders for user groups.
func (m *UserGroupsMock) RegisterMocks() {
	mockState.Lock()
	mockState.userGroups = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestUserGroup()

	// GET /userGroups — list user groups
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/userGroups", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_user_groups.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /userGroups/{id}/relationships/users — get user IDs for group
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/userGroups/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupID := parts[len(parts)-3]

		mockState.Lock()
		_, exists := mockState.userGroups[groupID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Group Not Found","detail":"The requested user group was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_user_group_user_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /userGroups/{id} — get user group by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/userGroups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.userGroups[groupID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Group Not Found","detail":"The requested user group was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_user_group_information.json")
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
func (m *UserGroupsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.userGroups = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/userGroups", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/userGroups/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Group Not Found","detail":"The requested user group was not found"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/userGroups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Group Not Found","detail":"The requested user group was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *UserGroupsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.userGroups {
		delete(mockState.userGroups, id)
	}
}

func (m *UserGroupsMock) seedTestUserGroup() {
	testGroup := map[string]any{
		"type": "userGroups",
		"id":   "UG123456",
	}
	mockState.Lock()
	mockState.userGroups["UG123456"] = testGroup
	mockState.Unlock()
}
