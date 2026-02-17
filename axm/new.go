package axm

import (
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

// NewClient creates a new Apple Business Manager client.
// Parameters:
//   - keyID: Your Apple Developer Key ID
//   - issuerID: Your Apple Developer Issuer ID (Team ID)
//   - privateKey: Your Apple Developer private key (*rsa.PrivateKey or *ecdsa.PrivateKey)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClient(keyID, issuerID string, privateKey any, options ...client.ClientOption) (*Client, error) {
	coreClient, err := client.NewTransport(keyID, issuerID, privateKey, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// NewClientFromFile creates a client using private key from file.
// Parameters:
//   - keyID: Your Apple Developer Key ID
//   - issuerID: Your Apple Developer Issuer ID (Team ID)
//   - privateKeyPath: Path to your Apple Developer private key file (.p8)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromFile(keyID, issuerID, privateKeyPath string, options ...client.ClientOption) (*Client, error) {
	privateKey, err := client.LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return NewClient(keyID, issuerID, privateKey, options...)
}

// NewClientFromEnv creates a client using environment variables.
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
// Parameters:
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromEnv(options ...client.ClientOption) (*Client, error) {
	coreClient, err := client.NewTransportFromEnv(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}
