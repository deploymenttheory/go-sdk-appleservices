package devices

import (
	"context"
	"fmt"
)

// GetOrganizationDevices retrieves a list of devices in an organization that enroll using Automated Device Enrollment
// URL: GET https://api-business.apple.com/v1/orgDevices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
func (c *Client) GetOrganizationDevices(ctx context.Context, opts *GetOrganizationDevicesOptions) (*OrgDevicesResponse, error) {
	if opts == nil {
		opts = &GetOrganizationDevicesOptions{}
	}

	endpoint := "/orgDevices"

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := c.client.QueryBuilder()

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

	err := c.client.GetPaginated(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeviceInformationByDeviceID retrieves information about a specific device in an organization
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
func (c *Client) GetDeviceInformationByDeviceID(ctx context.Context, deviceID string, opts *GetDeviceInformationOptions) (*OrgDeviceResponse, error) {
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

	queryParams := c.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		queryParams.AddStringSlice("fields[orgDevices]", opts.Fields)
	}

	var result OrgDeviceResponse

	err := c.client.Get(ctx, endpoint, queryParams.Build(), headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
