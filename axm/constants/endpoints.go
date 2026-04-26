package constants

// API base URL
const (
	DefaultBaseURL = "https://api-business.apple.com"
)

// API version prefix
const (
	APIVersionV1 = "/v1"
)

// Endpoint path constants for the Apple Business Manager API
const (
	EndpointOrgDevices          = APIVersionV1 + "/orgDevices"
	EndpointMDMServers          = APIVersionV1 + "/mdmServers"
	EndpointOrgDeviceActivities = APIVersionV1 + "/orgDeviceActivities"
	EndpointAuditEvents         = APIVersionV1 + "/auditEvents"
	EndpointUsers               = APIVersionV1 + "/users"
	EndpointUserGroups          = APIVersionV1 + "/userGroups"
	EndpointApps                = APIVersionV1 + "/apps"
	EndpointPackages            = APIVersionV1 + "/packages"
	EndpointConfigurations      = APIVersionV1 + "/configurations"
	EndpointBlueprints          = APIVersionV1 + "/blueprints"
)
