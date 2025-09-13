package axm2

const (
	// Apple School and Business Manager API base URLs
	AppleBusinessManagerBaseURL = "https://api-business.apple.com"
	AppleSchoolManagerBaseURL   = "https://api-school.apple.com"

	// OAuth endpoints
	TokenEndpoint = "https://account.apple.com/auth/oauth2/token"

	// OAuth constants
	GrantType           = "client_credentials"
	ClientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	BusinessScope       = "business.api"
	SchoolScope         = "school.api"

	// JWT Constants
	TokenAudience = "https://account.apple.com/auth/oauth2/v2/token"

	// API Types
	APITypeABM = "abm" // Apple Business Manager
	APITypeASM = "asm" // Apple School Manager

	// API Endpoints
	OrgDevicesEndpoint          = "/v1/orgDevices"
	MdmServersEndpoint          = "/v1/mdmServers"
	OrgDeviceActivitiesEndpoint = "/v1/orgDeviceActivities"
)
