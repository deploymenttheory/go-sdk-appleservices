package devicemanagement

// Activity type constants
const (
	ActivityTypeAssignDevices   = "ASSIGN_DEVICES"
	ActivityTypeUnassignDevices = "UNASSIGN_DEVICES"
)

// Activity status constants
const (
	ActivityStatusInProgress = "IN_PROGRESS"
	ActivityStatusCompleted  = "COMPLETED"
	ActivityStatusFailed     = "FAILED"
)

// Activity sub-status constants
const (
	ActivitySubStatusSubmitted  = "SUBMITTED"
	ActivitySubStatusProcessing = "PROCESSING"
)

// MDM Server field constants for field selection
const (
	FieldServerName             = "serverName"
	FieldServerType             = "serverType"
	FieldEnableMdmDisownFlag    = "enableMdmDisownFlag"
	FieldDefaultProductFamilies = "defaultProductFamilies"
	FieldStatus                 = "status"
	FieldDeviceCount            = "deviceCount"
	FieldLastConnectedDateTime  = "lastConnectedDateTime"
	FieldLastConnectedIp        = "lastConnectedIp"
	FieldCreatedDateTime        = "createdDateTime"
	FieldUpdatedDateTime        = "updatedDateTime"
	FieldDevices                = "devices"
)

// MDM server status constants
const (
	MDMServerStatusActive   = "ACTIVE"
	MDMServerStatusInactive = "INACTIVE"
)
