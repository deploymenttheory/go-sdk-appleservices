package axm2

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-apple/client/shared"
	"resty.dev/v3"
)

// HTTPClientWrapper wraps resty.Client to implement service interfaces
type HTTPClientWrapper struct {
	client *resty.Client
	parent *Client
}

// GetHTTPClient returns a wrapped HTTP client interface
func (c *Client) GetHTTPClient() shared.HTTPClientInterface {
	return &HTTPClientWrapper{client: c.httpClient, parent: c}
}

// R returns a wrapped request interface
func (h *HTTPClientWrapper) R() shared.RequestInterface {
	return &RequestWrapper{req: h.client.R(), parent: h.parent}
}

// RequestWrapper wraps resty.Request to implement RequestInterface
type RequestWrapper struct {
	req    *resty.Request
	parent *Client
}

func (r *RequestWrapper) SetContext(ctx context.Context) shared.RequestInterface {
	r.req.SetContext(ctx)
	return r
}

func (r *RequestWrapper) SetResult(result any) shared.RequestInterface {
	r.req.SetResult(result)
	return r
}

func (r *RequestWrapper) SetError(err any) shared.RequestInterface {
	r.req.SetError(err)
	return r
}

func (r *RequestWrapper) SetBody(body any) shared.RequestInterface {
	r.req.SetBody(body)
	return r
}

func (r *RequestWrapper) SetQueryParam(param, value string) shared.RequestInterface {
	r.req.SetQueryParam(param, value)
	return r
}

func (r *RequestWrapper) Get(url string) (shared.ResponseInterface, error) {
	resp, err := r.req.Get(url)
	if err != nil {
		return nil, err
	}
	return &ResponseWrapper{resp: resp}, nil
}

func (r *RequestWrapper) Post(url string) (shared.ResponseInterface, error) {
	resp, err := r.req.Post(url)
	if err != nil {
		return nil, err
	}
	return &ResponseWrapper{resp: resp}, nil
}

// ResponseWrapper wraps resty.Response to implement ResponseInterface
type ResponseWrapper struct {
	resp *resty.Response
}

func (r *ResponseWrapper) IsError() bool {
	return r.resp.IsError()
}

func (r *ResponseWrapper) StatusCode() int {
	return r.resp.StatusCode()
}

func (r *ResponseWrapper) String() string {
	return r.resp.String()
}
