package msapps

import (
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/client"
	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/services/apps"
)

// Client provides unified access to all Microsoft Mac Apps API services
type Client struct {
	*client.Client

	// Apps service
	Apps apps.AppsServiceInterface
}

// NewClient creates a new Microsoft Mac Apps API client
func NewClient(config *client.Config) (*Client, error) {
	if config == nil {
		config = &client.Config{}
	}

	if config.BaseURL == "" {
		config.BaseURL = client.DefaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryWait == 0 {
		config.RetryWait = 1 * time.Second
	}
	if config.UserAgent == "" {
		config.UserAgent = "go-api-sdk-apple-msapps/1.0.0"
	}

	coreClient, err := client.NewTransport(*config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: coreClient,
		Apps:   apps.NewService(coreClient),
	}, nil
}
