package client

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryBuilder provides a fluent interface for building query parameters
type QueryBuilder struct {
	params map[string]string
}

// NewQueryBuilder creates a new query parameter builder
func (c *AXMClient) NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		params: make(map[string]string),
	}
}

// Fields adds field filtering to the query
func (q *QueryBuilder) Fields(resource string, fields []string) *QueryBuilder {
	if len(fields) > 0 {
		q.params[fmt.Sprintf("fields[%s]", resource)] = strings.Join(fields, ",")
	}
	return q
}

// Limit sets the number of resources to return
func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	if limit > 0 {
		q.params["limit"] = strconv.Itoa(limit)
	}
	return q
}

// Filter adds a filter parameter
func (q *QueryBuilder) Filter(key, value string) *QueryBuilder {
	if key != "" && value != "" {
		q.params[fmt.Sprintf("filter[%s]", key)] = value
	}
	return q
}

// Sort adds sorting to the query
func (q *QueryBuilder) Sort(field string) *QueryBuilder {
	if field != "" {
		q.params["sort"] = field
	}
	return q
}

// Include adds included relationships
func (q *QueryBuilder) Include(relationships []string) *QueryBuilder {
	if len(relationships) > 0 {
		q.params["include"] = strings.Join(relationships, ",")
	}
	return q
}

// Param adds a custom parameter
func (q *QueryBuilder) Param(key, value string) *QueryBuilder {
	if key != "" && value != "" {
		q.params[key] = value
	}
	return q
}

// Build returns the constructed query parameters
func (q *QueryBuilder) Build() map[string]string {
	return q.params
}