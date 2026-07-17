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
	organizationalUnits map[string]map[string]any
}

func init() {
	mockState.organizationalUnits = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// OrganizationalUnitsMock provides httpmock responders for organizational unit endpoints.
type OrganizationalUnitsMock struct{}

// RegisterMocks registers all HTTP mock responders for organizational units.
func (m *OrganizationalUnitsMock) RegisterMocks() {
	mockState.Lock()
	mockState.organizationalUnits = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestOrganizationalUnit()

	// GET /organizationalUnits — list organizational units
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/organizationalUnits", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_organizational_units.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /organizationalUnits/{id}/relationships/users — get user IDs for organizational unit
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/organizationalUnits/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		unitID := parts[len(parts)-3]

		mockState.Lock()
		_, exists := mockState.organizationalUnits[unitID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Organizational Unit Not Found","detail":"The requested organizational unit was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_organizational_unit_user_ids.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /organizationalUnits/{id} — get organizational unit by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/organizationalUnits/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		unitID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.organizationalUnits[unitID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Organizational Unit Not Found","detail":"The requested organizational unit was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_organizational_unit_information.json")
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
func (m *OrganizationalUnitsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.organizationalUnits = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/organizationalUnits", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/organizationalUnits/[^/]+/relationships/users$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Organizational Unit Not Found","detail":"The requested organizational unit was not found"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/organizationalUnits/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Organizational Unit Not Found","detail":"The requested organizational unit was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *OrganizationalUnitsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.organizationalUnits {
		delete(mockState.organizationalUnits, id)
	}
}

func (m *OrganizationalUnitsMock) seedTestOrganizationalUnit() {
	testUnit := map[string]any{
		"type": "organizationalUnits",
		"id":   "OU789012",
	}
	mockState.Lock()
	mockState.organizationalUnits["OU789012"] = testUnit
	mockState.Unlock()
}
