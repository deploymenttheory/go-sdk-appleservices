package devices

import "time"

// Shared types for pagination and links
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// OrgDevice represents a device in the Apple Business Manager system based on the API specification
type OrgDevice struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes *OrgDeviceAttributes `json:"attributes,omitempty"`
}

// OrgDeviceAttributes contains the device attributes
type OrgDeviceAttributes struct {
	SerialNumber        string     `json:"serialNumber,omitempty"`
	AddedToOrgDateTime  *time.Time `json:"addedToOrgDateTime,omitempty"`
	UpdatedDateTime     *time.Time `json:"updatedDateTime,omitempty"`
	DeviceModel         string     `json:"deviceModel,omitempty"`
	ProductFamily       string     `json:"productFamily,omitempty"`
	ProductType         string     `json:"productType,omitempty"`
	DeviceCapacity      string     `json:"deviceCapacity,omitempty"`
	PartNumber          string     `json:"partNumber,omitempty"`
	OrderNumber         string     `json:"orderNumber,omitempty"`
	Color               string     `json:"color,omitempty"`
	Status              string     `json:"status,omitempty"`
	OrderDateTime       *time.Time `json:"orderDateTime,omitempty"`
	IMEI                []string   `json:"imei,omitempty"`
	MEID                []string   `json:"meid,omitempty"`
	EID                 string     `json:"eid,omitempty"`
	WiFiMACAddress      string     `json:"wifiMacAddress,omitempty"`
	BluetoothMACAddress string     `json:"bluetoothMacAddress,omitempty"`
	PurchaseSourceUid   string     `json:"purchaseSourceUid,omitempty"`
	PurchaseSourceType  string     `json:"purchaseSourceType,omitempty"`
	AssignedServer      string     `json:"assignedServer,omitempty"`
}

// OrgDeviceResponse represents the response for a single device
type OrgDeviceResponse struct {
	Data OrgDevice `json:"data"`
}

// OrgDeviceFilter represents filter options for organization devices
type OrgDeviceFilter struct {
	SerialNumber  string `json:"serialNumber,omitempty"`
	DeviceModel   string `json:"deviceModel,omitempty"`
	ProductFamily string `json:"productFamily,omitempty"`
	Color         string `json:"color,omitempty"`
	Status        string `json:"status,omitempty"`
}

// OrgDevicesResponse represents the response for getting organization devices
type OrgDevicesResponse struct {
	Data  []OrgDevice `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// RequestQueryOptions represents the query parameters for getting organization devices
type RequestQueryOptions struct {
	// Field selection - fields to return for orgDevices
	// Possible values: serialNumber, addedToOrgDateTime, updatedDateTime, deviceModel,
	// productFamily, productType, deviceCapacity, partNumber, orderNumber, color, status,
	// orderDateTime, imei, meid, eid, wifiMacAddress, bluetoothMacAddress, purchaseSourceId,
	// purchaseSourceType, assignedServer
	Fields []string `json:"fields,omitempty"`

	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// Legacy Device struct for backward compatibility
type Device struct {
	SerialNumber       string    `json:"serial_number"`
	Model              string    `json:"model"`
	Description        string    `json:"description"`
	Color              string    `json:"color"`
	AssetTag           string    `json:"asset_tag"`
	ProfileStatus      string    `json:"profile_status"`
	ProfileUUID        string    `json:"profile_uuid"`
	ProfileAssignTime  time.Time `json:"profile_assign_time"`
	DeviceAssignedBy   string    `json:"device_assigned_by"`
	DeviceAssignedDate time.Time `json:"device_assigned_date"`
	OpType             string    `json:"op_type,omitempty"`
	OpDate             time.Time `json:"op_date,omitempty"`
}
