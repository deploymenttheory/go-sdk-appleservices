package apps

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/client"
	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/services/apps/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) *AppsService {
	config := client.Config{
		BaseURL:    "https://appledevicepolicy.tools/api",
		RetryCount: 0, // Disable retries for tests
		Debug:      false,
	}

	coreClient, err := client.NewTransport(config)
	require.NoError(t, err)

	// Activate httpmock for the client's HTTP client
	httpClient := coreClient.GetHTTPClient().(*resty.Client)
	httpmock.ActivateNonDefault(httpClient.Client())

	// Setup cleanup
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(coreClient)
}

func TestGetLatestApps_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, err := svc.GetLatestApps(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Apps)
	assert.NotEmpty(t, result.Generated)

	// Verify we have multiple apps
	assert.Greater(t, len(result.Apps), 0)

	// Find and verify Excel
	var excel *App
	for _, app := range result.Apps {
		if app.BundleID == BundleIDExcel {
			excel = &app
			break
		}
	}

	require.NotNil(t, excel, "Excel app should be in the response")
	assert.Equal(t, "Microsoft Excel", excel.Name)
	assert.NotEmpty(t, excel.Version)
	assert.Equal(t, TypeApplication, excel.Type)
	assert.NotEmpty(t, excel.DirectURL)
	assert.NotEmpty(t, excel.DownloadURL)
	assert.NotEmpty(t, excel.SHA256)
	assert.Greater(t, excel.SizeBytes, int64(0))
	assert.Greater(t, excel.SizeMB, 0.0)

	// Verify exactly one HTTP call was made
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetLatestApps_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	// Reset httpmock to clear any previous registrations
	httpmock.Reset()

	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, err := svc.GetLatestApps(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API returned error status")

	assert.Equal(t, 4, httpmock.GetTotalCallCount()) // 1 original + 3 retries
}

func TestGetLatestApps_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := svc.GetLatestApps(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetLatestApps_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep to ensure timeout
	time.Sleep(1 * time.Millisecond)

	result, err := svc.GetLatestApps(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestGetAppByBundleID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test finding Excel
	app, err := svc.GetAppByBundleID(ctx, BundleIDExcel)

	require.NoError(t, err)
	require.NotNil(t, app)
	assert.Equal(t, BundleIDExcel, app.BundleID)
	assert.Equal(t, "Microsoft Excel", app.Name)
	assert.NotEmpty(t, app.Version)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppByBundleID_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	app, err := svc.GetAppByBundleID(ctx, "com.nonexistent.app")

	require.Error(t, err)
	assert.Nil(t, app)
	assert.Contains(t, err.Error(), "not found")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppByName_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test finding Teams
	app, err := svc.GetAppByName(ctx, AppNameTeams)

	require.NoError(t, err)
	require.NotNil(t, app)
	assert.Equal(t, BundleIDTeams, app.BundleID)
	assert.Equal(t, AppNameTeams, app.Name)
	assert.NotEmpty(t, app.Version)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppByName_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	app, err := svc.GetAppByName(ctx, "Non-existent App")

	require.Error(t, err)
	assert.Nil(t, app)
	assert.Contains(t, err.Error(), "not found")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestParseDetectedTime(t *testing.T) {
	app := &App{
		Detected: "2025-11-04T18:37:44.058059",
	}

	parsed, err := app.ParseDetectedTime()

	require.NoError(t, err)
	assert.Equal(t, 2025, parsed.Year())
	assert.Equal(t, time.November, parsed.Month())
	assert.Equal(t, 4, parsed.Day())
	assert.Equal(t, 18, parsed.Hour())
	assert.Equal(t, 37, parsed.Minute())
}

func TestParseLastModifiedTime(t *testing.T) {
	app := &App{
		LastModified: "Tue, 04 Nov 2025 16:40:19 GMT",
	}

	parsed, err := app.ParseLastModifiedTime()

	require.NoError(t, err)
	assert.Equal(t, 2025, parsed.Year())
	assert.Equal(t, time.November, parsed.Month())
	assert.Equal(t, 4, parsed.Day())
	assert.Equal(t, 16, parsed.Hour())
	assert.Equal(t, 40, parsed.Minute())
}

func TestParseGeneratedTime(t *testing.T) {
	response := &AppsResponse{
		Generated: "2025-11-08 19:59:18",
	}

	parsed, err := response.ParseGeneratedTime()

	require.NoError(t, err)
	assert.Equal(t, 2025, parsed.Year())
	assert.Equal(t, time.November, parsed.Month())
	assert.Equal(t, 8, parsed.Day())
	assert.Equal(t, 19, parsed.Hour())
	assert.Equal(t, 59, parsed.Minute())
	assert.Equal(t, 18, parsed.Second())
}

func TestBundleIDConstants(t *testing.T) {
	// Test that bundle ID constants are properly defined
	assert.Equal(t, "com.microsoft.Excel", BundleIDExcel)
	assert.Equal(t, "com.microsoft.teams2", BundleIDTeams)
	assert.Equal(t, "com.microsoft.VSCode", BundleIDVSCode)
	assert.Equal(t, "com.microsoft.Outlook", BundleIDOutlook)
	assert.Equal(t, "com.microsoft.Word", BundleIDWord)
	assert.Equal(t, "com.microsoft.wdav", BundleIDDefender)
	assert.Equal(t, "com.microsoft.autoupdate2", BundleIDAutoUpdate)
}

func TestAppNameConstants(t *testing.T) {
	// Test that app name constants are properly defined
	assert.Equal(t, "Microsoft Excel", AppNameExcel)
	assert.Equal(t, "Microsoft Teams", AppNameTeams)
	assert.Equal(t, "Visual Studio Code", AppNameVSCode)
	assert.Equal(t, "Microsoft Outlook", AppNameOutlook)
	assert.Equal(t, "Microsoft Word", AppNameWord)
	assert.Equal(t, "Defender for Mac", AppNameDefender)
	assert.Equal(t, "Microsoft AutoUpdate", AppNameAutoUpdate)
}

func TestMultipleSequentialRequests(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.MSAppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Make multiple requests
	for i := 0; i < 3; i++ {
		result, err := svc.GetLatestApps(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Apps)
	}

	assert.Equal(t, 3, httpmock.GetTotalCallCount())
}
