package axm2

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"resty.dev/v3"
)

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

// TokenProvider handles JWT-based authentication for Apple AXM APIs
type TokenProvider struct {
	clientID   string
	keyID      string
	privateKey *ecdsa.PrivateKey
	scope      string
	logger     *zap.Logger

	// Token management
	mutex          sync.RWMutex
	accessToken    string
	tokenExpiresAt time.Time
	httpClient     *resty.Client
}

// NewTokenProvider creates a new JWT token provider for Apple authentication
func NewTokenProvider(config Config, logger *zap.Logger) (*TokenProvider, error) {
	logger.Info("Creating Apple JWT token provider",
		zap.String("api_type", config.APIType),
		zap.String("client_id", config.ClientID))

	// Parse private key
	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &TokenProvider{
		clientID:   config.ClientID,
		keyID:      config.KeyID,
		privateKey: privateKey,
		scope:      config.Scope,
		logger:     logger,
		httpClient: resty.New().SetTimeout(30 * time.Second),
	}, nil
}

// GetToken returns a valid access token, refreshing if necessary
func (p *TokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mutex.RLock()
	if p.accessToken != "" && time.Now().Before(p.tokenExpiresAt.Add(-5*time.Minute)) {
		token := p.accessToken
		p.mutex.RUnlock()
		return token, nil
	}
	p.mutex.RUnlock()

	// Need to refresh token
	return p.refreshToken(ctx)
}

// IsValid returns true if the current token is valid
func (p *TokenProvider) IsValid() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.accessToken != "" && time.Now().Before(p.tokenExpiresAt)
}

// ForceRefresh forces a token refresh
func (p *TokenProvider) ForceRefresh(ctx context.Context) error {
	_, err := p.refreshToken(ctx)
	return err
}

// refreshToken performs the OAuth 2.0 client credentials flow
func (p *TokenProvider) refreshToken(ctx context.Context) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Double-check after acquiring lock
	if p.accessToken != "" && time.Now().Before(p.tokenExpiresAt.Add(-5*time.Minute)) {
		return p.accessToken, nil
	}

	p.logger.Debug("Refreshing access token")

	// Generate JWT assertion
	jwtToken, err := p.generateJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Prepare OAuth request
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope,omitempty"`
	}

	formData := map[string]string{
		"grant_type":            GrantType,
		"client_id":             p.clientID,
		"client_assertion_type": ClientAssertionType,
		"client_assertion":      jwtToken,
		"scope":                 p.scope,
	}

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		SetResult(&tokenResponse).
		Post(TokenEndpoint)

	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.IsError() {
		p.logger.Error("Authentication failed",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response", resp.String()))
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Update token information
	p.accessToken = tokenResponse.AccessToken
	p.tokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	p.logger.Info("Successfully refreshed access token",
		zap.String("token_type", tokenResponse.TokenType),
		zap.Int("expires_in", tokenResponse.ExpiresIn),
		zap.String("scope", tokenResponse.Scope))

	return p.accessToken, nil
}

// generateJWT creates a JWT token for authentication
func (p *TokenProvider) generateJWT() (string, error) {
	now := time.Now()

	// Generate unique JWT ID
	jti := fmt.Sprintf("%s-%d", p.clientID, now.Unix())

	claims := Claims{
		Issuer:    p.clientID,
		Subject:   p.clientID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(20 * time.Minute).Unix(), // Apple requires max 20 minutes
		Audience:  TokenAudience,
		JTI:       jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = p.keyID

	return token.SignedString(p.privateKey)
}

// parsePrivateKey parses a PEM-encoded ECDSA private key (Apple AXM API requirement)
func parsePrivateKey(keyData string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Try different parsing methods based on block type and fallback on errors
	switch block.Type {
	case "EC PRIVATE KEY":
		// Try EC private key format first
		if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
			return key, nil
		}
		// Fallback to PKCS#8 if EC parsing fails
		fallthrough
	case "PRIVATE KEY":
		// PKCS#8 format - extract ECDSA key
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
		}

		ecdsaKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not an ECDSA key (found %T). Apple AXM API requires ECDSA keys", key)
		}
		return ecdsaKey, nil
	default:
		// For unknown types, try both parsing methods
		// First try PKCS#8 (most common)
		if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if ecdsaKey, ok := key.(*ecdsa.PrivateKey); ok {
				return ecdsaKey, nil
			}
		}
		// Then try EC private key format
		if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
			return key, nil
		}

		return nil, fmt.Errorf("unsupported key type %q or invalid key format. Apple AXM API requires ECDSA keys in 'EC PRIVATE KEY' or 'PRIVATE KEY' (PKCS#8) format", block.Type)
	}
}
