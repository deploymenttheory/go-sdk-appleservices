package client

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// PaginatedResponse represents a response that supports pagination
type PaginatedResponse struct {
	Data  any               `json:"data"`
	Links map[string]string `json:"links"`
	Meta  any               `json:"meta,omitempty"`
}

// HasNextPage checks if there are more pages available
func (p *PaginatedResponse) HasNextPage() bool {
	next, exists := p.Links["next"]
	return exists && next != ""
}

// GetNextURL returns the next page URL
func (p *PaginatedResponse) GetNextURL() string {
	if next, exists := p.Links["next"]; exists {
		return next
	}
	return ""
}

// Get makes a GET request using a QueryBuilder and custom headers
func (c *AXMClient) Get(endpoint string, queryBuilder *QueryBuilder, headers map[string]string) (*resty.Response, []any, error) {
	var allItems []any
	currentParams := make(map[string]string)

	// Build parameters from QueryBuilder
	if queryBuilder != nil {
		params := queryBuilder.Build()
		maps.Copy(currentParams, params)
	}

	// Set default limit if not specified
	if _, hasLimit := currentParams["limit"]; !hasLimit {
		currentParams["limit"] = "200"
	}

	c.Logger.Debug("Starting GET request with automatic pagination",
		zap.String("endpoint", endpoint),
		zap.Any("params", currentParams))

	var lastResponse *resty.Response
	pageCount := 0

	for {
		pageCount++

		// Make authenticated request
		if err := c.ensureAuthenticated(); err != nil {
			return nil, nil, fmt.Errorf("authentication failed: %w", err)
		}

		req := c.HTTP.R().
			SetHeader("Authorization", "Bearer "+c.accessToken).
			SetQueryParams(currentParams)

		if headers != nil {
			req.SetHeaders(headers)
		}

		resp, err := req.Get(endpoint)
		if err != nil {
			return lastResponse, allItems, fmt.Errorf("GET request failed on page %d: %w", pageCount, err)
		}

		if !resp.IsSuccess() {
			return resp, allItems, fmt.Errorf("GET request failed with status %d", resp.StatusCode())
		}

		lastResponse = resp

		// Try to parse as paginated response
		var paginatedResp PaginatedResponse
		if err := json.Unmarshal(resp.Body(), &paginatedResp); err != nil {
			// Not a paginated response, just return the single result
			var singleResult any
			if err := json.Unmarshal(resp.Body(), &singleResult); err != nil {
				return resp, nil, fmt.Errorf("failed to unmarshal response: %w", err)
			}
			c.Logger.Debug("Single page response", zap.String("endpoint", endpoint))
			return resp, []any{singleResult}, nil
		}

		if paginatedResp.Data != nil {
			if dataSlice, ok := paginatedResp.Data.([]any); ok {
				allItems = append(allItems, dataSlice...)
				c.Logger.Debug("Retrieved page",
					zap.Int("page_number", pageCount),
					zap.Int("page_items", len(dataSlice)),
					zap.Int("total_items", len(allItems)))
			}
		}

		if !paginatedResp.HasNextPage() {
			break
		}

		// Parse next URL and update parameters
		nextURL := paginatedResp.GetNextURL()
		updatedParams, err := c.parseURLParams(nextURL)
		if err != nil {
			c.Logger.Warn("Failed to parse next URL",
				zap.String("next_url", nextURL),
				zap.Error(err))
			break
		}
		maps.Copy(currentParams, updatedParams)
	}

	c.Logger.Debug("Completed GET request",
		zap.String("endpoint", endpoint),
		zap.Int("total_pages", pageCount),
		zap.Int("total_items", len(allItems)))

	return lastResponse, allItems, nil
}

// parseURLParams extracts query parameters from a URL
func (c *AXMClient) parseURLParams(urlStr string) (map[string]string, error) {
	if urlStr == "" {
		return nil, fmt.Errorf("empty URL")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	params := make(map[string]string)
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params, nil
}
