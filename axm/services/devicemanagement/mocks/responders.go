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
	mdmServers       map[string]map[string]any
	deviceActivities map[string]map[string]any
}

func init() {
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

type DeviceManagementMock struct{}

// RegisterMocks registers all the HTTP mock responders for device management services
func (m *DeviceManagementMock) RegisterMocks() {
	mockState.Lock()
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	mockState.Unlock()

	// Seed with default test MDM server
	m.seedTestMDMServer()

	// GET /mdmServers - List device management services
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_device_management_services.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// If we have MDM servers in mock state, return them
		if len(mockState.mdmServers) > 0 {
			serverList := make([]map[string]any, 0, len(mockState.mdmServers))
			for _, server := range mockState.mdmServers {
				serverCopy := make(map[string]any)
				for k, v := range server {
					serverCopy[k] = v
				}
				serverList = append(serverList, serverCopy)
			}
			responseObj["data"] = serverList
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[mdmServers]"); fields != "" {
			// In a real implementation, we would filter fields here
		}

		// Handle limit parameter
		if limit := query.Get("limit"); limit != "" {
			// In a real implementation, we would apply pagination here
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /mdmServers/{id}/relationships/devices - Get MDM server device linkages
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices", func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_mdm_server_device_linkages.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id}/relationships/assignedServer - Get assigned server linkage
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/relationships/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_assigned_server_linkage.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id}/assignedServer - Get assigned server information
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_assigned_server_info.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[mdmServers]"); fields != "" {
			// In a real implementation, we would filter fields here
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /orgDeviceActivities - Assign/Unassign devices
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/orgDeviceActivities", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid request body"}]}`), nil
		}

		// Extract activity type from request
		data, ok := requestBody["data"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid request structure"}]}`), nil
		}

		attributes, ok := data["attributes"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Missing attributes"}]}`), nil
		}

		activityType, ok := attributes["activityType"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Missing activityType"}]}`), nil
		}

		// Choose appropriate response based on activity type
		var mockFile string
		if activityType == "ASSIGN_DEVICES" {
			mockFile = "validate_assign_devices_response.json"
		} else if activityType == "UNASSIGN_DEVICES" {
			mockFile = "validate_unassign_devices_response.json"
		} else {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid activityType"}]}`), nil
		}

		// Load the mock response template
		mockData, err := loadMockResponse(mockFile)
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(201, responseObj)
	})
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *DeviceManagementMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	mockState.Unlock()

	// GET /mdmServers - Return error
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	// GET /mdmServers/{id}/relationships/devices - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/mdmServers/.+/relationships/devices$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
	})

	// GET /orgDevices/{id}/relationships/assignedServer - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/relationships/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})

	// GET /orgDevices/{id}/assignedServer - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})

	// POST /orgDeviceActivities - Return error
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/orgDeviceActivities", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Mock error for testing"}]}`), nil
	})
}

// CleanupMockState clears all mock state data
func (m *DeviceManagementMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.mdmServers {
		delete(mockState.mdmServers, id)
	}
	for id := range mockState.deviceActivities {
		delete(mockState.deviceActivities, id)
	}
}

// AddMockMDMServer adds an MDM server to the mock state
func (m *DeviceManagementMock) AddMockMDMServer(serverID string, server map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.mdmServers[serverID] = server
}

// GetMockMDMServer retrieves an MDM server from the mock state
func (m *DeviceManagementMock) GetMockMDMServer(serverID string) (map[string]any, bool) {
	mockState.Lock()
	defer mockState.Unlock()
	server, exists := mockState.mdmServers[serverID]
	return server, exists
}

// seedTestMDMServer adds a default test MDM server to the mock state
func (m *DeviceManagementMock) seedTestMDMServer() {
	testServer := map[string]any{
		"type": "mdmServers",
		"id":   "1F97349736CF4614A94F624E705841AD",
		"attributes": map[string]any{
			"serverName":      "Test Device Management Service",
			"serverType":      "MDM",
			"createdDateTime": "2025-05-01T03:21:44.685Z",
			"updatedDateTime": "2025-05-01T03:21:46.284Z",
		},
		"relationships": map[string]any{
			"devices": map[string]any{
				"links": map[string]any{
					"self": "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices",
				},
			},
		},
	}

	m.AddMockMDMServer("1F97349736CF4614A94F624E705841AD", testServer)
}
