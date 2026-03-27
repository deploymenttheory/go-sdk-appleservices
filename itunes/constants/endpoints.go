package constants

const (
	// DefaultBaseURL is the base URL for the iTunes Search API.
	DefaultBaseURL = "https://itunes.apple.com"

	// EndpointSearch is the path for the iTunes Search endpoint.
	EndpointSearch = "/search"

	// EndpointLookup is the path for the iTunes Lookup endpoint.
	EndpointLookup = "/lookup"

	// ApplicationJSON is the MIME type for JSON content.
	ApplicationJSON = "application/json"

	// DefaultLimit is the default number of results to return.
	DefaultLimit = 50

	// MaxLimit is the maximum number of results the API accepts.
	MaxLimit = 200
)
