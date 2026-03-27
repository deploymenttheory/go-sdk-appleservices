package constants

// OAuth 2.0 / JWT Authentication endpoints and configuration
const (
	DefaultOAuthTokenEndpoint = "https://account.apple.com/auth/oauth2/v2/token"
	DefaultJWTAudience        = "appstoreconnect-v1"
)

// OAuth scope constants
const (
	ScopeBusinessAPI = "business.api"
	ScopeSchoolAPI   = "school.api"
)
