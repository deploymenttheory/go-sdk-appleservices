package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ClientOption is a function type for configuring the Transport.
type ClientOption func(*Transport) error

// WithBaseURL sets the base URL for API requests to a custom endpoint.
func WithBaseURL(urlStr string) ClientOption {
	return func(c *Transport) error {
		if urlStr == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		c.baseURL = urlStr
		c.logger.Info("Base URL configured", zap.String("base_url", urlStr))
		return nil
	}
}

// WithLogger can be used to configure a custom logger.
func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Transport) error {
		if logger == nil {
			return fmt.Errorf("logger cannot be nil")
		}
		c.logger = logger
		c.logger.Info("Custom logger configured")
		return nil
	}
}

// WithAuth sets the authentication provider for the client.
func WithAuth(auth AuthProvider) ClientOption {
	return func(c *Transport) error {
		if auth == nil {
			return fmt.Errorf("auth provider cannot be nil")
		}
		c.auth = auth
		c.logger.Info("Custom auth provider configured")
		return nil
	}
}

// WithTimeout sets the timeout for all HTTP requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Transport) error {
		if timeout < 0 {
			return fmt.Errorf("timeout cannot be negative")
		}
		c.httpClient.SetTimeout(timeout)
		c.logger.Info("HTTP timeout configured", zap.Duration("timeout", timeout))
		return nil
	}
}

// WithRetryCount sets the maximum number of retries for failed requests.
func WithRetryCount(retryCount int) ClientOption {
	return func(c *Transport) error {
		if retryCount < 0 {
			return fmt.Errorf("retry count cannot be negative")
		}
		c.httpClient.SetRetryCount(retryCount)
		c.logger.Info("Retry count configured", zap.Int("retry_count", retryCount))
		return nil
	}
}

// WithRetryWaitTime sets the default wait time between retry attempts.
func WithRetryWaitTime(retryWait time.Duration) ClientOption {
	return func(c *Transport) error {
		if retryWait < 0 {
			return fmt.Errorf("retry wait time cannot be negative")
		}
		c.httpClient.SetRetryWaitTime(retryWait)
		c.logger.Info("Retry wait time configured", zap.Duration("wait_time", retryWait))
		return nil
	}
}

// WithRetryMaxWaitTime sets the maximum wait time between retry attempts.
func WithRetryMaxWaitTime(maxWait time.Duration) ClientOption {
	return func(c *Transport) error {
		if maxWait < 0 {
			return fmt.Errorf("retry max wait time cannot be negative")
		}
		c.httpClient.SetRetryMaxWaitTime(maxWait)
		c.logger.Info("Retry max wait time configured", zap.Duration("max_wait_time", maxWait))
		return nil
	}
}

// WithUserAgent sets a custom user agent string for all requests.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Transport) error {
		if userAgent == "" {
			return fmt.Errorf("user agent cannot be empty")
		}
		c.httpClient.SetHeader("User-Agent", userAgent)
		c.logger.Info("User agent configured", zap.String("user_agent", userAgent))
		return nil
	}
}

// WithCustomAgent allows appending a custom identifier to the default user agent.
// Format: "go-api-sdk-apple/1.0.0; <customAgent>"
func WithCustomAgent(customAgent string) ClientOption {
	return func(c *Transport) error {
		enhancedUA := fmt.Sprintf("%s; %s", DefaultUserAgent, customAgent)
		c.httpClient.SetHeader("User-Agent", enhancedUA)
		c.logger.Info("Custom agent configured", zap.String("user_agent", enhancedUA))
		return nil
	}
}

// WithDebug enables debug mode for the HTTP client.
func WithDebug() ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetDebug(true)
		c.logger.Info("Debug mode enabled")
		return nil
	}
}

// WithErrorHandler sets a custom error handler.
func WithErrorHandler(handler *ErrorHandler) ClientOption {
	return func(c *Transport) error {
		if handler == nil {
			return fmt.Errorf("error handler cannot be nil")
		}
		c.errorHandler = handler
		c.logger.Info("Custom error handler configured")
		return nil
	}
}

