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

// JWTAuth implements direct JWT Bearer authentication for the Apple Notary API.
// The Notary API uses App Store Connect API keys — a signed JWT is used directly
// as a Bearer token without an OAuth token exchange step.
type JWTAuth struct {
	keyID       string
	issuerID    string
	privateKey  any // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	audience    string
	token       string
	tokenExpiry time.Time
	mutex       sync.RWMutex
}

// JWTAuthConfig holds configuration for JWT authentication
type JWTAuthConfig struct {
	KeyID      string
	IssuerID   string
	PrivateKey any    // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	Audience   string // Usually "appstoreconnect-v1"
}

// NewJWTAuth creates a new direct JWT authentication provider
func NewJWTAuth(config JWTAuthConfig) *JWTAuth {
	if config.Audience == "" {
		config.Audience = DefaultJWTAudience
	}

	return &JWTAuth{
		keyID:      config.KeyID,
		issuerID:   config.IssuerID,
		privateKey: config.PrivateKey,
		audience:   config.Audience,
	}
}

// ApplyAuth applies JWT Bearer authentication to the request
func (j *JWTAuth) ApplyAuth(req *resty.Request) error {
	token, err := j.getToken()
	if err != nil {
		return fmt.Errorf("failed to get JWT token: %w", err)
	}

	req.SetAuthToken(token)
	return nil
}

// getToken returns a valid JWT, generating a new one if expired
func (j *JWTAuth) getToken() (string, error) {
	j.mutex.RLock()
	if j.token != "" && time.Now().Before(j.tokenExpiry.Add(-5*time.Minute)) {
		token := j.token
		j.mutex.RUnlock()
		return token, nil
	}
	j.mutex.RUnlock()

	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Double-check after acquiring write lock
	if j.token != "" && time.Now().Before(j.tokenExpiry.Add(-5*time.Minute)) {
		return j.token, nil
	}

	return j.generateToken()
}

// generateToken creates a signed JWT for the Notary API
func (j *JWTAuth) generateToken() (string, error) {
	now := time.Now()
	expiry := now.Add(20 * time.Minute)

	claims := jwt.MapClaims{
		"iss": j.issuerID,
		"iat": now.Unix(),
		"exp": expiry.Unix(),
		"aud": j.audience,
	}

	var signingMethod jwt.SigningMethod
	switch j.privateKey.(type) {
	case *ecdsa.PrivateKey:
		signingMethod = jwt.SigningMethodES256
	case *rsa.PrivateKey:
		signingMethod = jwt.SigningMethodRS256
	default:
		return "", fmt.Errorf("unsupported private key type: %T", j.privateKey)
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = j.keyID

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	j.token = tokenString
	j.tokenExpiry = expiry

	return j.token, nil
}

// ForceRefresh forces a token refresh on the next request
func (j *JWTAuth) ForceRefresh() {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.token = ""
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
