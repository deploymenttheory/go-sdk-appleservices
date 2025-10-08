package devices

import (
	"time"
)

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
	PurchaseSourceID    string     `json:"purchaseSourceId,omitempty"`
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

// Meta contains pagination and other metadata
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

// Paging contains pagination information
type Paging struct {
	Total int `json:"total,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// Links contains pagination navigation links
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// OrgDevicesResponse represents the response for getting organization devices
type OrgDevicesResponse struct {
	Data  []OrgDevice `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// GetOrganizationDevicesOptions represents the query parameters for getting organization devices
type GetOrganizationDevicesOptions struct {
	// Field selection - fields to return for orgDevices
	// Possible values: serialNumber, addedToOrgDateTime, updatedDateTime, deviceModel,
	// productFamily, productType, deviceCapacity, partNumber, orderNumber, color, status,
	// orderDateTime, imei, meid, eid, wifiMacAddress, bluetoothMacAddress, purchaseSourceId,
	// purchaseSourceType, assignedServer
	Fields []string `json:"fields,omitempty"`

	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// GetDeviceInformationOptions represents the query parameters for getting device information
type GetDeviceInformationOptions struct {
	// Field selection for the device information response
	Fields []string `json:"fields,omitempty"`
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

// ProfileStatus constants for device profile status
const (
	ProfileStatusEmpty    = "empty"
	ProfileStatusAssigned = "assigned"
	ProfileStatusPushed   = "pushed"
	ProfileStatusRemoved  = "removed"
)

// DeviceModel constants for common device models
const (
	ModeliPhone     = "iPhone"
	ModeliPad       = "iPad"
	ModelMac        = "Mac"
	ModelAppleTV    = "Apple TV"
	ModelAppleWatch = "Apple Watch"
)

// DeviceColor constants for device colors
const (
	ColorSilver    = "Silver"
	ColorGold      = "Gold"
	ColorSpaceGray = "Space Gray"
	ColorRoseGold  = "Rose Gold"
	ColorBlack     = "Black"
	ColorWhite     = "White"
	ColorRed       = "Red"
	ColorBlue      = "Blue"
	ColorGreen     = "Green"
	ColorYellow    = "Yellow"
	ColorPurple    = "Purple"
)

// OrgDevice field constants for field selection
const (
	FieldSerialNumber        = "serialNumber"
	FieldAddedToOrgDateTime  = "addedToOrgDateTime"
	FieldUpdatedDateTime     = "updatedDateTime"
	FieldDeviceModel         = "deviceModel"
	FieldProductFamily       = "productFamily"
	FieldProductType         = "productType"
	FieldDeviceCapacity      = "deviceCapacity"
	FieldPartNumber          = "partNumber"
	FieldOrderNumber         = "orderNumber"
	FieldColor               = "color"
	FieldStatus              = "status"
	FieldOrderDateTime       = "orderDateTime"
	FieldIMEI                = "imei"
	FieldMEID                = "meid"
	FieldEID                 = "eid"
	FieldWiFiMACAddress      = "wifiMacAddress"
	FieldBluetoothMACAddress = "bluetoothMacAddress"
	FieldPurchaseSourceID    = "purchaseSourceId"
	FieldPurchaseSourceType  = "purchaseSourceType"
	FieldAssignedServer      = "assignedServer"
)

// Device status constants
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// Product family constants
const (
	ProductFamilyiPhone = "iPhone"
	ProductFamilyiPad   = "iPad"
	ProductFamilyMac    = "Mac"
)
