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
	packages map[string]map[string]any
}

func init() {
	mockState.packages = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// PackagesMock provides httpmock responders for package endpoints.
type PackagesMock struct{}

// RegisterMocks registers all HTTP mock responders for packages.
func (m *PackagesMock) RegisterMocks() {
	mockState.Lock()
	mockState.packages = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestPackage()

	// GET /packages — list packages
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/packages", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_packages.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /packages/{id} — get package by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/packages/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		packageID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.packages[packageID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Package Not Found","detail":"The requested package was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_package_information.json")
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
func (m *PackagesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.packages = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/packages", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/packages/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Package Not Found","detail":"The requested package was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *PackagesMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.packages {
		delete(mockState.packages, id)
	}
}

func (m *PackagesMock) seedTestPackage() {
	testPackage := map[string]any{
		"type": "packages",
		"id":   "pkg-12345",
	}
	mockState.Lock()
	mockState.packages["pkg-12345"] = testPackage
	mockState.Unlock()
}
