package axm

import (
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devicemanagement"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
)

// Client provides unified access to all Apple Business Manager API services
type Client struct {
	*client.Client

	// Imports for services
	DeviceManagement devicemanagement.DeviceManagementServiceInterface
	Devices          devices.DevicesServiceInterface
}

// NewClient creates a new Apple Business Manager client
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

	coreClient, err := client.NewTransport(config)
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
	coreClient, err := client.NewTransportFromEnv()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}
