package client

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ClientOptionFunc can be used to customize a new Apple Business Manager API client.
type ClientOptionFunc func(*Client) error

// WithBaseURL sets the base URL for API requests to a custom endpoint.
func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		if urlStr == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		c.baseURL = urlStr
		return nil
	}
}

// WithLogger can be used to configure a custom logger.
func WithLogger(logger *zap.Logger) ClientOptionFunc {
	return func(c *Client) error {
		if logger == nil {
			return fmt.Errorf("logger cannot be nil")
		}
		c.logger = logger
		return nil
	}
}

// WithAuth sets the authentication provider for the client.
func WithAuth(auth AuthProvider) ClientOptionFunc {
	return func(c *Client) error {
		if auth == nil {
			return fmt.Errorf("auth provider cannot be nil")
		}
		c.auth = auth
		return nil
	}
}

// WithTimeout sets the timeout for all HTTP requests.
func WithTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		if timeout < 0 {
			return fmt.Errorf("timeout cannot be negative")
		}
		c.httpClient.SetTimeout(timeout)
		return nil
	}
}

// WithRetryCount sets the maximum number of retries for failed requests.
func WithRetryCount(retryCount int) ClientOptionFunc {
	return func(c *Client) error {
		if retryCount < 0 {
			return fmt.Errorf("retry count cannot be negative")
		}
		c.httpClient.SetRetryCount(retryCount)
		return nil
	}
}

// WithRetryWaitTime sets the wait time between retries.
func WithRetryWaitTime(retryWait time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		if retryWait < 0 {
			return fmt.Errorf("retry wait time cannot be negative")
		}
		c.httpClient.SetRetryWaitTime(retryWait)
		return nil
	}
}

// WithRetryMaxWaitTime sets the maximum wait time between retries.
func WithRetryMaxWaitTime(maxWait time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		if maxWait < 0 {
			return fmt.Errorf("retry max wait time cannot be negative")
		}
		c.httpClient.SetRetryMaxWaitTime(maxWait)
		return nil
	}
}

// WithUserAgent sets a custom user agent string for all requests.
func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		if userAgent == "" {
			return fmt.Errorf("user agent cannot be empty")
		}
		c.httpClient.SetHeader("User-Agent", userAgent)
		return nil
	}
}

// WithDebug enables debug mode for the HTTP client.
func WithDebug(debug bool) ClientOptionFunc {
	return func(c *Client) error {
		c.httpClient.SetDebug(debug)
		return nil
	}
}

// WithErrorHandler sets a custom error handler.
func WithErrorHandler(handler *ErrorHandler) ClientOptionFunc {
	return func(c *Client) error {
		if handler == nil {
			return fmt.Errorf("error handler cannot be nil")
		}
		c.errorHandler = handler
		return nil
	}
}

// WithProxy sets a proxy for the HTTP client.
func WithProxy(proxyURL string) ClientOptionFunc {
	return func(c *Client) error {
		if proxyURL == "" {
			return fmt.Errorf("proxy URL cannot be empty")
		}
		c.httpClient.SetProxy(proxyURL)
		return nil
	}
}

// ApplyOptions applies the given options to the client.
func (c *Client) ApplyOptions(opts ...ClientOptionFunc) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}
