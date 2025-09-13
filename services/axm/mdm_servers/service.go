package mdm_servers

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/client/shared"
	"go.uber.org/zap"
)

// Constants for API endpoints
const (
	MdmServersEndpoint = "/v1/mdmServers"
)

// Service handles MDM servers operations
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

// NewService creates a new MDM servers service
func NewService(client HTTPClient, logger *zap.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}

// GetMdmServers retrieves a list of device management services in the organization
func (s *Service) GetMdmServers(ctx context.Context, opts ...interface{}) ([]MdmServer, error) {
	s.logger.Debug("Getting MDM servers with centralized pagination")

	// Use centralized pagination helper
	result, err := s.client.DoRequestWithPagination(ctx, MdmServersEndpoint, func() shared.PaginatedResponse {
		return &MdmServersResponse{}
	}, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to get MDM servers: %w", err)
	}

	servers, ok := result.([]MdmServer)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from paginated request")
	}

	s.logger.Info("Successfully retrieved MDM servers",
		zap.Int("server_count", len(servers)))

	return servers, nil
}

// GetMdmServer retrieves a specific MDM server by ID
func (s *Service) GetMdmServer(ctx context.Context, serverID string, opts ...interface{}) (*MdmServer, error) {
	s.logger.Debug("Getting MDM server", zap.String("server_id", serverID))

	if serverID == "" {
		return nil, fmt.Errorf("MDM server ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s", MdmServersEndpoint, serverID)

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
		return nil, fmt.Errorf("failed to get MDM server: %w", err)
	}

	if response.IsError() {
		s.logger.Error("API error getting MDM server",
			zap.String("server_id", serverID),
			zap.Int("status_code", response.StatusCode()),
			zap.Any("error", errorResponse))
		return nil, fmt.Errorf("API error %d: %s", response.StatusCode(), response.String())
	}

	s.logger.Debug("Successfully retrieved MDM server",
		zap.String("server_id", serverID),
		zap.String("server_name", serverResponse.Data.Attributes.ServerName))

	return &serverResponse.Data, nil
}

// GetDevices retrieves device serial numbers assigned to a specific MDM server
func (s *Service) GetDevices(ctx context.Context, serverID string, opts ...interface{}) ([]string, error) {
	s.logger.Debug("Getting devices for MDM server", zap.String("server_id", serverID))

	if serverID == "" {
		return nil, fmt.Errorf("MDM server ID cannot be empty")
	}

	var allDeviceIDs []string
	nextURL := fmt.Sprintf("%s/%s/relationships/devices", MdmServersEndpoint, serverID)

	pageCount := 0
	for nextURL != "" {
		pageCount++
		s.logger.Debug("Fetching MDM server devices page",
			zap.String("server_id", serverID),
			zap.Int("page", pageCount),
			zap.String("url", nextURL))

		var pageResponse MdmServerDevicesLinkagesResponse
		var apiError APIError

		request := s.client.GetHTTPClient().R().
			SetContext(ctx).
			SetResult(&pageResponse).
			SetError(&apiError)

		// Apply RequestOption parameters for first page only
		if pageCount == 1 {
			s.client.ApplyRequestOptions(request, opts...)
		}

		response, err := request.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to execute GET request for MDM server devices (page %d): %w", pageCount, err)
		}

		if response.IsError() {
			s.logger.Error("API error getting MDM server devices",
				zap.String("server_id", serverID),
				zap.Int("page", pageCount),
				zap.Int("status_code", response.StatusCode()),
				zap.Any("error", apiError))
			return nil, fmt.Errorf("API error getting MDM server devices (page %d): %d %s", pageCount, response.StatusCode(), response.String())
		}

		// Extract device IDs from linkages
		for _, linkage := range pageResponse.Data {
			if linkage.Type == "orgDevices" {
				allDeviceIDs = append(allDeviceIDs, linkage.ID)
			}
		}

		nextURL = pageResponse.Links.Next

		s.logger.Debug("MDM server devices page fetched successfully",
			zap.String("server_id", serverID),
			zap.Int("page", pageCount),
			zap.Int("items_this_page", len(pageResponse.Data)),
			zap.Int("total_devices_so_far", len(allDeviceIDs)),
			zap.String("next_url", nextURL))
	}

	s.logger.Info("Successfully retrieved MDM server devices",
		zap.String("server_id", serverID),
		zap.Int("total_pages", pageCount),
		zap.Int("device_count", len(allDeviceIDs)))

	return allDeviceIDs, nil
}
