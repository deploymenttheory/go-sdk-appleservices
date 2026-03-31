package devicemanagement

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/constants"
	"resty.dev/v3"
)

// DeviceManagement handles communication with the device management
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessmanagerapi/
type (
	DeviceManagement struct {
		client client.Client
	}
)

// NewService creates a new device management service.
func NewService(c client.Client) *DeviceManagement {
	return &DeviceManagement{client: c}
}

// GetV1 retrieves a list of device management services (MDM servers) in an organization.
// URL: GET https://api-business.apple.com/v1/mdmServers
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-device-management-services
func (s *DeviceManagement) GetV1(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allServers []MDMServer
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointMDMServers, func(pageData []byte) error {
			var pageResponse ResponseMDMServers
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allServers = append(allServers, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &ResponseMDMServers{
		Data:  allServers,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetDeviceSerialNumbersByServerIDV1 retrieves a list of device IDs assigned to an MDM server.
// URL: GET https://api-business.apple.com/v1/mdmServers/{id}/relationships/devices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-device-management-service
func (s *DeviceManagement) GetDeviceSerialNumbersByServerIDV1(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointMDMServers+"/%s/relationships/devices", mdmServerID)

	params := s.client.QueryBuilder()

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allDevices []MDMServerDeviceLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse ResponseMDMServerDevicesLinkages
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allDevices = append(allDevices, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &ResponseMDMServerDevicesLinkages{
		Data:  allDevices,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetAssignedServerIDByDeviceIDV1 retrieves the assigned device management service ID linkage for a device.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/relationships/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-id-for-an-orgdevice
func (s *DeviceManagement) GetAssignedServerIDByDeviceIDV1(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointOrgDevices+"/%s/relationships/assignedServer", deviceID)

	var result ResponseOrgDeviceAssignedServerLinkage

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetAssignedServerInfoByDeviceIDV1 retrieves the assigned device management service information for a device.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-information-for-an-orgdevice
func (s *DeviceManagement) GetAssignedServerInfoByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointOrgDevices+"/%s/assignedServer", deviceID)

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	var result MDMServerResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// AssignDevicesV1 assigns devices to an MDM server.
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagement) AssignDevicesV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, nil, fmt.Errorf("at least one device ID is required")
	}

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

	var result ResponseOrgDeviceActivity

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(request).
		SetResult(&result).
		Post(constants.EndpointOrgDeviceActivities)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UnassignDevicesV1 unassigns devices from an MDM server.
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-an-orgdeviceactivity
func (s *DeviceManagement) UnassignDevicesV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}
	if len(deviceIDs) == 0 {
		return nil, nil, fmt.Errorf("at least one device ID is required")
	}

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

	var result ResponseOrgDeviceActivity

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(request).
		SetResult(&result).
		Post(constants.EndpointOrgDeviceActivities)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
