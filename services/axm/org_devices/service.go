package org_devices

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/client/shared"
	"go.uber.org/zap"
)

// Constants for API endpoints
const (
	OrgDevicesEndpoint = "/v1/orgDevices"
)

// Service handles org devices operations
type Service struct {
	client HTTPClient
	logger *zap.Logger
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	DoRequestWithPagination(ctx context.Context, endpoint string, newResponseFunc func() shared.PaginatedResponse, opts ...interface{}) (interface{}, error)
	GetHTTPClient() shared.HTTPClientInterface
	ApplyRequestOptions(req shared.RequestInterface, opts ...interface{})
}

// NewService creates a new org devices service
func NewService(client HTTPClient, logger *zap.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}

// GetOrgDevices retrieves organization devices with automatic pagination using centralized helper
func (s *Service) GetOrgDevices(ctx context.Context, opts ...interface{}) ([]OrgDevice, error) {
	s.logger.Debug("Getting organization devices with centralized pagination")

	// Use centralized pagination helper
	result, err := s.client.DoRequestWithPagination(ctx, OrgDevicesEndpoint, func() shared.PaginatedResponse {
		return &OrgDevicesResponse{}
	}, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to get organization devices: %w", err)
	}

	devices, ok := result.([]OrgDevice)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from paginated request")
	}

	s.logger.Info("Successfully retrieved organization devices",
		zap.Int("device_count", len(devices)))

	return devices, nil
}

// GetOrgDevice retrieves a single organization device by ID using Resty v3 patterns
func (s *Service) GetOrgDevice(ctx context.Context, deviceID string, opts ...interface{}) (*OrgDevice, error) {
	s.logger.Debug("Getting organization device", zap.String("device_id", deviceID))

	if deviceID == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s", OrgDevicesEndpoint, deviceID)

	// Use Resty v3 pattern with SetResult for automatic unmarshaling
	var deviceResponse OrgDeviceResponse
	var errorResponse APIError

	request := s.client.GetHTTPClient().R().SetContext(ctx).
		SetResult(&deviceResponse).
		SetError(&errorResponse)

	// Apply RequestOption parameters via client
	s.client.ApplyRequestOptions(request, opts...)

	response, err := request.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization device: %w", err)
	}

	if response.IsError() {
		s.logger.Error("API error getting device",
			zap.String("device_id", deviceID),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", errorResponse))
		return nil, fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	s.logger.Debug("Successfully retrieved organization device",
		zap.String("device_id", deviceID),
		zap.String("serial_number", deviceResponse.Data.Attributes.SerialNumber))

	return &deviceResponse.Data, nil
}

// GetAssignedMdmServer retrieves the assigned MDM server ID for a specific device
// Apple API endpoint: GET /v1/orgDevices/{id}/relationships/assignedServer
func (s *Service) GetAssignedMdmServer(ctx context.Context, deviceID string, opts ...interface{}) (string, error) {
	s.logger.Debug("Getting assigned MDM server for device", zap.String("device_id", deviceID))

	if deviceID == "" {
		return "", fmt.Errorf("device ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s/relationships/assignedServer", OrgDevicesEndpoint, deviceID)

	var linkageResponse OrgDeviceAssignedServerLinkageResponse
	var errorResponse APIError

	request := s.client.GetHTTPClient().R().
		SetContext(ctx).
		SetResult(&linkageResponse).
		SetError(&errorResponse)

	// Apply RequestOption parameters via client
	s.client.ApplyRequestOptions(request, opts...)

	response, err := request.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get device assigned MDM server: %w", err)
	}

	if response.IsError() {
		s.logger.Error("API error getting device assigned MDM server",
			zap.String("device_id", deviceID),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", errorResponse))
		return "", fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	if linkageResponse.Data == nil {
		s.logger.Debug("Device has no assigned MDM server", zap.String("device_id", deviceID))
		return "", nil
	}

	if linkageResponse.Data.Type != "mdmServers" {
		s.logger.Warn("Unexpected assigned server type",
			zap.String("device_id", deviceID),
			zap.String("expected_type", "mdmServers"),
			zap.String("actual_type", linkageResponse.Data.Type))
	}

	mdmServerID := linkageResponse.Data.ID
	s.logger.Debug("Successfully retrieved device assigned MDM server",
		zap.String("device_id", deviceID),
		zap.String("mdm_server_id", mdmServerID),
		zap.String("self_link", linkageResponse.Links.Self),
		zap.String("related_link", linkageResponse.Links.Related))

	return mdmServerID, nil
}

// GetAssignedMdmServerInfo retrieves the assigned device management service information for a device
// Apple API endpoint: GET /v1/orgDevices/{id}/assignedServer
func (s *Service) GetAssignedMdmServerInfo(ctx context.Context, deviceID string, opts ...interface{}) (*MdmServer, error) {
	s.logger.Debug("Getting assigned MDM server info for device", zap.String("device_id", deviceID))

	if deviceID == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s/assignedServer", OrgDevicesEndpoint, deviceID)

	var serverResponse MdmServerResponse
	var errorResponse APIError

	request := s.client.GetHTTPClient().R().
		SetContext(ctx).
		SetResult(&serverResponse).
		SetError(&errorResponse)

	// Apply RequestOption parameters via client
	s.client.ApplyRequestOptions(request, opts...)

	response, err := request.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get device assigned MDM server info: %w", err)
	}

	if response.IsError() {
		s.logger.Error("API error getting device assigned MDM server info",
			zap.String("device_id", deviceID),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", errorResponse))
		return nil, fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	s.logger.Debug("Successfully retrieved device assigned MDM server info",
		zap.String("device_id", deviceID),
		zap.String("mdm_server_id", serverResponse.Data.ID),
		zap.String("server_name", serverResponse.Data.Attributes.ServerName))

	return &serverResponse.Data, nil
}
