package constants

// API base URL
const (
	DefaultBaseURL = "https://appstoreconnect.apple.com"
)

// API version prefix
const (
	APIVersionV2 = "/notary/v2"
)

// Endpoint path constants for the Apple Notary API
const (
	EndpointSubmissions = APIVersionV2 + "/submissions"
	EndpointLogs        = "/logs"
)
