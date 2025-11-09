package client

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// Get executes a GET request
func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "GET", path)
}

// executeRequest is a centralized request executor
func (c *Client) executeRequest(req *resty.Request, method, path string) error {
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API returned error status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
