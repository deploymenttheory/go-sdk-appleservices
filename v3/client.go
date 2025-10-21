package axm

import (
	core "github.com/deploymenttheory/go-api-sdk-apple/v3/core"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devicemanagement"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devices"
)

// Client provides unified access to all Apple Business Manager API services
// This eliminates the need for wrapper clients and provides direct service access
type Client struct {
	*core.Client
	DeviceManagement *devicemanagement.Service
	Devices          *devices.Service
}

// NewClient creates a new unified client with all services embedded
func NewClient(config core.Config) (*Client, error) {
	coreClient, err := core.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// NewClientBuilder returns a new client builder for fluent configuration
func NewClientBuilder() *core.ClientBuilder {
	return core.NewClientBuilder()
}

// NewClientFromEnv creates a client using environment variables
func NewClientFromEnv() (*Client, error) {
	coreClient, err := core.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// NewClientFromFile creates a client using credentials from files
func NewClientFromFile(keyID, issuerID, privateKeyPath string) (*Client, error) {
	coreClient, err := core.NewClientFromFile(keyID, issuerID, privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}

// Crypto functions exported from core
var (
	LoadPrivateKeyFromFile = core.LoadPrivateKeyFromFile
	ParsePrivateKey        = core.ParsePrivateKey
	ValidatePrivateKey     = core.ValidatePrivateKey
)

// Config type exported from core
type Config = core.Config

// NewClientWithBuilder creates a client using the builder pattern
func NewClientWithBuilder() *core.ClientBuilder {
	return core.NewClientBuilder()
}

// BuildClient creates a unified client from a configured builder
func BuildClient(builder *core.ClientBuilder) (*Client, error) {
	coreClient, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:           coreClient,
		DeviceManagement: devicemanagement.NewService(coreClient),
		Devices:          devices.NewService(coreClient),
	}, nil
}