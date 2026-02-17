package client

// API Base URL
const (
	DefaultBaseURL = "https://api-business.apple.com"
)

// OAuth scope constants
const (
	ScopeBusinessAPI = "business.api"
	ScopeSchoolAPI   = "school.api"
)

// Default OAuth endpoints
const (
	DefaultOAuthTokenEndpoint = "https://account.apple.com/auth/oauth2/v2/token"
)

// Default OAuth audience
const (
	DefaultJWTAudience = "appstoreconnect-v1"
)

// User Agent
const (
	DefaultUserAgent = "go-api-sdk-apple/1.0.0"
	Version          = "1.0.0"
)
