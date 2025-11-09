package apps

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/interfaces"
)

type (
	// AppsServiceInterface defines the interface for apps operations
	AppsServiceInterface interface {
		// GetLatestApps retrieves the latest Microsoft Mac application versions
		//
		// API docs: https://appledevicepolicy.tools/microsoft-apps
		GetLatestApps(ctx context.Context) (*AppsResponse, error)

		// GetAppByBundleID finds an application by its bundle ID
		GetAppByBundleID(ctx context.Context, bundleID string) (*App, error)

		// GetAppByName finds an application by its name
		GetAppByName(ctx context.Context, name string) (*App, error)
	}

	// AppsService handles communication with the apps
	// related methods of the Microsoft Mac Apps API
	AppsService struct {
		client interfaces.HTTPClient
	}
)

var _ AppsServiceInterface = (*AppsService)(nil)

// NewService creates a new apps service
func NewService(client interfaces.HTTPClient) *AppsService {
	return &AppsService{
		client: client,
	}
}

// GetLatestApps retrieves the latest Microsoft Mac application versions
// URL: GET https://appledevicepolicy.tools/api/latest
func (s *AppsService) GetLatestApps(ctx context.Context) (*AppsResponse, error) {
	var result AppsResponse

	headers := map[string]string{
		"Accept": "application/json",
	}

	err := s.client.Get(ctx, EndpointLatest, nil, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAppByBundleID finds an application by its bundle ID
func (s *AppsService) GetAppByBundleID(ctx context.Context, bundleID string) (*App, error) {
	apps, err := s.GetLatestApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps.Apps {
		if app.BundleID == bundleID {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("app with bundle ID %s not found", bundleID)
}

// GetAppByName finds an application by its name
func (s *AppsService) GetAppByName(ctx context.Context, name string) (*App, error) {
	apps, err := s.GetLatestApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps.Apps {
		if app.Name == name {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("app with name %s not found", name)
}
