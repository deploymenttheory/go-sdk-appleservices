package client

import (
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/interfaces"
)

// QueryBuilder provides a fluent interface for building query parameters
type QueryBuilder struct {
	params map[string]string
}

// Ensure QueryBuilder implements the interface
var _ interfaces.QueryBuilder = (*QueryBuilder)(nil)
var _ interfaces.ServiceQueryBuilder = (*QueryBuilder)(nil)

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		params: make(map[string]string),
	}
}

// AddString adds a string parameter if the value is not empty
func (qb *QueryBuilder) AddString(key, value string) interfaces.QueryBuilder {
	if value != "" {
		qb.params[key] = value
	}
	return qb
}

// AddInt adds an integer parameter if the value is greater than 0
func (qb *QueryBuilder) AddInt(key string, value int) interfaces.QueryBuilder {
	if value > 0 {
		qb.params[key] = strconv.Itoa(value)
	}
	return qb
}

// AddInt64 adds an int64 parameter if the value is greater than 0
func (qb *QueryBuilder) AddInt64(key string, value int64) interfaces.QueryBuilder {
	if value > 0 {
		qb.params[key] = strconv.FormatInt(value, 10)
	}
	return qb
}

// AddBool adds a boolean parameter
func (qb *QueryBuilder) AddBool(key string, value bool) interfaces.QueryBuilder {
	qb.params[key] = strconv.FormatBool(value)
	return qb
}

// AddTime adds a time parameter in RFC3339 format if the time is not zero
func (qb *QueryBuilder) AddTime(key string, value time.Time) interfaces.QueryBuilder {
	if !value.IsZero() {
		qb.params[key] = value.Format(time.RFC3339)
	}
	return qb
}

// AddStringSlice adds a string slice parameter as comma-separated values
func (qb *QueryBuilder) AddStringSlice(key string, values []string) interfaces.QueryBuilder {
	if len(values) > 0 {
		// Join multiple values with comma
		result := ""
		for i, v := range values {
			if v != "" {
				if i > 0 {
					result += ","
				}
				result += v
			}
		}
		if result != "" {
			qb.params[key] = result
		}
	}
	return qb
}

// AddIntSlice adds an integer slice parameter as comma-separated values
func (qb *QueryBuilder) AddIntSlice(key string, values []int) interfaces.QueryBuilder {
	if len(values) > 0 {
		result := ""
		for i, v := range values {
			if i > 0 {
				result += ","
			}
			result += strconv.Itoa(v)
		}
		qb.params[key] = result
	}
	return qb
}

// AddCustom adds a custom parameter with any value
func (qb *QueryBuilder) AddCustom(key, value string) interfaces.QueryBuilder {
	qb.params[key] = value
	return qb
}

// AddIfNotEmpty adds a parameter only if the value is not empty
func (qb *QueryBuilder) AddIfNotEmpty(key, value string) interfaces.QueryBuilder {
	if value != "" {
		qb.params[key] = value
	}
	return qb
}

// AddIfTrue adds a parameter only if the condition is true
func (qb *QueryBuilder) AddIfTrue(condition bool, key, value string) interfaces.QueryBuilder {
	if condition {
		qb.params[key] = value
	}
	return qb
}

// Merge merges parameters from another query builder or map
func (qb *QueryBuilder) Merge(other map[string]string) interfaces.QueryBuilder {
	for k, v := range other {
		qb.params[k] = v
	}
	return qb
}

// Remove removes a parameter
func (qb *QueryBuilder) Remove(key string) interfaces.QueryBuilder {
	delete(qb.params, key)
	return qb
}

// Has checks if a parameter exists
func (qb *QueryBuilder) Has(key string) bool {
	_, exists := qb.params[key]
	return exists
}

// Get retrieves a parameter value
func (qb *QueryBuilder) Get(key string) string {
	return qb.params[key]
}

// Build returns the final map of query parameters
func (qb *QueryBuilder) Build() map[string]string {
	// Return a copy to prevent external modification
	result := make(map[string]string, len(qb.params))
	for k, v := range qb.params {
		result[k] = v
	}
	return result
}

// BuildString returns the query parameters as a URL-encoded string
func (qb *QueryBuilder) BuildString() string {
	if len(qb.params) == 0 {
		return ""
	}

	result := ""
	first := true
	for k, v := range qb.params {
		if !first {
			result += "&"
		}
		result += k + "=" + v
		first = false
	}
	return result
}

// Clear removes all parameters
func (qb *QueryBuilder) Clear() interfaces.QueryBuilder {
	qb.params = make(map[string]string)
	return qb
}

// Count returns the number of parameters
func (qb *QueryBuilder) Count() int {
	return len(qb.params)
}

// IsEmpty returns true if no parameters are set
func (qb *QueryBuilder) IsEmpty() bool {
	return len(qb.params) == 0
}
