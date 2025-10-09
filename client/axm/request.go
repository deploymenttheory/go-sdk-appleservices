package axm

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// Get executes a GET request following Resty v3 best practices.
// Supports query parameters, custom headers, and automatic response parsing.
//
// This method uses Resty's SetQueryParam for URL parameters and SetHeader for custom headers.
// The response is automatically unmarshaled into the provided result interface{}.
//
// Reference: https://resty.dev/docs/request-query-params/
// Reference: https://resty.dev/docs/response-auto-parse/
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

// Post executes a POST request with JSON body following Resty v3 best practices.
// Automatically sets Content-Type to application/json and handles request body serialization.
//
// The body parameter is automatically marshaled to JSON using Resty's SetBody method.
// Custom headers can be provided to override defaults or add additional headers.
//
// Reference: https://resty.dev/docs/request-body-types/
// Reference: https://resty.dev/docs/response-auto-parse/
func (c *Client) Post(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "POST", path)
}

// PostWithQuery executes a POST request with both body and query parameters.
// Combines the functionality of POST requests with URL query parameter support.
//
// This method is useful for APIs that require both request body data and URL parameters.
// The body is JSON-serialized while query parameters are added to the URL.
//
// Reference: https://resty.dev/docs/request-query-params/
// Reference: https://resty.dev/docs/request-body-types/
func (c *Client) PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "POST", path)
}

// Put executes a PUT request following Resty v3 best practices.
// Typically used for updating existing resources with complete replacement.
//
// The body parameter is automatically marshaled to JSON using Resty's SetBody method.
// Follows RESTful conventions for resource updates and modifications.
//
// Reference: https://resty.dev/docs/request-body-types/
// Reference: https://resty.dev/docs/response-auto-parse/
func (c *Client) Put(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "PUT", path)
}

// Patch executes a PATCH request following Resty v3 best practices.
// Typically used for partial updates of existing resources.
//
// The body parameter is automatically marshaled to JSON using Resty's SetBody method.
// Follows RESTful conventions for partial resource modifications.
//
// Reference: https://resty.dev/docs/request-body-types/
// Reference: https://resty.dev/docs/response-auto-parse/
func (c *Client) Patch(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "PATCH", path)
}

// Delete executes a DELETE request following Resty v3 best practices.
// Supports query parameters for filtering or specifying deletion criteria.
//
// This method uses Resty's SetQueryParam for URL parameters and follows
// RESTful conventions for resource deletion operations.
//
// Reference: https://resty.dev/docs/request-query-params/
// Reference: https://resty.dev/docs/response-auto-parse/
func (c *Client) Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error {
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

	return c.executeRequest(req, "DELETE", path)
}

// DeleteWithBody executes a DELETE request with body (for bulk operations).
// Useful for bulk deletion operations where multiple resources are specified in the request body.
//
// While not always RESTful, some APIs require request bodies for DELETE operations
// to specify multiple resources or complex deletion criteria.
//
// Reference: https://resty.dev/docs/request-body-types/
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	for k, v := range headers {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	return c.executeRequest(req, "DELETE", path)
}

// PostMultipart executes a POST request with multipart form data following Resty v3 best practices.
// Supports file uploads and form fields in a single multipart/form-data request.
//
// This method uses Resty's SetFile for file uploads and SetFormData for form fields.
// Automatically sets the appropriate Content-Type header for multipart requests.
//
// Files map: fieldName -> filePath
// Fields map: fieldName -> fieldValue
//
// Reference: https://resty.dev/docs/multipart/
// Reference: https://resty.dev/docs/form-data/
func (c *Client) PostMultipart(ctx context.Context, path string, files map[string]string, fields map[string]string, result any) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range fields {
		req.SetFormData(map[string]string{k: v})
	}

	for fieldName, filePath := range files {
		req.SetFile(fieldName, filePath)
	}

	return c.executeRequest(req, "POST", path)
}

// executeRequest is a centralized request executor that handles error processing.
// Provides consistent error handling, response processing, and HTTP method routing.
//
// This method:
// - Sets up error response handling using Resty's SetError
// - Routes requests to appropriate HTTP methods
// - Processes API errors through the configured error handler
// - Ensures consistent error formatting across all request types
//
// Reference: https://resty.dev/docs/response-auto-parse/
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
