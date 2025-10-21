package axm

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
)

// Meta contains pagination metadata matching Apple's API format
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

// Paging contains pagination information matching Apple's API format
type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Links contains pagination navigation links matching Apple's API format
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// PaginationOptions represents common pagination parameters for Apple's API
type PaginationOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// AddToQueryBuilder adds pagination options to a query builder
func (opts *PaginationOptions) AddToQueryBuilder(qb ServiceQueryBuilder) ServiceQueryBuilder {
	if opts == nil {
		return qb
	}

	return qb.
		AddInt("limit", opts.Limit).
		AddString("cursor", opts.Cursor)
}

// GetPaginated executes a paginated GET request
func (c *Client) GetPaginated(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error {
	return c.Get(ctx, path, queryParams, headers, result)
}

// GetNextPage extracts the next page URL from links and makes a request
func (c *Client) GetNextPage(ctx context.Context, nextURL string, headers map[string]string, result any) error {
	if nextURL == "" {
		return fmt.Errorf("no next page URL provided")
	}

	parsedURL, err := url.Parse(nextURL)
	if err != nil {
		return fmt.Errorf("failed to parse next URL: %w", err)
	}

	queryParams := make(map[string]string)
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// Make the request using the path and extracted parameters
	return c.Get(ctx, parsedURL.Path, queryParams, headers, result)
}

// HasNextPage checks if there is a next page available
func HasNextPage(links *Links) bool {
	return links != nil && links.Next != ""
}

// HasPrevPage checks if there is a previous page available
func HasPrevPage(links *Links) bool {
	return links != nil && links.Prev != ""
}

// GetAllPages retrieves all pages of results by following pagination links
func (c *Client) GetAllPages(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, processPage func([]byte) error) error {
	currentParams := make(map[string]string)
	maps.Copy(currentParams, queryParams)

	for {
		var rawResponse json.RawMessage
		err := c.GetPaginated(ctx, path, currentParams, headers, &rawResponse)
		if err != nil {
			return err
		}

		if err := processPage(rawResponse); err != nil {
			return err
		}

		var pageInfo struct {
			Links *Links `json:"links,omitempty"`
		}
		if err := json.Unmarshal(rawResponse, &pageInfo); err != nil {
			return fmt.Errorf("failed to parse pagination info: %w", err)
		}

		if !HasNextPage(pageInfo.Links) {
			break
		}

		nextParams, err := extractParamsFromURL(pageInfo.Links.Next)
		if err != nil {
			return fmt.Errorf("failed to parse next URL: %w", err)
		}

		maps.Copy(currentParams, nextParams)
	}

	return nil
}

// extractParamsFromURL extracts query parameters from a URL string
func extractParamsFromURL(urlStr string) (map[string]string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params, nil
}