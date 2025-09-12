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

	// OAuth endpoints - Apple uses different OAuth server
	AuthServerURL = "https://account.apple.com"
	TokenEndpoint = "/auth/oauth2/token"
	
	// OAuth constants
	GrantType           = "client_credentials"
	ClientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	BusinessScope       = "business.api"
	SchoolScope         = "school.api"
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
	TeamID     string        // Team ID from Apple (typically same as ClientID)
	KeyID      string        // Key ID from Apple
	PrivateKey string        // Private key content (PEM format)
	Scope      string        // OAuth scope ("business.api" or "school.api")
	Timeout    time.Duration // HTTP timeout
	RetryCount int           // Number of retries
	RetryDelay time.Duration // Delay between retries
	UserAgent  string        // User agent string
	Debug      bool          // Enable debug logging
}

// Claims represents JWT claims for Apple authentication
type Claims struct {
	Issuer    string `json:"iss"` // Team ID
	Subject   string `json:"sub"` // Client ID
	IssuedAt  int64  `json:"iat"` // Issued at timestamp
	ExpiresAt int64  `json:"exp"` // Expiration timestamp
	Audience  string `json:"aud"` // Always "https://account.apple.com/auth/oauth2/v2/token"
	JTI       string `json:"jti"` // Unique identifier
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

// generateJWT creates a JWT token for authentication using ES256 (Apple Business Manager API requirement)
func (c *AXMClient) generateJWT() (string, error) {
	now := time.Now()
	claims := Claims{
		Issuer:    c.Config.OrgID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(20 * time.Minute).Unix(), // Apple requires max 20 minutes
		Audience:  "appstoreconnect-v1",
	}

	// Apple Business Manager API requires ES256 signing with ECDSA keys
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.Config.KeyID

	return token.SignedString(c.privateKey)
}

// authenticate performs OAuth authentication with Apple
func (c *AXMClient) authenticate() error {
	jwtToken, err := c.generateJWT()
	if err != nil {
		return fmt.Errorf("failed to generate JWT: %w", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	resp, err := c.HTTP.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&tokenResponse).
		Post(TokenEndpoint)

	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	c.accessToken = tokenResponse.AccessToken
	c.tokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	c.Logger.Info("Successfully authenticated with Apple School and Business Manager API",
		zap.String("token_type", tokenResponse.TokenType),
		zap.Int("expires_in", tokenResponse.ExpiresIn),
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

// GetOrgID returns the configured organization ID
func (c *AXMClient) GetOrgID() string {
	return c.Config.OrgID
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
