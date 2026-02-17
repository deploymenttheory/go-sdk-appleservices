package devicemanagement

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/interfaces"
)

type (
	// DeviceManagementServiceInterface defines the interface for device management operations
	DeviceManagementServiceInterface interface {
		// GetDeviceManagementServicesV1 retrieves a list of device management services (MDM servers) in an organization
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-device-management-services
		GetDeviceManagementServicesV1(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, error)

		// GetDeviceSerialNumbersForDeviceManagementServiceV1 retrieves a list of device IDs assigned to an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-device-management-service
		GetDeviceSerialNumbersForDeviceManagementServiceV1(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, error)

		// GetAssignedDeviceManagementServiceIDForADeviceV1 retrieves the assigned device management service ID linkage for a device
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-id-for-an-orgdevice
		GetAssignedDeviceManagementServiceIDForADeviceV1(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, error)

		// GetAssignedDeviceManagementServiceInformationByDeviceIDV1 retrieves the assigned device management service information for a device
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-information-for-an-orgdevice
		GetAssignedDeviceManagementServiceInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, error)

		// AssignDevicesToServerV1 assigns devices to an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
		AssignDevicesToServerV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error)

		// UnassignDevicesFromServerV1 unassigns devices from an MDM server
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
		UnassignDevicesFromServerV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error)
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

// GetDeviceManagementServicesV1 retrieves a list of device management services (MDM servers) in an organization
// URL: GET https://api-business.apple.com/v1/mdmServers
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-device-management-services
func (s *DeviceManagementService) GetDeviceManagementServicesV1(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := APIVersionV1 + "/mdmServers"

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

	var allServers []MDMServer

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, func(pageData []byte) error {
		var pageResponse ResponseMDMServers
		if err := json.Unmarshal(pageData, &pageResponse); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}
		allServers = append(allServers, pageResponse.Data...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ResponseMDMServers{
		Data: allServers,
	}, nil
}

// GetDeviceSerialNumbersForDeviceManagementServiceV1 retrieves a list of device IDs assigned to an MDM server
// URL: GET https://api-business.apple.com/v1/mdmServers/{id}/relationships/devices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-device-management-service
func (s *DeviceManagementService) GetDeviceSerialNumbersForDeviceManagementServiceV1(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}

	endpoint := fmt.Sprintf(APIVersionV1+"/mdmServers/%s/relationships/devices", mdmServerID)

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

	var allDevices []MDMServerDeviceLinkage

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, func(pageData []byte) error {
		var pageResponse ResponseMDMServerDevicesLinkages
		if err := json.Unmarshal(pageData, &pageResponse); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}
		allDevices = append(allDevices, pageResponse.Data...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ResponseMDMServerDevicesLinkages{
		Data: allDevices,
	}, nil
}

// GetAssignedDeviceManagementServiceIDForADeviceV1 retrieves the assigned device management service ID linkage for a device
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/relationships/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-id-for-an-orgdevice
func (s *DeviceManagementService) GetAssignedDeviceManagementServiceIDForADeviceV1(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf(APIVersionV1+"/orgDevices/%s/relationships/assignedServer", deviceID)

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

// GetAssignedDeviceManagementServiceInformationByDeviceIDV1 retrieves the assigned device management service information for a device
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-information-for-an-orgdevice
func (s *DeviceManagementService) GetAssignedDeviceManagementServiceInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf(APIVersionV1+"/orgDevices/%s/assignedServer", deviceID)

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

// AssignDevicesToServerV1 assigns devices to an MDM server
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagementService) AssignDevicesToServerV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("at least one device ID is required")
	}

	endpoint := APIVersionV1 + "/orgDeviceActivities"

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

// UnassignDevicesFromServerV1 unassigns devices from an MDM server
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagementService) UnassignDevicesFromServerV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, error) {
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("at least one device ID is required")
	}

	endpoint := APIVersionV1 + "/orgDeviceActivities"

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
