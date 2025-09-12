package axm2

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// OrgDevice represents an organization device resource
type OrgDevice struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id"`
	Attributes    OrgDeviceAttributes    `json:"attributes"`
	Relationships OrgDeviceRelationships `json:"relationships,omitempty"`
	Links         OrgDeviceLinks         `json:"links,omitempty"`
}

// OrgDeviceAttributes contains the device attributes
type OrgDeviceAttributes struct {
	SerialNumber        string   `json:"serialNumber,omitempty"`
	AddedToOrgDateTime  string   `json:"addedToOrgDateTime,omitempty"`
	UpdatedDateTime     string   `json:"updatedDateTime,omitempty"`
	DeviceModel         string   `json:"deviceModel,omitempty"`
	ProductFamily       string   `json:"productFamily,omitempty"`
	ProductType         string   `json:"productType,omitempty"`
	DeviceCapacity      string   `json:"deviceCapacity,omitempty"`
	PartNumber          string   `json:"partNumber,omitempty"`
	OrderNumber         string   `json:"orderNumber,omitempty"`
	Color               string   `json:"color,omitempty"`
	Status              string   `json:"status,omitempty"`
	OrderDateTime       string   `json:"orderDateTime,omitempty"`
	IMEI                []string `json:"imei,omitempty"`
	MEID                []string `json:"meid,omitempty"`
	EID                 string   `json:"eid,omitempty"`
	WifiMacAddress      string   `json:"wifiMacAddress,omitempty"`
	BluetoothMacAddress string   `json:"bluetoothMacAddress,omitempty"`
	PurchaseSourceUid   string   `json:"purchaseSourceUid,omitempty"`
	PurchaseSourceType  string   `json:"purchaseSourceType,omitempty"`
}

// OrgDeviceRelationships contains relationship links
type OrgDeviceRelationships struct {
	AssignedServer *RelationshipLinks `json:"assignedServer,omitempty"`
}

// RelationshipLinks contains links for a relationship
type RelationshipLinks struct {
	Links map[string]string `json:"links,omitempty"`
}

// OrgDeviceLinks contains resource links
type OrgDeviceLinks struct {
	Self string `json:"self,omitempty"`
}

// OrgDevicesResponse represents the API response structure
type OrgDevicesResponse struct {
	Data  []OrgDevice        `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// PagedDocumentLinks contains navigational links for paged responses
type PagedDocumentLinks struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// PagingInformation contains metadata about pagination
type PagingInformation struct {
	Paging *PagingDetails `json:"paging,omitempty"`
}

// PagingDetails contains specific paging information
type PagingDetails struct {
	NextCursor string `json:"nextCursor,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// GetOrgDevices retrieves organization devices with automatic pagination using Resty v3 patterns
func (c *Client) GetOrgDevices(ctx context.Context, queryBuilder *QueryBuilder) ([]OrgDevice, error) {
	c.logger.Debug("Getting organization devices with automatic pagination")

	var allDevices []OrgDevice
	nextURL := OrgDevicesEndpoint

	// Build initial query parameters
	var queryParams map[string]string
	if queryBuilder != nil {
		queryParams = queryBuilder.Build()
	}

	pageCount := 0
	for nextURL != "" {
		pageCount++
		c.logger.Debug("Fetching page", zap.Int("page", pageCount), zap.String("url", nextURL))

		// Use Resty v3 SetResult pattern for automatic unmarshaling
		var pageResponse OrgDevicesResponse
		var apiError APIError

		request := c.httpClient.R().
			SetContext(ctx).
			SetResult(&pageResponse).
			SetError(&apiError)

		// Add query parameters for first page only
		if pageCount == 1 && queryParams != nil {
			for key, value := range queryParams {
				request.SetQueryParam(key, value)
			}
		}

		response, err := request.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to execute GET request (page %d): %w", pageCount, err)
		}

		if response.IsError() {
			c.logger.Error("API error getting devices",
				zap.Int("page", pageCount),
				zap.Int("status_code", response.StatusCode()),
				zap.Any("error", apiError))
			return nil, fmt.Errorf("API error (page %d): %d %s", pageCount, response.StatusCode(), response.String())
		}

		// Add devices from this page
		allDevices = append(allDevices, pageResponse.Data...)

		// Determine next URL for pagination
		nextURL = pageResponse.Links.Next

		c.logger.Debug("Page fetched successfully",
			zap.Int("page", pageCount),
			zap.Int("items_this_page", len(pageResponse.Data)),
			zap.Int("total_items_so_far", len(allDevices)),
			zap.String("next_url", nextURL))
	}

	c.logger.Info("Successfully retrieved organization devices",
		zap.Int("total_pages", pageCount),
		zap.Int("device_count", len(allDevices)))

	return allDevices, nil
}

// GetOrgDevice retrieves a single organization device by ID using Resty v3 patterns
func (c *Client) GetOrgDevice(ctx context.Context, deviceID string, queryBuilder *QueryBuilder) (*OrgDevice, error) {
	c.logger.Debug("Getting organization device", zap.String("device_id", deviceID))

	if deviceID == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s", OrgDevicesEndpoint, deviceID)

	// Use Resty v3 pattern with SetResult for automatic unmarshaling
	var deviceResponse struct {
		Data OrgDevice `json:"data"`
	}
	var errorResponse APIError

	request := c.httpClient.R().SetContext(ctx).
		SetResult(&deviceResponse).
		SetError(&errorResponse)

	// Add query parameters if provided
	if queryBuilder != nil {
		queryParams := queryBuilder.Build()
		for key, value := range queryParams {
			request.SetQueryParam(key, value)
		}
	}

	response, err := request.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization device: %w", err)
	}

	if response.IsError() {
		c.logger.Error("API error getting device",
			zap.String("device_id", deviceID),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", errorResponse))
		return nil, fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	c.logger.Debug("Successfully retrieved organization device",
		zap.String("device_id", deviceID),
		zap.String("serial_number", deviceResponse.Data.Attributes.SerialNumber))

	return &deviceResponse.Data, nil
}

// NewQueryBuilder creates a new query parameter builder
func (c *Client) NewQueryBuilder() *QueryBuilder {
	return NewQueryBuilder()
}
