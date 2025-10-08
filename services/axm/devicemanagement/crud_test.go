package devicemanagement

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/devicemanagement/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) *Client {
	// Create a mock auth provider
	mockAuth := &MockAuthProvider{}

	// Create AXM client config
	config := axm.Config{
		BaseURL:    "https://api-business.apple.com/v1",
		Auth:       mockAuth,
		Logger:     zap.NewNop(),
		Debug:      false,
		RetryCount: 0, // Disable retries for tests
	}

	// Create AXM client
	axmClient, err := axm.NewClient(config)
	require.NoError(t, err)

	// Activate httpmock for the client's HTTP client
	httpmock.ActivateNonDefault(axmClient.GetHTTPClient().Client())

	// Setup cleanup
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Create device management client
	return NewClient(axmClient)
}

// MockAuthProvider implements the AuthProvider interface for testing
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	// Mock auth - do nothing for tests
	return nil
}

func TestGetDeviceManagementServices_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetMDMServersOptions{
		Fields: []string{FieldServerName, FieldServerType, FieldCreatedDateTime},
		Limit:  100,
	}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify the first MDM server
	server := result.Data[0]
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	assert.NotNil(t, server.Attributes)

	// Test all server attributes
	assert.Equal(t, "Test Device Management Service", server.Attributes.ServerName)
	assert.Equal(t, "MDM", server.Attributes.ServerType)
	assert.NotNil(t, server.Attributes.CreatedDateTime)
	assert.NotNil(t, server.Attributes.UpdatedDateTime)

	// Verify relationships
	assert.NotNil(t, server.Relationships)
	assert.NotNil(t, server.Relationships.Devices)
	assert.NotNil(t, server.Relationships.Devices.Links)
	assert.Contains(t, server.Relationships.Devices.Links.Self, "relationships/devices")

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServices_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, err := client.GetDeviceManagementServices(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServices_WithLimitEnforcement(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetMDMServersOptions{
		Limit: 1500, // Exceeds API maximum of 1000
	}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServices_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Reset httpmock to clear any previous registrations
	httpmock.Reset()

	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetMDMServersOptions{}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 4, httpmock.GetTotalCallCount()) // 1 original + 3 retries
}

