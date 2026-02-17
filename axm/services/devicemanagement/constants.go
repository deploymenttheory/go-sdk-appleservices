package devicemanagement

// API Version
const (
	APIVersionV1 = "/v1"
)

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
	FieldServerName      = "serverName"
	FieldServerType      = "serverType"
	FieldCreatedDateTime = "createdDateTime"
	FieldUpdatedDateTime = "updatedDateTime"
	FieldDevices         = "devices"
)
