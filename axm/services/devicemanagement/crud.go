package devicemanagement

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/interfaces"
)

type (
	// DeviceManagementServiceInterface defines the interface for device management operations
	DeviceManagementServiceInterface interface {
		// GetDeviceManagementServices retrieves a list of device management services (MDM servers) in an organization
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-mdm-servers
		GetDeviceManagementServices(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, error)

		// GetMDMServerDeviceLinkages retrieves a list of device IDs assigned to an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-mdmserver
		GetMDMServerDeviceLinkages(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, error)

		// GetAssignedDeviceManagementServiceIDForADevice retrieves the assigned device management service ID linkage for a device
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-server-id-for-an-orgdevice
		GetAssignedDeviceManagementServiceIDForADevice(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, error)

		// GetAssignedDeviceManagementServiceInformationByDeviceID retrieves the assigned device management service information for a device
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-server-information-for-an-orgdevice
		GetAssignedDeviceManagementServiceInformationByDeviceID(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, error)

		// AssignDevicesToServer assigns devices to an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
		AssignDevicesToServer(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error)

		// UnassignDevicesFromServer unassigns devices from an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
		UnassignDevicesFromServer(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error)
	}

	// DeviceManagementService handles communication with the device management
	// related methods of the Apple Business Manager API.
	//
	// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessmanagerapi/
	DeviceManagementService struct {
		client interfaces.HTTPClient
	}
)

var _ DeviceManagementServiceInterface = (*DeviceManagementService)(nil)

// NewService creates a new device management service
func NewService(client interfaces.HTTPClient) *DeviceManagementService {
	return &DeviceManagementService{
		client: client,
	}
}

// GetDeviceManagementServices retrieves a list of device management services (MDM servers) in an organization
// URL: GET https://api-business.apple.com/v1/mdmServers
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-mdm-servers
func (s *DeviceManagementService) GetDeviceManagementServices(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := "/mdmServers"

	queryParams := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		queryParams.AddInt("limit", opts.Limit)
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ResponseMDMServers
	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMDMServerDeviceLinkages retrieves a list of device IDs assigned to an MDM server
// URL: GET https://api-business.apple.com/v1/mdmServers/{id}/relationships/devices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-mdmserver
func (s *DeviceManagementService) GetMDMServerDeviceLinkages(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}

	endpoint := fmt.Sprintf("/mdmServers/%s/relationships/devices", mdmServerID)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	queryParams := s.client.QueryBuilder()

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		queryParams.AddInt("limit", opts.Limit)
	}

	var result ResponseMDMServerDevicesLinkages

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAssignedDeviceManagementServiceIDForADevice retrieves the assigned device management service ID linkage for a device
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/relationships/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-server-id-for-an-orgdevice
func (s *DeviceManagementService) GetAssignedDeviceManagementServiceIDForADevice(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf("/orgDevices/%s/relationships/assignedServer", deviceID)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ResponseOrgDeviceAssignedServerLinkage
	err := s.client.Get(ctx, endpoint, nil, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAssignedDeviceManagementServiceInformationByDeviceID retrieves the assigned device management service information for a device
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-server-information-for-an-orgdevice
func (s *DeviceManagementService) GetAssignedDeviceManagementServiceInformationByDeviceID(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf("/orgDevices/%s/assignedServer", deviceID)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	queryParams := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	var result MDMServerResponse
	err := s.client.Get(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// AssignDevicesToServer assigns devices to an MDM server
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagementService) AssignDevicesToServer(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("at least one device ID is required")
	}

	endpoint := "/orgDeviceActivities"

	deviceLinkages := make([]OrgDeviceActivityDeviceLinkage, len(deviceIDs))
	for i, deviceID := range deviceIDs {
		deviceLinkages[i] = OrgDeviceActivityDeviceLinkage{
			Type: "orgDevices",
			ID:   deviceID,
		}
	}

	request := &OrgDeviceActivityCreateRequest{
		Data: OrgDeviceActivityData{
			Type: "orgDeviceActivities",
			Attributes: OrgDeviceActivityCreateAttributes{
				ActivityType: ActivityTypeAssignDevices,
			},
			Relationships: OrgDeviceActivityCreateRelationships{
				MDMServer: &OrgDeviceActivityMDMServerRelationship{
					Data: OrgDeviceActivityMDMServerLinkage{
						Type: "mdmServers",
						ID:   mdmServerID,
					},
				},
				Devices: &OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages,
				},
			},
		},
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ResponseOrgDeviceActivity

	err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UnassignDevicesFromServer unassigns devices from an MDM server
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagementService) UnassignDevicesFromServer(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("at least one device ID is required")
	}

	endpoint := "/orgDeviceActivities"

	deviceLinkages := make([]OrgDeviceActivityDeviceLinkage, len(deviceIDs))
	for i, deviceID := range deviceIDs {
		deviceLinkages[i] = OrgDeviceActivityDeviceLinkage{
			Type: "orgDevices",
			ID:   deviceID,
		}
	}

	request := &OrgDeviceActivityCreateRequest{
		Data: OrgDeviceActivityData{
			Type: "orgDeviceActivities",
			Attributes: OrgDeviceActivityCreateAttributes{
				ActivityType: ActivityTypeUnassignDevices,
			},
			Relationships: OrgDeviceActivityCreateRelationships{
				MDMServer: &OrgDeviceActivityMDMServerRelationship{
					Data: OrgDeviceActivityMDMServerLinkage{
						Type: "mdmServers",
						ID:   mdmServerID,
					},
				},
				Devices: &OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages,
				},
			},
		},
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ResponseOrgDeviceActivity
	err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