// WithGlobalHeader sets a global header that will be included in all requests.
func WithGlobalHeader(key, value string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetHeader(key, value)
		c.logger.Info("Global header configured", zap.String("key", key), zap.String("value", value))
		return nil
	}
}

// WithGlobalHeaders sets multiple global headers at once.
func WithGlobalHeaders(headers map[string]string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetHeaders(headers)
		c.logger.Info("Multiple global headers configured", zap.Int("count", len(headers)))
		return nil
	}
}

// WithProxy sets an HTTP proxy for all requests.
func WithProxy(proxyURL string) ClientOption {
	return func(c *Transport) error {
		if proxyURL == "" {
			return fmt.Errorf("proxy URL cannot be empty")
		}
		c.httpClient.SetProxy(proxyURL)
		c.logger.Info("Proxy configured", zap.String("proxy", proxyURL))
		return nil
	}
}

// WithTLSClientConfig sets custom TLS configuration.
func WithTLSClientConfig(tlsConfig *tls.Config) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetTLSClientConfig(tlsConfig)
		c.logger.Info("TLS client config configured",
			zap.Uint16("min_version", tlsConfig.MinVersion),
			zap.Bool("insecure_skip_verify", tlsConfig.InsecureSkipVerify))
		return nil
	}
}

// WithClientCertificate sets a client certificate for mutual TLS authentication.
func WithClientCertificate(certFile, keyFile string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetCertificateFromFile(certFile, keyFile)
		c.logger.Info("Client certificate configured",
			zap.String("cert_file", certFile),
			zap.String("key_file", keyFile))
		return nil
	}
}

// WithClientCertificateFromString sets a client certificate from PEM-encoded strings.
func WithClientCertificateFromString(certPEM, keyPEM string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetCertificateFromString(certPEM, keyPEM)
		c.logger.Info("Client certificate configured from string")
		return nil
	}
}

// WithRootCertificates adds custom root CA certificates for server validation.
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetClientRootCertificates(pemFilePaths...)
		c.logger.Info("Root certificates configured", zap.Int("count", len(pemFilePaths)))
		return nil
	}
}

// WithRootCertificateFromString adds a custom root CA certificate from PEM string.
func WithRootCertificateFromString(pemContent string) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetClientRootCertificateFromString(pemContent)
		c.logger.Info("Root certificate configured from string")
		return nil
	}
}

// WithTransport sets a custom HTTP transport (http.RoundTripper).
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *Transport) error {
		c.httpClient.SetTransport(transport)
		c.logger.Info("Custom transport configured")
		return nil
	}
}

// WithInsecureSkipVerify disables TLS certificate verification (USE WITH CAUTION).
func WithInsecureSkipVerify() ClientOption {
	return func(c *Transport) error {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		c.httpClient.SetTLSClientConfig(tlsConfig)
		c.logger.Warn("TLS certificate verification DISABLED - use only for testing")
		return nil
	}
}

// WithMinTLSVersion sets the minimum TLS version for connections.
func WithMinTLSVersion(minVersion uint16) ClientOption {
	return func(c *Transport) error {
		tlsConfig := &tls.Config{
			MinVersion: minVersion,
		}
		c.httpClient.SetTLSClientConfig(tlsConfig)

		versionName := "unknown"
		switch minVersion {
		case tls.VersionTLS10:
			versionName = "TLS 1.0"
		case tls.VersionTLS11:
			versionName = "TLS 1.1"
		case tls.VersionTLS12:
			versionName = "TLS 1.2"
		case tls.VersionTLS13:
			versionName = "TLS 1.3"
		}

		c.logger.Info("Minimum TLS version configured",
			zap.String("version", versionName),
			zap.Uint16("version_code", minVersion))
		return nil
	}
}

// WithAudience sets a custom JWT audience (default: "appstoreconnect-v1").
func WithAudience(audience string) ClientOption {
	return func(c *Transport) error {
		if jwtAuth, ok := c.auth.(*JWTAuth); ok {
			jwtAuth.audience = audience
			c.logger.Info("JWT audience configured", zap.String("audience", audience))
		}
		return nil
	}
}
