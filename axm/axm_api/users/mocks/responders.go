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
	users map[string]map[string]any
}

func init() {
	mockState.users = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// UsersMock provides httpmock responders for user endpoints.
type UsersMock struct{}

// RegisterMocks registers all HTTP mock responders for users.
func (m *UsersMock) RegisterMocks() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestUser()

	// GET /users — list users
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/users", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_users.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /users/{id} — get user by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/users/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		userID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.users[userID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Not Found","detail":"The requested user was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_user_information.json")
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
func (m *UsersMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/users", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/users/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"User Not Found","detail":"The requested user was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *UsersMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.users {
		delete(mockState.users, id)
	}
}

func (m *UsersMock) seedTestUser() {
	testUser := map[string]any{
		"type": "users",
		"id":   "1234567890",
	}
	mockState.Lock()
	mockState.users["1234567890"] = testUser
	mockState.Unlock()
}
