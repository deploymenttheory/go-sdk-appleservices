package shared

import "context"

// HTTPClientInterface represents the HTTP client interface
type HTTPClientInterface interface {
	R() RequestInterface
}

// RequestInterface represents a request interface
type RequestInterface interface {
	SetContext(ctx context.Context) RequestInterface
	SetResult(result any) RequestInterface
	SetError(err any) RequestInterface
	SetBody(body any) RequestInterface
	SetQueryParam(param, value string) RequestInterface
	Get(url string) (ResponseInterface, error)
	Post(url string) (ResponseInterface, error)
}

// ResponseInterface represents a response interface
type ResponseInterface interface {
	IsError() bool
	StatusCode() int
	String() string
}

// PaginatedResponse represents a response that supports pagination
type PaginatedResponse interface {
	GetData() any
	GetNextURL() string
	AppendData(any)
}
