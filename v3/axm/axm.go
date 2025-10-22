package axm

import (
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/client"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devicemanagement"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devices"
)

// Client provides unified access to all Apple Business Manager API services
// following the GitLab go-gitlab pattern
type Client struct {
	*client.Client

	// Services using interfaces for better testability and no import cycles
	DeviceManagement devicemanagement.DeviceManagementServiceInterface
	Devices          devices.DevicesServiceInterface
}

// NewClient creates a new Apple Business Manager client following GitLab pattern
// Usage: client, err := axm.NewClient("keyID", "issuerID", privateKey)
func NewClient(keyID, issuerID string, privateKey any) (*Client, error) {
	config := client.Config{
		BaseURL: "https://api-business.apple.com/v1",
		Auth: client.NewJWTAuth(client.JWTAuthConfig{
			KeyID:      keyID,
			IssuerID:   issuerID,
			PrivateKey: privateKey,
			Audience:   "appstoreconnect-v1",
		}),
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryWait:  1 * time.Second,
		UserAgent:  "go-api-sdk-apple/3.0.0",
		Debug:      true,
	}

	coreClient, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// NewClientFromFile creates a client using private key from file
func NewClientFromFile(keyID, issuerID, privateKeyPath string) (*Client, error) {
	privateKey, err := client.LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return NewClient(keyID, issuerID, privateKey)
}

// NewClientFromEnv creates a client using environment variables
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
func NewClientFromEnv() (*Client, error) {
	coreClient, err := client.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// Helper functions following GitLab pattern for pointer values

// Ptr returns a pointer to the provided value
func Ptr[T any](v T) *T {
	return &v
}

// String returns a pointer to the provided string value
func String(v string) *string {
	return &v
}

// Int returns a pointer to the provided int value
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the provided bool value
func Bool(v bool) *bool {
	return &v
}

// Time returns a pointer to the provided time value
func Time(v time.Time) *time.Time {
	return &v
}

// Type aliases for common types
type (
	Config        = client.Config
	JWTAuthConfig = client.JWTAuthConfig
)

// Re-export common functions
var (
	LoadPrivateKeyFromFile = client.LoadPrivateKeyFromFile
	ParsePrivateKey        = client.ParsePrivateKey
	ValidatePrivateKey     = client.ValidatePrivateKey
)
