package devices

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/interfaces"
)

type (
	// DevicesServiceInterface defines the interface for device operations.
	DevicesServiceInterface interface {
		// GetOrganizationDevices retrieves a list of devices in an organization
		// that enroll using Automated Device Enrollment.
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
		GetOrganizationDevices(ctx context.Context, opts *GetOrganizationDevicesOptions) (*OrgDevicesResponse, error)

		// GetDeviceInformationByDeviceID retrieves information about a specific device
		// in an organization.
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
		GetDeviceInformationByDeviceID(ctx context.Context, deviceID string, opts *GetDeviceInformationOptions) (*OrgDeviceResponse, error)
	}

	// DevicetService handles communication with the device
	// related methods of the Apple Business Manager API.
	//
	// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessmanagerapi/
	DevicesService struct {
		client interfaces.HTTPClient
	}
)

var _ DevicesServiceInterface = (*DevicesService)(nil)

// NewService creates a new devices service
func NewService(client interfaces.HTTPClient) *DevicesService {
	return &DevicesService{
		client: client,
	}
}

// GetOrganizationDevices retrieves a list of devices in an organization that enroll using Automated Device Enrollment
// URL: GET https://api-business.apple.com/v1/orgDevices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
func (s *DevicesService) GetOrganizationDevices(ctx context.Context, opts *GetOrganizationDevicesOptions) (*OrgDevicesResponse, error) {
	if opts == nil {
		opts = &GetOrganizationDevicesOptions{}
	}

	endpoint := "/orgDevices"

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[orgDevices]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		queryParams.AddInt("limit", opts.Limit)
	}

	var result OrgDevicesResponse

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeviceInformationByDeviceID retrieves information about a specific device in an organization
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
func (s *DevicesService) GetDeviceInformationByDeviceID(ctx context.Context, deviceID string, opts *GetDeviceInformationOptions) (*OrgDeviceResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := "/orgDevices/" + deviceID

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &GetDeviceInformationOptions{}
	}

	queryParams := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[orgDevices]", opts.Fields)
	}

	var result OrgDeviceResponse

	err := s.client.Get(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
