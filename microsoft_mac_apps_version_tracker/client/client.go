package client

import (
	"time"

	"resty.dev/v3"
)

// Client represents the core HTTP client for Microsoft Mac Apps API
type Client struct {
	httpClient *resty.Client
	baseURL    string
}

// Config holds configuration for the client
type Config struct {
	BaseURL    string
	Timeout    time.Duration
	RetryCount int
	RetryWait  time.Duration
	UserAgent  string
	Debug      bool
}

// NewTransport creates a new HTTP transport for Microsoft Mac Apps API
func NewTransport(config Config) (*Client, error) {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
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

	httpClient := resty.New()

	httpClient.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWait).
		SetRetryMaxWaitTime(config.RetryWait*10).
		SetHeader("User-Agent", config.UserAgent).
		SetHeader("Accept", "application/json")

	if config.Debug {
		httpClient.SetDebug(true)
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
	}, nil
}

// GetHTTPClient returns the underlying HTTP client for testing purposes
func (c *Client) GetHTTPClient() any {
	return c.httpClient
}

// Close closes the HTTP client and cleans up resources
func (c *Client) Close() error {
	if c.httpClient != nil {
		c.httpClient.Close()
	}
	return nil
}
