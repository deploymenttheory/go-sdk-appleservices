package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// PostWithParams makes a POST request using a map of parameters
func (c *AXMClient) PostWithParams(endpoint string, body interface{}, params map[string]string) (*resty.Response, error) {
	return c.post(endpoint, body, params)
}

// Post makes a POST request using a QueryBuilder
func (c *AXMClient) Post(endpoint string, body interface{}, queryBuilder *QueryBuilder) (*resty.Response, error) {
	var params map[string]string
	if queryBuilder != nil {
		params = queryBuilder.Build()
	}
	return c.post(endpoint, body, params)
}

// post is the internal method that handles the actual POST request
func (c *AXMClient) post(endpoint string, body interface{}, params map[string]string) (*resty.Response, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	req := c.HTTP.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetHeader("Content-Type", "application/json")

	if params != nil {
		req.SetQueryParams(params)
	}

	if body != nil {
		req.SetBody(body)
	}

	c.Logger.Debug("Making POST request", 
		zap.String("endpoint", endpoint),
		zap.Any("params", params))

	resp, err := req.Post(endpoint)
	if err != nil {
		c.Logger.Error("POST request failed", zap.String("endpoint", endpoint), zap.Error(err))
		return nil, fmt.Errorf("POST request failed: %w", err)
	}

	c.Logger.Debug("POST request completed",
		zap.String("endpoint", endpoint),
		zap.Int("status_code", resp.StatusCode()))

	return resp, nil
}

// PatchWithParams makes a PATCH request using a map of parameters
func (c *AXMClient) PatchWithParams(endpoint string, body interface{}, params map[string]string) (*resty.Response, error) {
	return c.patch(endpoint, body, params)
}

// Patch makes a PATCH request using a QueryBuilder
func (c *AXMClient) Patch(endpoint string, body interface{}, queryBuilder *QueryBuilder) (*resty.Response, error) {
	var params map[string]string
	if queryBuilder != nil {
		params = queryBuilder.Build()
	}
	return c.patch(endpoint, body, params)
}

// patch is the internal method that handles the actual PATCH request
func (c *AXMClient) patch(endpoint string, body interface{}, params map[string]string) (*resty.Response, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	req := c.HTTP.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetHeader("Content-Type", "application/json")

	if params != nil {
		req.SetQueryParams(params)
	}

	if body != nil {
		req.SetBody(body)
	}

	c.Logger.Debug("Making PATCH request", 
		zap.String("endpoint", endpoint),
		zap.Any("params", params))

	resp, err := req.Patch(endpoint)
	if err != nil {
		c.Logger.Error("PATCH request failed", zap.String("endpoint", endpoint), zap.Error(err))
		return nil, fmt.Errorf("PATCH request failed: %w", err)
	}

	c.Logger.Debug("PATCH request completed",
		zap.String("endpoint", endpoint),
		zap.Int("status_code", resp.StatusCode()))

	return resp, nil
}

// DeleteWithParams makes a DELETE request using a map of parameters
func (c *AXMClient) DeleteWithParams(endpoint string, params map[string]string) (*resty.Response, error) {
	return c.delete(endpoint, params)
}

// Delete makes a DELETE request using a QueryBuilder
func (c *AXMClient) Delete(endpoint string, queryBuilder *QueryBuilder) (*resty.Response, error) {
	var params map[string]string
	if queryBuilder != nil {
		params = queryBuilder.Build()
	}
	return c.delete(endpoint, params)
}

// delete is the internal method that handles the actual DELETE request
func (c *AXMClient) delete(endpoint string, params map[string]string) (*resty.Response, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	req := c.HTTP.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetHeader("Content-Type", "application/json")

	if params != nil {
		req.SetQueryParams(params)
	}

	c.Logger.Debug("Making DELETE request", 
		zap.String("endpoint", endpoint),
		zap.Any("params", params))

	resp, err := req.Delete(endpoint)
	if err != nil {
		c.Logger.Error("DELETE request failed", zap.String("endpoint", endpoint), zap.Error(err))
		return nil, fmt.Errorf("DELETE request failed: %w", err)
	}

	c.Logger.Debug("DELETE request completed",
		zap.String("endpoint", endpoint),
		zap.Int("status_code", resp.StatusCode()))

	return resp, nil
}