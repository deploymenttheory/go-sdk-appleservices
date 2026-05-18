package client

import (
	"context"

	"go.uber.org/zap"
)

// Client is the interface service implementations depend on.
// The Transport struct in this package satisfies this interface.
type Client interface {
	// NewRequest returns a RequestBuilder that the service layer uses to
	// construct a complete request — headers, body, query params, result
	// target — before executing it via Get/Post/Put/Patch/Delete/GetPaginated.
	// Auth, retry, and concurrency limiting are applied by the transport at
	// execution time.
	NewRequest(ctx context.Context) *RequestBuilder

	// QueryBuilder returns a new query parameter builder instance.
	// Use this to build complex query parameter sets before passing
	// them to SetQueryParams on the RequestBuilder.
	QueryBuilder() *QueryBuilder

	// GetLogger returns the configured zap logger instance.
	GetLogger() *zap.Logger
}
