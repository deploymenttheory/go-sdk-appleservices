package itunes_search

import (
	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"go.uber.org/zap"
)

type Client struct {
	baseClient *client.Client
	logger     *zap.Logger
}

func NewClient(baseClient *client.Client) *Client {
	return &Client{
		baseClient: baseClient,
		logger:     baseClient.Logger,
	}
}

func NewDefaultClient() *Client {
	return NewClient(client.NewDefaultClient())
}
