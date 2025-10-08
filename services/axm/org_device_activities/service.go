package org_device_activities

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Constants for API endpoints
const (
	OrgDeviceActivitiesEndpoint = "/v1/orgDeviceActivities"
)

// Service handles org device activities operations
type Service struct {
	client HTTPClient
	logger *zap.Logger
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Get(ctx context.Context, endpoint string, result any, opts ...any) error
	Post(ctx context.Context, endpoint string, body, result any, opts ...any) error
}

// NewService creates a new org device activities service
func NewService(client HTTPClient, logger *zap.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}

// GetActivity retrieves information about a specific organization device activity
// This gets activity information for device management actions that were performed
// Activities are only available for the past 30 days
// Apple API endpoint: GET /v1/orgDeviceActivities/{id}
func (s *Service) GetActivity(ctx context.Context, activityID string, opts ...any) (*OrgDeviceActivity, error) {
	s.logger.Debug("Getting organization device activity", zap.String("activity_id", activityID))

	if activityID == "" {
		return nil, fmt.Errorf("activity ID cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/%s", OrgDeviceActivitiesEndpoint, activityID)

	var activityResponse OrgDeviceActivityResponse
	if err := s.client.Get(ctx, endpoint, &activityResponse, opts...); err != nil {
		return nil, fmt.Errorf("failed to get organization device activity: %w", err)
	}

	s.logger.Debug("Successfully retrieved organization device activity",
		zap.String("activity_id", activityID),
		zap.String("status", activityResponse.Data.Attributes.Status),
		zap.String("sub_status", activityResponse.Data.Attributes.SubStatus))

	return &activityResponse.Data, nil
}

// buildAssignDevicesRequest creates a request for assigning devices to an MDM server
func buildAssignDevicesRequest(deviceIDs []string, mdmServerID string) *OrgDeviceActivitiesRequest {
	deviceLinkages := make([]ResourceLinkage, len(deviceIDs))
	for i, deviceID := range deviceIDs {
		deviceLinkages[i] = ResourceLinkage{
			Type: "orgDevices",
			ID:   deviceID,
		}
	}

	attributes := OrgDeviceActivitiesRequestAttributes{
		ActivityType: ActivityTypeAssignDevices,
	}

	mdmServerRelationship := &MdmServerRelationship{
		Data: ResourceLinkage{
			Type: "mdmServers",
			ID:   mdmServerID,
		},
	}

	devicesRelationship := &DevicesRelationship{
		Data: deviceLinkages,
	}

	relationships := OrgDeviceActivitiesRequestRelationships{
		MdmServer: mdmServerRelationship,
		Devices:   devicesRelationship,
	}

	data := OrgDeviceActivitiesRequestData{
		Type:          "orgDeviceActivities",
		Attributes:    attributes,
		Relationships: relationships,
	}

	return &OrgDeviceActivitiesRequest{
		Data: data,
	}
}

// AssignDevices assigns devices to a device management service
func (s *Service) AssignDevices(ctx context.Context, deviceIDs []string, mdmServerID string) (*OrgDeviceActivity, error) {
	s.logger.Debug("Assigning devices to MDM server",
		zap.Strings("device_ids", deviceIDs),
		zap.String("mdm_server_id", mdmServerID))

	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("device IDs cannot be empty")
	}
	if mdmServerID == "" {
		return nil, fmt.Errorf("MDM server ID cannot be empty for assignment")
	}

	request := buildAssignDevicesRequest(deviceIDs, mdmServerID)

	var activityResponse OrgDeviceActivityResponse
	if err := s.client.Post(ctx, OrgDeviceActivitiesEndpoint, request, &activityResponse); err != nil {
		return nil, fmt.Errorf("failed to assign devices: %w", err)
	}

	s.logger.Info("Successfully created device assignment activity",
		zap.String("activity_id", activityResponse.Data.ID),
		zap.String("status", activityResponse.Data.Attributes.Status),
		zap.String("sub_status", activityResponse.Data.Attributes.SubStatus),
		zap.Int("device_count", len(deviceIDs)))

	return &activityResponse.Data, nil
}

// buildUnassignDevicesRequest creates a request for unassigning devices
func buildUnassignDevicesRequest(deviceIDs []string) *OrgDeviceActivitiesRequest {
	deviceLinkages := make([]ResourceLinkage, len(deviceIDs))
	for i, deviceID := range deviceIDs {
		deviceLinkages[i] = ResourceLinkage{
			Type: "orgDevices",
			ID:   deviceID,
		}
	}

	attributes := OrgDeviceActivitiesRequestAttributes{
		ActivityType: ActivityTypeUnassignDevices,
	}

	devicesRelationship := &DevicesRelationship{
		Data: deviceLinkages,
	}

	relationships := OrgDeviceActivitiesRequestRelationships{
		Devices: devicesRelationship,
	}

	data := OrgDeviceActivitiesRequestData{
		Type:          "orgDeviceActivities",
		Attributes:    attributes,
		Relationships: relationships,
	}

	return &OrgDeviceActivitiesRequest{
		Data: data,
	}
}

// UnassignDevices unassigns devices from their current device management service
func (s *Service) UnassignDevices(ctx context.Context, deviceIDs []string) (*OrgDeviceActivity, error) {
	s.logger.Debug("Unassigning devices from MDM server",
		zap.Strings("device_ids", deviceIDs))

	if len(deviceIDs) == 0 {
		return nil, fmt.Errorf("device IDs cannot be empty")
	}

	request := buildUnassignDevicesRequest(deviceIDs)

	var activityResponse OrgDeviceActivityResponse
	if err := s.client.Post(ctx, OrgDeviceActivitiesEndpoint, request, &activityResponse); err != nil {
		return nil, fmt.Errorf("failed to unassign devices: %w", err)
	}

	s.logger.Info("Successfully created device unassignment activity",
		zap.String("activity_id", activityResponse.Data.ID),
		zap.String("status", activityResponse.Data.Attributes.Status),
		zap.String("sub_status", activityResponse.Data.Attributes.SubStatus),
		zap.Int("device_count", len(deviceIDs)))

	return &activityResponse.Data, nil
}

// AssignDevice assigns a single device to a device management service
func (s *Service) AssignDevice(ctx context.Context, deviceID, mdmServerID string) (*OrgDeviceActivity, error) {
	return s.AssignDevices(ctx, []string{deviceID}, mdmServerID)
}

// UnassignDevice unassigns a single device from its current device management service
func (s *Service) UnassignDevice(ctx context.Context, deviceID string) (*OrgDeviceActivity, error) {
	return s.UnassignDevices(ctx, []string{deviceID})
}
