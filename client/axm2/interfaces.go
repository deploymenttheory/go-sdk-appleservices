package axm2

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-apple/client/shared"
)

// OrgDevicesService provides operations for the /v1/orgDevices endpoint
type OrgDevicesService interface {
	// Device listing and details
	GetOrgDevices(ctx context.Context, opts ...interface{}) ([]OrgDevice, error)
	GetOrgDevice(ctx context.Context, deviceID string, opts ...interface{}) (*OrgDevice, error)

	// Device assignment relationships
	GetAssignedMdmServer(ctx context.Context, deviceID string, opts ...interface{}) (string, error)
	GetAssignedMdmServerInfo(ctx context.Context, deviceID string, opts ...interface{}) (*MdmServer, error)
}

// MdmServersService provides operations for the /v1/mdmServers endpoint
type MdmServersService interface {
	// MDM server listing and details
	GetMdmServers(ctx context.Context, opts ...interface{}) ([]MdmServer, error)
	GetMdmServer(ctx context.Context, serverID string, opts ...interface{}) (*MdmServer, error)

	// MDM server device relationships
	GetDevices(ctx context.Context, serverID string, opts ...interface{}) ([]string, error)
}

// OrgDeviceActivitiesService provides operations for the /v1/orgDeviceActivities endpoint
type OrgDeviceActivitiesService interface {
	// Activity tracking (read-only, for past 30 days)
	GetActivity(ctx context.Context, activityID string, opts ...interface{}) (*OrgDeviceActivity, error)

	// Device assignment operations (POST to /v1/orgDeviceActivities)
	AssignDevice(ctx context.Context, deviceID, mdmServerID string) (*OrgDeviceActivity, error)
	UnassignDevice(ctx context.Context, deviceID string) (*OrgDeviceActivity, error)
	AssignDevices(ctx context.Context, deviceIDs []string, mdmServerID string) (*OrgDeviceActivity, error)
	UnassignDevices(ctx context.Context, deviceIDs []string) (*OrgDeviceActivity, error)
}

// AXMClient defines the main interface for Apple School/Business Manager API client
// organized by Apple API endpoints
type AXMClient interface {
	// Apple API endpoint services
	OrgDevices() OrgDevicesService
	MdmServers() MdmServersService
	OrgDeviceActivities() OrgDeviceActivitiesService

	// Authentication
	IsAuthenticated() bool
	ForceReauthenticate() error

	// Client information
	GetClientID() string
	GetBaseURL() string
	GetAPIType() string

	// Health and diagnostics
	HealthCheck(ctx context.Context) error
	GetDiagnostics() map[string]any

	// Lifecycle
	Close()
}

// TokenProviderInterface defines the interface for token management
type TokenProviderInterface interface {
	GetToken(ctx context.Context) (string, error)
	IsValid() bool
	ForceRefresh(ctx context.Context) error
}

// RequestOption defines functional options for configuring requests
type RequestOption func(*RequestBuilder)

// RequestBuilder provides context for applying request options
type RequestBuilder struct {
	req RequestInterface
}

// RequestInterface represents a request interface that services can use
// Use shared interfaces
type RequestInterface = shared.RequestInterface
type ResponseInterface = shared.ResponseInterface
type HTTPClientInterface = shared.HTTPClientInterface

// WithFieldsOption creates a RequestOption for field filtering
func WithFieldsOption(resource string, fields []string) RequestOption {
	return func(rb *RequestBuilder) {
		fieldsParam := strings.Join(fields, ",")
		rb.req.SetQueryParam(fmt.Sprintf("fields[%s]", resource), fieldsParam)
	}
}

// WithLimitOption creates a RequestOption for pagination limit
func WithLimitOption(limit int) RequestOption {
	return func(rb *RequestBuilder) {
		rb.req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

// LoadBalancerInterface defines the interface for custom load balancing
type LoadBalancerInterface interface {
	Next() string
	Feedback(url string, success bool)
	Close()
}
