package devices

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
	FieldPurchaseSourceUid   = "purchaseSourceUid"
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