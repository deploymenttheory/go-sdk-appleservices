package axm

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// Get executes a GET request following Resty v3 best practices
func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	// Add query parameters
	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "GET", path)
}

// Post executes a POST request with JSON body following Resty v3 best practices
func (c *Client) Post(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "POST", path)
}

// PostWithQuery executes a POST request with both body and query parameters
func (c *Client) PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	// Add query parameters
	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "POST", path)
}

// Put executes a PUT request following Resty v3 best practices
func (c *Client) Put(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "PUT", path)
}

// Patch executes a PATCH request following Resty v3 best practices
func (c *Client) Patch(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "PATCH", path)
}

// Delete executes a DELETE request following Resty v3 best practices
func (c *Client) Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	// Add query parameters
	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "DELETE", path)
}

// DeleteWithBody executes a DELETE request with body (for bulk operations)
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	// Add headers
	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "DELETE", path)
}

// PostMultipart executes a POST request with multipart data following Resty v3 best practices
func (c *Client) PostMultipart(ctx context.Context, path string, files map[string]string, fields map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	// Add form fields
	for k, v := range fields {
		req.SetFormData(map[string]string{k: v})
	}

	// Add files
	for fieldName, filePath := range files {
		req.SetFile(fieldName, filePath)
	}

	return c.executeRequest(req, "POST", path)
}

// executeRequest is a centralized request executor that handles error processing
func (c *Client) executeRequest(req *resty.Request, method, path string) error {
	var apiErr ErrorResponse
	req.SetError(&apiErr)

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Handle API errors
	if resp.IsError() {
		return c.errorHandler.HandleError(resp, &apiErr)
	}

	return nil
}
