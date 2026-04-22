package client

import (
	"context"

	"go.uber.org/zap"
)

// Client is the interface service implementations depend on.
// The Transport struct in this package satisfies this interface.
type Client interface {
	// NewRequest returns a RequestBuilder for constructing API requests.
	// Retry and logging are applied by the transport at execution time.
	NewRequest(ctx context.Context) *RequestBuilder

	// QueryBuilder returns a new query parameter builder instance.
	QueryBuilder() *QueryBuilder

	// GetLogger returns the configured zap logger instance.
	GetLogger() *zap.Logger
}
