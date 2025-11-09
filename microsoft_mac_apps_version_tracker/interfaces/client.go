package interfaces

import "context"

// HTTPClient interface that services will use
type HTTPClient interface {
	Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	GetHTTPClient() any
}
