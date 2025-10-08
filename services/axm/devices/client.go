package devices

import (
	"github.com/deploymenttheory/go-api-sdk-apple/client/axm"
)

// Client provides access to Apple Business Manager Device Management Services API
type Client struct {
	client *axm.Client
}

// NewClient creates a new device management client
func NewClient(c *axm.Client) *Client {
	return &Client{
		client: c,
	}
}
