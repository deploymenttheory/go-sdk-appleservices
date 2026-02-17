package client

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"resty.dev/v3"
)

// AuthProvider interface for different authentication methods
type AuthProvider interface {
	ApplyAuth(req *resty.Request) error
}

// JWTAuth implements OAuth 2.0 JWT-based authentication for Apple Business Manager API
type JWTAuth struct {
	keyID       string
	issuerID    string
	privateKey  any // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	audience    string
	scope       string
	accessToken string
	tokenExpiry time.Time
	mutex       sync.RWMutex
	httpClient  *resty.Client
}

// JWTAuthConfig holds configuration for JWT authentication
type JWTAuthConfig struct {
	KeyID      string
	IssuerID   string
	PrivateKey any    // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	Audience   string // Usually "appstoreconnect-v1"
	Scope      string // "business.api" or "school.api"
}

// NewJWTAuth creates a new OAuth 2.0 JWT authentication provider
func NewJWTAuth(config JWTAuthConfig) *JWTAuth {
	if config.Audience == "" {
		config.Audience = DefaultJWTAudience
	}
	if config.Scope == "" {
		config.Scope = ScopeBusinessAPI
	}

	return &JWTAuth{
		keyID:      config.KeyID,
		issuerID:   config.IssuerID,
		privateKey: config.PrivateKey,
		audience:   config.Audience,
		scope:      config.Scope,
		httpClient: resty.New(),
	}
}

// ApplyAuth applies OAuth 2.0 authentication to the request
func (j *JWTAuth) ApplyAuth(req *resty.Request) error {
	accessToken, err := j.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	req.SetAuthToken(accessToken)
	return nil
}

// getAccessToken returns a valid access token, refreshing if necessary
func (j *JWTAuth) getAccessToken() (string, error) {
	j.mutex.RLock()
	if j.accessToken != "" && time.Now().Before(j.tokenExpiry.Add(-5*time.Minute)) {
		token := j.accessToken
		j.mutex.RUnlock()
		return token, nil
	}
	j.mutex.RUnlock()

	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Double-check after acquiring write lock
	if j.accessToken != "" && time.Now().Before(j.tokenExpiry.Add(-5*time.Minute)) {
		return j.accessToken, nil
	}

	clientAssertion, err := j.generateClientAssertion()
	if err != nil {
		return "", fmt.Errorf("failed to generate client assertion: %w", err)
	}

	tokenResp, err := j.exchangeForAccessToken(clientAssertion)
	if err != nil {
		return "", fmt.Errorf("failed to exchange for access token: %w", err)
	}

	j.accessToken = tokenResp.AccessToken
	j.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return j.accessToken, nil
}

// generateClientAssertion creates a JWT client assertion for OAuth 2.0 authentication
func (j *JWTAuth) generateClientAssertion() (string, error) {
	now := time.Now()

	// Create client assertion claims as per Apple's OAuth 2.0 spec
	claims := jwt.MapClaims{
		"iss": j.issuerID,                           // team_id (issuer)
		"sub": j.issuerID,                           // client_id (subject) - same as issuer for Apple
		"aud": DefaultOAuthTokenEndpoint,            // OAuth 2.0 token endpoint
		"iat": now.Unix(),                           // Issued at time
		"exp": now.Add(180 * 24 * time.Hour).Unix(), // Max 180 days as per Apple docs
		"jti": fmt.Sprintf("%d", now.UnixNano()),    // Unique identifier
	}

	// Determine signing method based on key type
	var signingMethod jwt.SigningMethod
	switch j.privateKey.(type) {
	case *ecdsa.PrivateKey:
		signingMethod = jwt.SigningMethodES256 // ES256 for ECDSA keys
	case *rsa.PrivateKey:
		signingMethod = jwt.SigningMethodRS256 // RS256 for RSA keys (fallback)
	default:
		return "", fmt.Errorf("unsupported private key type: %T", j.privateKey)
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = j.keyID

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT client assertion: %w", err)
	}

	return tokenString, nil
}

// TokenResponse represents the OAuth 2.0 token response from Apple
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// exchangeForAccessToken exchanges the client assertion for an access token
func (j *JWTAuth) exchangeForAccessToken(clientAssertion string) (*TokenResponse, error) {
	var tokenResp TokenResponse
	resp, err := j.httpClient.R().
		SetFormData(map[string]string{
			"grant_type":            "client_credentials",
			"client_id":             j.issuerID,
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      clientAssertion,
			"scope":                 j.scope,
		}).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Host", "account.apple.com").
		SetResult(&tokenResp).
		Post(DefaultOAuthTokenEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return &tokenResp, nil
}

// ForceRefresh forces a token refresh on the next request
func (j *JWTAuth) ForceRefresh() {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.accessToken = ""
	j.tokenExpiry = time.Time{}
}

// APIKeyAuth implements simple API key authentication
type APIKeyAuth struct {
	apiKey string
	header string
}

// NewAPIKeyAuth creates a new API key authentication provider
func NewAPIKeyAuth(apiKey, header string) *APIKeyAuth {
	if header == "" {
		header = "Authorization"
	}
	return &APIKeyAuth{
		apiKey: apiKey,
		header: header,
	}
}

// ApplyAuth applies API key authentication to the request
func (a *APIKeyAuth) ApplyAuth(req *resty.Request) error {
	if a.header == "Authorization" {
		req.SetAuthToken(a.apiKey)
	} else {
		req.SetHeader(a.header, a.apiKey)
	}
	return nil
}