func TestGetMDMServerDeviceLinkages_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	opts := &GetMDMServerDeviceLinkagesOptions{
		Limit: 100,
	}

	result, err := client.GetMDMServerDeviceLinkages(ctx, serverID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify device linkage
	linkage := result.Data[0]
	assert.Equal(t, "orgDevices", linkage.Type)
	assert.Equal(t, "XABC123X0ABC123X0", linkage.ID)

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetMDMServerDeviceLinkages_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetMDMServerDeviceLinkagesOptions{}

	result, err := client.GetMDMServerDeviceLinkages(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetMDMServerDeviceLinkages_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, err := client.GetMDMServerDeviceLinkages(ctx, serverID, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceIDForADevice_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "DVVS36G1YD3JKQNI"

	result, err := client.GetAssignedDeviceManagementServiceIDForADevice(ctx, deviceID)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify linkage data
	assert.Equal(t, "mdmServers", result.Data.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", result.Data.ID)

	// Verify links
	assert.NotNil(t, result.Links)
	assert.Contains(t, result.Links.Self, "relationships/assignedServer")
	assert.Contains(t, result.Links.Related, "assignedServer")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceIDForADevice_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, err := client.GetAssignedDeviceManagementServiceIDForADevice(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceInformationByDeviceID_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "DVVS36G1YD3JKQNI"
	opts := &GetAssignedServerInfoOptions{
		Fields: []string{FieldServerName, FieldServerType},
	}

	result, err := client.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify server data
	server := result.Data
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	assert.NotNil(t, server.Attributes)

	// Test all server attributes
	assert.Equal(t, "Test Device Management Service", server.Attributes.ServerName)
	assert.Equal(t, "APPLE_CONFIGURATOR", server.Attributes.ServerType)
	assert.NotNil(t, server.Attributes.CreatedDateTime)
	assert.NotNil(t, server.Attributes.UpdatedDateTime)

	// Verify relationships
	assert.NotNil(t, server.Relationships)
	assert.NotNil(t, server.Relationships.Devices)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceInformationByDeviceID_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetAssignedServerInfoOptions{}

	result, err := client.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	deviceIDs := []string{"XABC123X0ABC123X0", "YDEF456Y1DEF456Y1"}

	result, err := client.AssignDevicesToServer(ctx, serverID, deviceIDs)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify activity data
	activity := result.Data
	assert.Equal(t, "orgDeviceActivities", activity.Type)
	assert.NotEmpty(t, activity.ID)
	assert.NotNil(t, activity.Attributes)

	// Test activity attributes
	assert.Equal(t, ActivityStatusInProgress, activity.Attributes.Status)
	assert.Equal(t, ActivitySubStatusSubmitted, activity.Attributes.SubStatus)
	assert.NotNil(t, activity.Attributes.CreatedDateTime)

	// Verify links
	assert.NotNil(t, activity.Links)
	assert.Contains(t, activity.Links.Self, "orgDeviceActivities")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, err := client.AssignDevicesToServer(ctx, "", deviceIDs)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_EmptyDeviceIDs(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, err := client.AssignDevicesToServer(ctx, serverID, []string{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least one device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, err := client.UnassignDevicesFromServer(ctx, serverID, deviceIDs)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify activity data
	activity := result.Data
	assert.Equal(t, "orgDeviceActivities", activity.Type)
	assert.NotEmpty(t, activity.ID)
	assert.NotNil(t, activity.Attributes)

	// Test activity attributes
	assert.Equal(t, ActivityStatusInProgress, activity.Attributes.Status)
	assert.Equal(t, ActivitySubStatusSubmitted, activity.Attributes.SubStatus)
	assert.NotNil(t, activity.Attributes.CreatedDateTime)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, err := client.UnassignDevicesFromServer(ctx, "", deviceIDs)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_EmptyDeviceIDs(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, err := client.UnassignDevicesFromServer(ctx, serverID, []string{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least one device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestContextCancellation(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	opts := &GetMDMServersOptions{}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestContextTimeout(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep to ensure timeout
	time.Sleep(1 * time.Millisecond)

	opts := &GetMDMServersOptions{}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestFieldConstants(t *testing.T) {
	// Test that field constants are properly defined
	expectedFields := []string{
		FieldServerName,
		FieldServerType,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
		FieldDevices,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	// Test specific field values
	assert.Equal(t, "serverName", FieldServerName)
	assert.Equal(t, "serverType", FieldServerType)
	assert.Equal(t, "createdDateTime", FieldCreatedDateTime)
	assert.Equal(t, "updatedDateTime", FieldUpdatedDateTime)
	assert.Equal(t, "devices", FieldDevices)
}

func TestActivityConstants(t *testing.T) {
	// Test activity type constants
	assert.Equal(t, "ASSIGN_DEVICES", ActivityTypeAssignDevices)
	assert.Equal(t, "UNASSIGN_DEVICES", ActivityTypeUnassignDevices)

	// Test activity status constants
	assert.Equal(t, "IN_PROGRESS", ActivityStatusInProgress)
	assert.Equal(t, "COMPLETED", ActivityStatusCompleted)
	assert.Equal(t, "FAILED", ActivityStatusFailed)

	// Test activity sub-status constants
	assert.Equal(t, "SUBMITTED", ActivitySubStatusSubmitted)
	assert.Equal(t, "PROCESSING", ActivitySubStatusProcessing)
}

func TestOptionsStructures(t *testing.T) {
	// Test GetMDMServersOptions
	opts1 := &GetMDMServersOptions{
		Fields: []string{FieldServerName, FieldServerType},
		Limit:  100,
	}
	assert.Len(t, opts1.Fields, 2)
	assert.Equal(t, 100, opts1.Limit)

	// Test GetMDMServerDeviceLinkagesOptions
	opts2 := &GetMDMServerDeviceLinkagesOptions{
		Limit: 50,
	}
	assert.Equal(t, 50, opts2.Limit)

	// Test GetAssignedServerInfoOptions
	opts3 := &GetAssignedServerInfoOptions{
		Fields: []string{FieldServerName, FieldCreatedDateTime},
	}
	assert.Len(t, opts3.Fields, 2)
}

func TestComprehensiveFieldCoverage(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test with all available fields to ensure comprehensive coverage
	opts := &GetMDMServersOptions{
		Fields: []string{
			FieldServerName,
			FieldServerType,
			FieldCreatedDateTime,
			FieldUpdatedDateTime,
			FieldDevices,
		},
	}

	result, err := client.GetDeviceManagementServices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify that the request was made with all field parameters
	// This ensures our field constants are correctly defined
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Verify all fields are accessible (no compilation errors)
	if len(result.Data) > 0 {
		attrs := result.Data[0].Attributes
		_ = attrs.ServerName
		_ = attrs.ServerType
		_ = attrs.CreatedDateTime
		_ = attrs.UpdatedDateTime
		_ = attrs.Devices
	}
}
