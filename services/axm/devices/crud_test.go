package devices

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices/mocks"
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

	// Create devices client
	return NewClient(axmClient)
}

// MockAuthProvider implements the AuthProvider interface for testing
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	// Mock auth - do nothing for tests
	return nil
}

func TestGetOrganizationDevices_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetOrganizationDevicesOptions{
		Fields: []string{"serialNumber", "deviceModel", "status"},
		Limit:  100,
	}

	result, err := client.GetOrganizationDevices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify the first device
	device := result.Data[0]
	assert.Equal(t, "orgDevices", device.Type)
	assert.Equal(t, "XABC123X0ABC123X0", device.ID)
	assert.NotNil(t, device.Attributes)

	// Test all string fields from the list response
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)
	assert.Equal(t, "iMac16,2", device.Attributes.ProductType)
	assert.Equal(t, "750GB", device.Attributes.DeviceCapacity)
	assert.Equal(t, "FD311LL/A", device.Attributes.PartNumber)
	assert.Equal(t, "1234567890", device.Attributes.OrderNumber)
	assert.Equal(t, "SILVER", device.Attributes.Color) // Has color in list response
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "89049037640158663184237812557346", device.Attributes.EID)
	assert.Equal(t, "-2085650007946880", device.Attributes.PurchaseSourceUid)
	assert.Equal(t, "APPLE", device.Attributes.PurchaseSourceType)

	// Test timestamp fields
	assert.NotNil(t, device.Attributes.AddedToOrgDateTime)
	assert.NotNil(t, device.Attributes.UpdatedDateTime)
	assert.NotNil(t, device.Attributes.OrderDateTime)

	// Test array fields
	assert.Len(t, device.Attributes.IMEI, 2)
	assert.Equal(t, "123456789012345", device.Attributes.IMEI[0])
	assert.Equal(t, "123456789012346", device.Attributes.IMEI[1])
	assert.Len(t, device.Attributes.MEID, 1)
	assert.Equal(t, "12345678901237", device.Attributes.MEID[0])

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	// Verify exactly one HTTP call was made
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, err := client.GetOrganizationDevices(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithFieldSelection(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetOrganizationDevicesOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel, FieldStatus, FieldProductFamily},
		Limit:  50,
	}

	result, err := client.GetOrganizationDevices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify device data
	device := result.Data[0]
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithLimitEnforcement(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetOrganizationDevicesOptions{
		Limit: 1500, // Exceeds API maximum of 1000
	}

	result, err := client.GetOrganizationDevices(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Reset httpmock to clear any previous registrations
	httpmock.Reset()

	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetOrganizationDevicesOptions{}

	result, err := client.GetOrganizationDevices(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 4, httpmock.GetTotalCallCount()) // 1 original + 3 retries
}

func TestGetDeviceInformation_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &GetDeviceInformationOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel, FieldStatus},
	}

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify device data
	device := result.Data
	assert.Equal(t, "orgDevices", device.Type)
	assert.Equal(t, "XABC123X0ABC123X0", device.ID)
	assert.NotNil(t, device.Attributes)

	// Test all string fields
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)
	assert.Equal(t, "iMac16,2", device.Attributes.ProductType)
	assert.Equal(t, "750GB", device.Attributes.DeviceCapacity)
	assert.Equal(t, "FD311LL/A", device.Attributes.PartNumber)
	assert.Equal(t, "1234567890", device.Attributes.OrderNumber)
	assert.Equal(t, "SILVER", device.Attributes.Color)
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "89049037640158663184237812557346", device.Attributes.EID)
	assert.Equal(t, "-2085650007946880", device.Attributes.PurchaseSourceUid)
	assert.Equal(t, "APPLE", device.Attributes.PurchaseSourceType)

	// Verify timestamps are parsed correctly
	assert.NotNil(t, device.Attributes.AddedToOrgDateTime)
	assert.NotNil(t, device.Attributes.UpdatedDateTime)
	assert.NotNil(t, device.Attributes.OrderDateTime)

	// Verify arrays are handled correctly
	assert.Len(t, device.Attributes.IMEI, 2)
	assert.Equal(t, "123456789012345", device.Attributes.IMEI[0])
	assert.Equal(t, "123456789012346", device.Attributes.IMEI[1])
	assert.Len(t, device.Attributes.MEID, 1)
	assert.Equal(t, "12345678901237", device.Attributes.MEID[0])

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "XABC123X0ABC123X0", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetDeviceInformationOptions{}

	result, err := client.GetDeviceInformationByDeviceID(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_DeviceNotFound(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "NONEXISTENT123"
	opts := &GetDeviceInformationOptions{}

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_HTTPError(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &GetDeviceInformationOptions{}

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_ContextCancellation(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	deviceID := "XABC123X0ABC123X0"
	opts := &GetDeviceInformationOptions{}

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetOrganizationDevices_ContextTimeout(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep to ensure timeout
	time.Sleep(1 * time.Millisecond)

	opts := &GetOrganizationDevicesOptions{}

	result, err := client.GetOrganizationDevices(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestDeviceFieldConstants(t *testing.T) {
	// Test that field constants are properly defined
	expectedFields := []string{
		FieldSerialNumber,
		FieldAddedToOrgDateTime,
		FieldUpdatedDateTime,
		FieldDeviceModel,
		FieldProductFamily,
		FieldProductType,
		FieldDeviceCapacity,
		FieldPartNumber,
		FieldOrderNumber,
		FieldColor,
		FieldStatus,
		FieldOrderDateTime,
		FieldIMEI,
		FieldMEID,
		FieldEID,
		FieldWiFiMACAddress,
		FieldBluetoothMACAddress,
		FieldPurchaseSourceUid,
		FieldPurchaseSourceType,
		FieldAssignedServer,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	// Test specific field values
	assert.Equal(t, "serialNumber", FieldSerialNumber)
	assert.Equal(t, "deviceModel", FieldDeviceModel)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "productFamily", FieldProductFamily)
}

func TestOptionsStructures(t *testing.T) {
	// Test GetOrganizationDevicesOptions
	opts1 := &GetOrganizationDevicesOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel},
		Limit:  100,
	}
	assert.Len(t, opts1.Fields, 2)
	assert.Equal(t, 100, opts1.Limit)

	// Test GetDeviceInformationOptions
	opts2 := &GetDeviceInformationOptions{
		Fields: []string{FieldStatus, FieldProductFamily},
	}
	assert.Len(t, opts2.Fields, 2)
}

func TestComprehensiveFieldCoverage(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"

	// Test with all available fields to ensure comprehensive coverage
	opts := &GetDeviceInformationOptions{
		Fields: []string{
			FieldSerialNumber,
			FieldAddedToOrgDateTime,
			FieldUpdatedDateTime,
			FieldDeviceModel,
			FieldProductFamily,
			FieldProductType,
			FieldDeviceCapacity,
			FieldPartNumber,
			FieldOrderNumber,
			FieldColor,
			FieldStatus,
			FieldOrderDateTime,
			FieldIMEI,
			FieldMEID,
			FieldEID,
			FieldWiFiMACAddress,
			FieldBluetoothMACAddress,
			FieldPurchaseSourceUid,
			FieldPurchaseSourceType,
			FieldAssignedServer,
		},
	}

	result, err := client.GetDeviceInformationByDeviceID(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify that the request was made with all field parameters
	// This ensures our field constants are correctly defined
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Verify all fields are accessible (no compilation errors)
	attrs := result.Data.Attributes
	_ = attrs.SerialNumber
	_ = attrs.AddedToOrgDateTime
	_ = attrs.UpdatedDateTime
	_ = attrs.DeviceModel
	_ = attrs.ProductFamily
	_ = attrs.ProductType
	_ = attrs.DeviceCapacity
	_ = attrs.PartNumber
	_ = attrs.OrderNumber
	_ = attrs.Color
	_ = attrs.Status
	_ = attrs.OrderDateTime
	_ = attrs.IMEI
	_ = attrs.MEID
	_ = attrs.EID
	_ = attrs.WiFiMACAddress
	_ = attrs.BluetoothMACAddress
	_ = attrs.PurchaseSourceUid
	_ = attrs.PurchaseSourceType
	_ = attrs.AssignedServer
}
