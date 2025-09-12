package axm

import (
	client "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"go.uber.org/zap"
)

type Client struct {
	axmClient *client.AXMClient
	logger    *zap.Logger
}

// Close cleans up client resources
func (c *Client) Close() {
	if c.axmClient != nil {
		c.axmClient.Close()
	}
}

// NewQueryBuilder creates a new query parameter builder
func (c *Client) NewQueryBuilder() *client.QueryBuilder {
	return c.axmClient.NewQueryBuilder()
}

func NewClient(axmClient *client.AXMClient) *Client {
	return &Client{
		axmClient: axmClient,
		logger:    axmClient.Logger,
	}
}

func NewClientWithConfig(config client.AXMConfig) (*Client, error) {
	axmClient, err := client.NewAXMClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		axmClient: axmClient,
		logger:    axmClient.Logger,
	}, nil
}
