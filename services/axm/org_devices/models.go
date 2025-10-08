package org_devices

import (
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/mdm_servers"
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

// OrgDevicesResponse represents the API response structure for multiple devices
type OrgDevicesResponse struct {
	Data  []OrgDevice        `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// GetData returns the data slice for pagination
func (r *OrgDevicesResponse) GetData() any {
	return r.Data
}

// GetNextURL returns the next URL for pagination
func (r *OrgDevicesResponse) GetNextURL() string {
	return r.Links.Next
}

// AppendData appends data from another response for pagination
func (r *OrgDevicesResponse) AppendData(data any) {
	if existingData, ok := data.([]OrgDevice); ok {
		r.Data = append(r.Data, existingData...)
	}
}

// OrgDeviceResponse represents the API response structure for a single device
type OrgDeviceResponse struct {
	Data  OrgDevice          `json:"data"`
	Links PagedDocumentLinks `json:"links"`
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

// ResourceLinkage represents a linkage between resources
type ResourceLinkage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// OrgDeviceAssignedServerLinkageResponse represents the response for device assigned server relationships
type OrgDeviceAssignedServerLinkageResponse struct {
	Data  *ResourceLinkage                            `json:"data"`
	Links OrgDeviceAssignedServerLinkageResponseLinks `json:"links"`
}

// OrgDeviceAssignedServerLinkageResponseLinks contains links for assigned server relationships
type OrgDeviceAssignedServerLinkageResponseLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// Use MdmServer types from the mdm_servers package
type MdmServer = mdm_servers.MdmServer
type MdmServerAttributes = mdm_servers.MdmServerAttributes
type MdmServerResponse = mdm_servers.MdmServerResponse

// ResourceLinks contains self-links to requested resources
type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Errors []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source *struct {
		Pointer   string `json:"pointer,omitempty"`
		Parameter string `json:"parameter,omitempty"`
	} `json:"source,omitempty"`
}
