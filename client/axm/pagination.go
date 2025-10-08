package axm

import (
	"context"
	"fmt"
	"net/url"
)

// PagedResponse represents a paginated API response
type PagedResponse[T any] struct {
	Data   []T   `json:"data"`
	Meta   Meta  `json:"meta"`
	Links  Links `json:"links"`
	client *Client
	path   string
	params map[string]string
}

// Meta contains pagination metadata
type Meta struct {
	TotalCount int `json:"total_count,omitempty"`
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	PageCount  int `json:"page_count,omitempty"`
}

// Links contains pagination navigation links
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

// IteratorResult represents a single item in pagination iteration
type IteratorResult[T any] struct {
	Item T
	Err  error
}

// HasNext returns true if there is a next page available
func (pr *PagedResponse[T]) HasNext() bool {
	return pr.Links.Next != ""
}

// HasPrev returns true if there is a previous page available
func (pr *PagedResponse[T]) HasPrev() bool {
	return pr.Links.Prev != ""
}

// Next retrieves the next page of results
func (pr *PagedResponse[T]) Next(ctx context.Context) (*PagedResponse[T], error) {
	if !pr.HasNext() {
		return nil, ErrNoNextPage
	}

	// Extract parameters from next URL
	nextParams, err := extractParamsFromURL(pr.Links.Next)
	if err != nil {
		return nil, fmt.Errorf("failed to parse next URL: %w", err)
	}

	// Merge with existing params
	params := make(map[string]string)
	for k, v := range pr.params {
		params[k] = v
	}
	for k, v := range nextParams {
		params[k] = v
	}

	// TODO: Implement using the new client architecture
	return nil, fmt.Errorf("pagination method not implemented with new client structure")
}

// Prev retrieves the previous page of results
func (pr *PagedResponse[T]) Prev(ctx context.Context) (*PagedResponse[T], error) {
	if !pr.HasPrev() {
		return nil, fmt.Errorf("no previous page available")
	}

	// Extract parameters from prev URL
	prevParams, err := extractParamsFromURL(pr.Links.Prev)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prev URL: %w", err)
	}

	// Merge with existing params
	params := make(map[string]string)
	for k, v := range pr.params {
		params[k] = v
	}
	for k, v := range prevParams {
		params[k] = v
	}

	// TODO: Implement using the new client architecture
	return nil, fmt.Errorf("pagination method not implemented with new client structure")
}

// Iterator returns a channel for iterating through all pages
func (pr *PagedResponse[T]) Iterator(ctx context.Context) <-chan IteratorResult[T] {
	ch := make(chan IteratorResult[T])

	go func() {
		defer close(ch)

		current := pr
		for {
			// Send current page items
			for _, item := range current.Data {
				select {
				case ch <- IteratorResult[T]{Item: item}:
				case <-ctx.Done():
					ch <- IteratorResult[T]{Err: ctx.Err()}
					return
				}
			}

			// Check if there's a next page
			if !current.HasNext() {
				return
			}

			// Get next page
			next, err := current.Next(ctx)
			if err != nil {
				ch <- IteratorResult[T]{Err: err}
				return
			}
			current = next
		}
	}()

	return ch
}

// CollectAll collects all items from all pages into a single slice
func (pr *PagedResponse[T]) CollectAll(ctx context.Context) ([]T, error) {
	var allItems []T

	current := pr
	for {
		allItems = append(allItems, current.Data...)

		if !current.HasNext() {
			break
		}

		next, err := current.Next(ctx)
		if err != nil {
			return nil, err
		}
		current = next
	}

	return allItems, nil
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

// PaginationOptions represents common pagination parameters
type PaginationOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Page   int    `json:"page,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// AddToQueryBuilder adds pagination options to a query builder
func (opts *PaginationOptions) AddToQueryBuilder(qb *QueryBuilder) *QueryBuilder {
	if opts == nil {
		return qb
	}

	return qb.
		AddInt("limit", opts.Limit).
		AddString("cursor", opts.Cursor).
		AddInt("page", opts.Page).
		AddInt("offset", opts.Offset)
}
