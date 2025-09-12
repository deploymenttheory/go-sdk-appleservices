package axm

import (
	"encoding/json"
	"fmt"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"go.uber.org/zap"
)

const (
	OrgDevicesEndpoint = "/v1/orgDevices"
)

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

// GetOrgDevices retrieves axm organization devices using the QueryBuilder
func (c *Client) GetOrgDevices(queryBuilder *client.QueryBuilder) ([]OrgDevice, error) {
	c.logger.Debug("Getting organization devices with automatic pagination using QueryBuilder")

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	resp, _, err := c.axmClient.Get(OrgDevicesEndpoint, queryBuilder, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization devices: %w", err)
	}

	var response OrgDevicesResponse

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		c.logger.Error("Failed to unmarshal organization devices response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Debug("Successfully retrieved organization devices",
		zap.Int("total_count", len(response.Data)))

	return response.Data, nil
}
