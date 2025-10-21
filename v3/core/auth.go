package axm

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"net/url"
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
		config.Audience = "appstoreconnect-v1"
	}
	if config.Scope == "" {
		config.Scope = "business.api"
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

	// Generate client assertion
	clientAssertion, err := j.generateClientAssertion()
	if err != nil {
		return "", fmt.Errorf("failed to generate client assertion: %w", err)
	}

	// Exchange for access token
	tokenResp, err := j.exchangeForAccessToken(clientAssertion)
	if err != nil {
		return "", fmt.Errorf("failed to exchange for access token: %w", err)
	}

	// Store the token
	j.accessToken = tokenResp.AccessToken
	j.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return j.accessToken, nil
}

// generateClientAssertion creates a JWT client assertion for OAuth 2.0 authentication
func (j *JWTAuth) generateClientAssertion() (string, error) {
	now := time.Now()

	// Create client assertion claims as per Apple's OAuth 2.0 spec
	claims := jwt.MapClaims{
		"iss": j.issuerID,                                       // team_id (issuer)
		"sub": j.issuerID,                                       // client_id (subject) - same as issuer for Apple
		"aud": "https://account.apple.com/auth/oauth2/v2/token", // OAuth 2.0 token endpoint
		"iat": now.Unix(),
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
		Post("https://account.apple.com/auth/oauth2/v2/token")

	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return &tokenResp, nil
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

// OAuth2Auth implements OAuth 2.0 authentication for Apple School and Business Manager API
type OAuth2Auth struct {
	clientID   string
	teamID     string
	keyID      string
	privateKey *ecdsa.PrivateKey
	scope      string
	tokenURL   string
	httpClient *resty.Client

	// Token management
	mu          sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

// OAuth2Config holds configuration for OAuth 2.0 authentication
type OAuth2Config struct {
	ClientID   string            // Your client ID from Apple
	TeamID     string            // Your team ID
	KeyID      string            // Your key ID from Apple
	PrivateKey *ecdsa.PrivateKey // Your ECDSA private key
	Scope      string            // "business.api" or "school.api"
	TokenURL   string            // OAuth token endpoint (optional, defaults to Apple's)
	HTTPClient *resty.Client     // HTTP client for token requests (optional)
}

// OAuth2TokenResponse represents the OAuth 2.0 token response
type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// NewOAuth2Auth creates a new OAuth 2.0 authentication provider
func NewOAuth2Auth(config OAuth2Config) (*OAuth2Auth, error) {
	// Validate required fields
	if config.ClientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}
	if config.TeamID == "" {
		return nil, fmt.Errorf("team ID is required")
	}
	if config.KeyID == "" {
		return nil, fmt.Errorf("key ID is required")
	}
	if config.PrivateKey == nil {
		return nil, fmt.Errorf("private key is required")
	}
	if config.Scope == "" {
		return nil, fmt.Errorf("scope is required (business.api or school.api)")
	}

	// Set defaults
	if config.TokenURL == "" {
		config.TokenURL = DefaultOAuthTokenURL
	}
	if config.HTTPClient == nil {
		config.HTTPClient = resty.New()
	}

	return &OAuth2Auth{
		clientID:   config.ClientID,
		teamID:     config.TeamID,
		keyID:      config.KeyID,
		privateKey: config.PrivateKey,
		scope:      config.Scope,
		tokenURL:   config.TokenURL,
		httpClient: config.HTTPClient,
	}, nil
}

// ApplyAuth applies OAuth 2.0 authentication to the request
func (o *OAuth2Auth) ApplyAuth(req *resty.Request) error {
	token, err := o.getValidAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	req.SetAuthToken(token)
	return nil
}

// getValidAccessToken returns a valid access token, refreshing if necessary
func (o *OAuth2Auth) getValidAccessToken() (string, error) {
	o.mu.RLock()
	if o.accessToken != "" && time.Now().Before(o.tokenExpiry.Add(-30*time.Second)) {
		token := o.accessToken
		o.mu.RUnlock()
		return token, nil
	}
	o.mu.RUnlock()

	// Need to refresh token
	return o.refreshAccessToken()
}

// refreshAccessToken obtains a new access token using OAuth 2.0 flow
func (o *OAuth2Auth) refreshAccessToken() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Double-check in case another goroutine already refreshed
	if o.accessToken != "" && time.Now().Before(o.tokenExpiry.Add(-30*time.Second)) {
		return o.accessToken, nil
	}

	// Generate client assertion JWT
	clientAssertion, err := o.generateClientAssertion()
	if err != nil {
		return "", fmt.Errorf("failed to generate client assertion: %w", err)
	}

	// Prepare token request
	formData := url.Values{
		"grant_type":            {"client_credentials"},
		"client_id":             {o.clientID},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"client_assertion":      {clientAssertion},
		"scope":                 {o.scope},
	}

	// Make token request
	var tokenResp OAuth2TokenResponse
	resp, err := o.httpClient.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Host", "account.apple.com").
		SetFormDataFromValues(formData).
		SetResult(&tokenResp).
		Post(o.tokenURL)

	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Update stored token
	o.accessToken = tokenResp.AccessToken
	o.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return o.accessToken, nil
}

// generateClientAssertion creates a JWT client assertion for OAuth 2.0
func (o *OAuth2Auth) generateClientAssertion() (string, error) {
	now := time.Now()

	// JWT header
	token := jwt.New(jwt.SigningMethodES256)
	token.Header["kid"] = o.keyID

	// JWT claims
	claims := jwt.MapClaims{
		"iss": o.teamID,                             // Issuer (team ID)
		"sub": o.clientID,                           // Subject (client ID)
		"aud": DefaultOAuthAudience,                 // Audience
		"iat": now.Unix(),                           // Issued at
		"exp": now.Add(180 * 24 * time.Hour).Unix(), // Expires (max 180 days)
		"jti": fmt.Sprintf("%d", now.UnixNano()),    // JWT ID (unique identifier)
	}
	token.Claims = claims

	// Sign the token
	tokenString, err := token.SignedString(o.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign client assertion: %w", err)
	}

	return tokenString, nil
}

// ForceRefresh forces a token refresh on the next request
func (o *OAuth2Auth) ForceRefresh() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.accessToken = ""
	o.tokenExpiry = time.Time{}
}

// TokenInfo represents information about the current access token
type TokenInfo struct {
	HasToken     bool      `json:"has_token"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
	IsExpired    bool      `json:"is_expired"`
	NeedsRefresh bool      `json:"needs_refresh"`
}

// GetTokenInfo returns information about the current access token
func (o *OAuth2Auth) GetTokenInfo() TokenInfo {
	o.mu.RLock()
	defer o.mu.RUnlock()

	info := TokenInfo{
		HasToken: o.accessToken != "",
	}

	if info.HasToken {
		info.ExpiresAt = o.tokenExpiry
		info.ExpiresIn = int64(time.Until(o.tokenExpiry).Seconds())
		info.IsExpired = time.Now().After(o.tokenExpiry)
		info.NeedsRefresh = time.Now().After(o.tokenExpiry.Add(-30 * time.Second))
	}

	return info
}
