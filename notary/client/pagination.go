package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Meta contains pagination metadata matching Apple's API format.
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

// Paging contains pagination information matching Apple's API format.
type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Links contains pagination navigation links matching Apple's API format.
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// PaginationOptions represents common pagination parameters for Apple's API.
type PaginationOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// AddToQueryBuilder adds pagination options to a query builder.
func (opts *PaginationOptions) AddToQueryBuilder(qb *QueryBuilder) *QueryBuilder {
	if opts == nil {
		return qb
	}

	return qb.
		AddInt("limit", opts.Limit).
		AddString("cursor", opts.Cursor)
}

// HasNextPage checks if there is a next page available.
func HasNextPage(links *Links) bool {
	return links != nil && links.Next != ""
}

// HasPrevPage checks if there is a previous page available.
func HasPrevPage(links *Links) bool {
	return links != nil && links.Prev != ""
}

// extractParamsFromURL extracts query parameters from a URL string.
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

// parseJSON is a helper that unmarshals raw JSON bytes into a target value.
func parseJSON(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}
	return nil
}
