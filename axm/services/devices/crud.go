package devices

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/interfaces"
)

type (
	// DevicesServiceInterface defines the interface for device operations.
	DevicesServiceInterface interface {
		// GetOrganizationDevicesV1 retrieves a list of devices in an organization
		// that enroll using Automated Device Enrollment.
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
		GetOrganizationDevicesV1(ctx context.Context, opts *RequestQueryOptions) (*OrgDevicesResponse, error)

		// GetDeviceInformationByDeviceIDV1 retrieves information about a specific device
		// in an organization.
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
		GetDeviceInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*OrgDeviceResponse, error)

		// GetAppleCareInformationByDeviceIDV1 retrieves AppleCare coverage information
		// for a specific device.
		//
		// Apple Business Manager API docs:
		// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-apple-care-coverage-for-an-orgdevice
		GetAppleCareInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*AppleCareCoverageResponse, error)
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

// GetOrganizationDevicesV1 retrieves a list of devices in an organization that enroll using Automated Device Enrollment
// URL: GET https://api-business.apple.com/v1/orgDevices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
func (s *DevicesService) GetOrganizationDevicesV1(ctx context.Context, opts *RequestQueryOptions) (*OrgDevicesResponse, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := APIVersionV1 + "/orgDevices"

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

	var allDevices []OrgDevice

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, func(pageData []byte) error {
		var pageResponse OrgDevicesResponse
		if err := json.Unmarshal(pageData, &pageResponse); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}
		allDevices = append(allDevices, pageResponse.Data...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &OrgDevicesResponse{
		Data: allDevices,
	}, nil
}

// GetDeviceInformationByDeviceIDV1 retrieves information about a specific device in an organization
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
func (s *DevicesService) GetDeviceInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*OrgDeviceResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := APIVersionV1 + "/orgDevices/" + deviceID

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
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

// GetAppleCareInformationByDeviceIDV1 retrieves AppleCare coverage information for a specific device
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/appleCareCoverage
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-apple-care-coverage-for-an-orgdevice
func (s *DevicesService) GetAppleCareInformationByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*AppleCareCoverageResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	endpoint := APIVersionV1 + "/orgDevices/" + deviceID + "/appleCareCoverage"

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	queryParams := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[appleCareCoverage]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		queryParams.AddInt("limit", opts.Limit)
	}

	var allCoverage []AppleCareCoverage

	err := s.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, func(pageData []byte) error {
		var pageResponse AppleCareCoverageResponse
		if err := json.Unmarshal(pageData, &pageResponse); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}
		allCoverage = append(allCoverage, pageResponse.Data...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AppleCareCoverageResponse{
		Data: allCoverage,
	}, nil
}
