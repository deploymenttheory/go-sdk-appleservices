package client

// DefaultUserAgent is the default User-Agent header value for all requests.
const (
	DefaultUserAgent = "go-api-sdk-apple/1.0.0"
	Version          = "1.0.0"
)

// The following constants are re-exported from the constants package so that
// existing code and tests in the client package can reference them without
// importing the constants package directly.
const (
	DefaultBaseURL            = "https://api-business.apple.com"
	DefaultOAuthTokenEndpoint = "https://account.apple.com/auth/oauth2/v2/token"
	DefaultJWTAudience        = "appstoreconnect-v1"
	ScopeBusinessAPI          = "business.api"
	ScopeSchoolAPI            = "school.api"
)
