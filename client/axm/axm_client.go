package client

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Apple School and Business Manager API base URLs
	AppleSchoolManagerBaseURL   = "https://axm-adm-enroll.apple.com"
	AppleBusinessManagerBaseURL = "https://axm-adm-enroll.apple.com"

	// OAuth endpoints
	TokenEndpoint = "https://account.apple.com/auth/oauth2/token"
	
	// OAuth constants
	GrantType           = "client_credentials"
	ClientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	BusinessScope       = "business.api"
	SchoolScope         = "school.api"
	
	// JWT Constants
	TokenAudience = "https://account.apple.com/auth/oauth2/v2/token"
)

// AXMClient represents a client for Apple School and Business Manager API
type AXMClient struct {
	HTTP           *resty.Client
	Logger         *zap.Logger
	Config         AXMConfig
	privateKey     *ecdsa.PrivateKey
	accessToken    string
	tokenExpiresAt time.Time
}

// AXMConfig holds configuration for the AXM client
type AXMConfig struct {
	BaseURL    string        // Apple School/Business Manager API base URL
	ClientID   string        // Client ID from Apple (e.g., "BUSINESSAPI.9703f56c-10ce-4876-8f59-e78e5e23a152")
	KeyID      string        // Key ID from Apple
	PrivateKey string        // Private key content (PEM format)
	Scope      string        // OAuth scope ("business.api" or "school.api")
	Timeout    time.Duration // HTTP timeout
	RetryCount int           // Number of retries
	RetryDelay time.Duration // Delay between retries
	UserAgent  string        // User agent string
	Debug      bool          // Enable debug logging
}

// Claims represents JWT claims for Apple AXM authentication
type Claims struct {
	Issuer    string `json:"iss"` // Client ID (same as Subject)
	Subject   string `json:"sub"` // Client ID (same as Issuer)
	IssuedAt  int64  `json:"iat"` // Issued at timestamp
	ExpiresAt int64  `json:"exp"` // Expiration timestamp (max 20 minutes from iat)
	Audience  string `json:"aud"` // Always "https://account.apple.com/auth/oauth2/v2/token"
	JTI       string `json:"jti"` // Unique identifier for this JWT
	jwt.RegisteredClaims
}

// NewAXMClient creates a new Apple School and Business Manager API client
func NewAXMClient(config AXMConfig) (*AXMClient, error) {
	var logger *zap.Logger
	var err error

	if config.Debug {
		// Development config with colors and console encoder
		developmentConfig := zap.NewDevelopmentConfig()
		developmentConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		developmentConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		developmentConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		logger, err = developmentConfig.Build()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		logger = zap.NewNop()
	}

	// Set defaults
	if config.BaseURL == "" {
		config.BaseURL = AppleSchoolManagerBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}
	if config.UserAgent == "" {
		config.UserAgent = "go-api-sdk-apple/1.0.0"
	}

	// Parse private key
	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create HTTP client
	httpClient := resty.New().
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryDelay).
		SetHeader("User-Agent", config.UserAgent).
		SetBaseURL(config.BaseURL)

	if config.Debug {
		httpClient.SetDebug(true)
	}

	client := &AXMClient{
		HTTP:       httpClient,
		Logger:     logger,
		Config:     config,
		privateKey: privateKey,
	}

	return client, nil
}

// parsePrivateKey parses a PEM-encoded ECDSA private key (Apple Business Manager API requirement)
func parsePrivateKey(keyData string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	switch block.Type {
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		// PKCS#8 format - extract ECDSA key
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		
		ecdsaKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not an ECDSA key (found %T). Apple Business Manager API requires ECDSA keys", key)
		}
		return ecdsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q. Apple Business Manager API requires ECDSA keys in 'EC PRIVATE KEY' or 'PRIVATE KEY' (PKCS#8) format", block.Type)
	}
}

// generateJWT creates a JWT token for authentication using ES256 (Apple AXM API requirement)
func (c *AXMClient) generateJWT() (string, error) {
	now := time.Now()
	
	// Generate unique JWT ID
	jti := fmt.Sprintf("%s-%d", c.Config.ClientID, now.Unix())
	
	claims := Claims{
		Issuer:    c.Config.ClientID,                    // Client ID
		Subject:   c.Config.ClientID,                    // Same as Issuer for AXM API
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(20 * time.Minute).Unix(),    // Apple requires max 20 minutes
		Audience:  TokenAudience,                       // Apple's OAuth token endpoint
		JTI:       jti,                                 // Unique identifier
	}

	// Apple AXM API requires ES256 signing with ECDSA keys
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.Config.KeyID

	return token.SignedString(c.privateKey)
}

// authenticate performs OAuth 2.0 client credentials authentication with Apple
func (c *AXMClient) authenticate() error {
	jwtToken, err := c.generateJWT()
	if err != nil {
		return fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Use default scope if not specified
	scope := c.Config.Scope
	if scope == "" {
		scope = BusinessScope
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope,omitempty"`
	}

	// Prepare form data for OAuth 2.0 client credentials flow
	formData := map[string]string{
		"grant_type":            GrantType,
		"client_id":             c.Config.ClientID,
		"client_assertion_type": ClientAssertionType,
		"client_assertion":      jwtToken,
		"scope":                 scope,
	}

	resp, err := c.HTTP.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		SetResult(&tokenResponse).
		Post(TokenEndpoint)

	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.IsError() {
		c.Logger.Error("Authentication failed", 
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response", resp.String()))
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	c.accessToken = tokenResponse.AccessToken
	c.tokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	c.Logger.Info("Successfully authenticated with Apple School and Business Manager API",
		zap.String("token_type", tokenResponse.TokenType),
		zap.Int("expires_in", tokenResponse.ExpiresIn),
		zap.String("scope", tokenResponse.Scope),
	)

	return nil
}

// ensureAuthenticated ensures we have a valid access token
func (c *AXMClient) ensureAuthenticated() error {
	// Check if we have a token and it's not expired (with 5 minute buffer)
	if c.accessToken != "" && time.Now().Before(c.tokenExpiresAt.Add(-5*time.Minute)) {
		return nil
	}

	return c.authenticate()
}

// Close cleans up resources
func (c *AXMClient) Close() {
	if c.Logger != nil {
		c.Logger.Sync()
	}
}

// GetClientID returns the configured client ID
func (c *AXMClient) GetClientID() string {
	return c.Config.ClientID
}

// IsAuthenticated returns true if we have a valid access token
func (c *AXMClient) IsAuthenticated() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiresAt)
}

// ForceReauthenticate forces a new authentication cycle
func (c *AXMClient) ForceReauthenticate() error {
	c.accessToken = ""
	c.tokenExpiresAt = time.Time{}
	return c.authenticate()
}
