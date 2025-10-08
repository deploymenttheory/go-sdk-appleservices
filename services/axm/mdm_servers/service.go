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
	Get(ctx context.Context, endpoint string, result any, opts ...any) error
	GetWithPagination(ctx context.Context, endpoint string, newResponseFunc func() shared.PaginatedResponse, opts ...any) (any, error)
}

// NewService creates a new MDM servers service
func NewService(client HTTPClient, logger *zap.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}

// GetMdmServers retrieves a list of device management services in the organization
func (s *Service) GetMdmServers(ctx context.Context, opts ...any) ([]MdmServer, error) {
	s.logger.Debug("Getting MDM servers with centralized pagination")

	// Use centralized pagination helper
	result, err := s.client.GetWithPagination(ctx, MdmServersEndpoint, func() shared.PaginatedResponse {
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
func (s *Service) GetMdmServer(ctx context.Context, serverID string, opts ...any) (*MdmServer, error) {
	s.logger.Debug("Getting MDM server", zap.String("server_id", serverID))

	if serverID == "" {
		return nil, fmt.Errorf("MDM server ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s", MdmServersEndpoint, serverID)

	var serverResponse MdmServerResponse
	if err := s.client.Get(ctx, endpoint, &serverResponse, opts...); err != nil {
		return nil, fmt.Errorf("failed to get MDM server: %w", err)
	}

	s.logger.Debug("Successfully retrieved MDM server",
		zap.String("server_id", serverID),
		zap.String("server_name", serverResponse.Data.Attributes.ServerName))

	return &serverResponse.Data, nil
}

// GetDevices retrieves device serial numbers assigned to a specific MDM server
func (s *Service) GetDevices(ctx context.Context, serverID string, opts ...any) ([]string, error) {
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
		
		// Apply options only for first page
		var requestOpts []any
		if pageCount == 1 {
			requestOpts = opts
		}
		
		if err := s.client.Get(ctx, nextURL, &pageResponse, requestOpts...); err != nil {
			return nil, fmt.Errorf("failed to get MDM server devices (page %d): %w", pageCount, err)
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
