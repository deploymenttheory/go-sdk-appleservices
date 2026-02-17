package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"resty.dev/v3"
)

func TestNewJWTAuth(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	config := JWTAuthConfig{
		KeyID:      "test-key-id",
		IssuerID:   "test-issuer",
		PrivateKey: privateKey,
	}

	auth := NewJWTAuth(config)

	if auth == nil {
		t.Fatal("NewJWTAuth returned nil")
	}

	if auth.keyID != "test-key-id" {
		t.Errorf("keyID = %v, want 'test-key-id'", auth.keyID)
	}

	if auth.issuerID != "test-issuer" {
		t.Errorf("issuerID = %v, want 'test-issuer'", auth.issuerID)
	}

	if auth.audience != DefaultJWTAudience {
		t.Errorf("audience = %v, want %v", auth.audience, DefaultJWTAudience)
	}

	if auth.scope != ScopeBusinessAPI {
		t.Errorf("scope = %v, want %v", auth.scope, ScopeBusinessAPI)
	}
}

func TestNewJWTAuth_CustomAudienceAndScope(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	config := JWTAuthConfig{
		KeyID:      "test-key-id",
		IssuerID:   "test-issuer",
		PrivateKey: privateKey,
		Audience:   "custom-audience",
		Scope:      ScopeSchoolAPI,
	}

	auth := NewJWTAuth(config)

	if auth.audience != "custom-audience" {
		t.Errorf("audience = %v, want 'custom-audience'", auth.audience)
	}

	if auth.scope != ScopeSchoolAPI {
		t.Errorf("scope = %v, want %v", auth.scope, ScopeSchoolAPI)
	}
}

func TestJWTAuth_GenerateClientAssertion_ECDSA(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:      "test-key-id",
		issuerID:   "test-issuer",
		privateKey: privateKey,
		audience:   DefaultJWTAudience,
		scope:      ScopeBusinessAPI,
	}

	assertion, err := auth.generateClientAssertion()
	if err != nil {
		t.Fatalf("generateClientAssertion failed: %v", err)
	}

	if assertion == "" {
		t.Error("generateClientAssertion returned empty string")
	}

	// Parse the JWT to verify structure
	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse generated JWT: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	// Verify required claims
	if claims["iss"] != "test-issuer" {
		t.Errorf("iss claim = %v, want 'test-issuer'", claims["iss"])
	}

	if claims["sub"] != "test-issuer" {
		t.Errorf("sub claim = %v, want 'test-issuer'", claims["sub"])
	}

	if claims["aud"] != DefaultOAuthTokenEndpoint {
		t.Errorf("aud claim = %v, want %v", claims["aud"], DefaultOAuthTokenEndpoint)
	}

	// Verify kid header
	if token.Header["kid"] != "test-key-id" {
		t.Errorf("kid header = %v, want 'test-key-id'", token.Header["kid"])
	}

	// Verify signing method
	if token.Method != jwt.SigningMethodES256 {
		t.Errorf("signing method = %v, want ES256", token.Method)
	}
}

func TestJWTAuth_GenerateClientAssertion_RSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:      "test-key-id",
		issuerID:   "test-issuer",
		privateKey: privateKey,
		audience:   DefaultJWTAudience,
		scope:      ScopeBusinessAPI,
	}

	assertion, err := auth.generateClientAssertion()
	if err != nil {
		t.Fatalf("generateClientAssertion failed: %v", err)
	}

	if assertion == "" {
		t.Error("generateClientAssertion returned empty string")
	}

	// Parse the JWT to verify structure
	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse generated JWT: %v", err)
	}

	// Verify signing method
	if token.Method != jwt.SigningMethodRS256 {
		t.Errorf("signing method = %v, want RS256", token.Method)
	}
}

func TestJWTAuth_GenerateClientAssertion_UnsupportedKeyType(t *testing.T) {
	auth := &JWTAuth{
		keyID:      "test-key-id",
		issuerID:   "test-issuer",
		privateKey: "not-a-valid-key",
		audience:   DefaultJWTAudience,
		scope:      ScopeBusinessAPI,
	}

	_, err := auth.generateClientAssertion()
	if err == nil {
		t.Error("Expected error for unsupported key type, got nil")
	}
}

func TestJWTAuth_ForceRefresh(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:       "test-key-id",
		issuerID:    "test-issuer",
		privateKey:  privateKey,
		accessToken: "existing-token",
		tokenExpiry: time.Now().Add(1 * time.Hour),
	}

	// Verify token exists
	if auth.accessToken == "" {
		t.Fatal("accessToken should not be empty before ForceRefresh")
	}

	auth.ForceRefresh()

	// Verify token was cleared
	if auth.accessToken != "" {
		t.Error("accessToken should be empty after ForceRefresh")
	}

	if !auth.tokenExpiry.IsZero() {
		t.Error("tokenExpiry should be zero after ForceRefresh")
	}
}

func TestNewAPIKeyAuth(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		header     string
		wantHeader string
	}{
		{
			name:       "Default header",
			apiKey:     "test-api-key",
			header:     "",
			wantHeader: "Authorization",
		},
		{
			name:       "Custom header",
			apiKey:     "test-api-key",
			header:     "X-API-Key",
			wantHeader: "X-API-Key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAPIKeyAuth(tt.apiKey, tt.header)

			if auth == nil {
				t.Fatal("NewAPIKeyAuth returned nil")
			}

			if auth.apiKey != tt.apiKey {
				t.Errorf("apiKey = %v, want %v", auth.apiKey, tt.apiKey)
			}

			if auth.header != tt.wantHeader {
				t.Errorf("header = %v, want %v", auth.header, tt.wantHeader)
			}
		})
	}
}

func TestAPIKeyAuth_ApplyAuth_AuthorizationHeader(t *testing.T) {
	auth := NewAPIKeyAuth("test-key", "Authorization")

	req := resty.New().R()
	err := auth.ApplyAuth(req)

	if err != nil {
		t.Fatalf("ApplyAuth failed: %v", err)
	}

	// Verify auth token was set (resty internal, can't easily verify directly)
	// Just ensuring no panic/error
}

func TestAPIKeyAuth_ApplyAuth_CustomHeader(t *testing.T) {
	auth := NewAPIKeyAuth("test-key", "X-Custom-API-Key")

	req := resty.New().R()
	err := auth.ApplyAuth(req)

	if err != nil {
		t.Fatalf("ApplyAuth failed: %v", err)
	}
}

func TestTokenResponse(t *testing.T) {
	tokenResp := &TokenResponse{
		AccessToken: "test-token-12345",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       "business.api",
	}

	if tokenResp.AccessToken != "test-token-12345" {
		t.Errorf("AccessToken = %v", tokenResp.AccessToken)
	}

	if tokenResp.TokenType != "Bearer" {
		t.Errorf("TokenType = %v", tokenResp.TokenType)
	}

	if tokenResp.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %v, want 3600", tokenResp.ExpiresIn)
	}

	if tokenResp.Scope != "business.api" {
		t.Errorf("Scope = %v", tokenResp.Scope)
	}
}

func TestJWTAuthConfig_Defaults(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	config := JWTAuthConfig{
		KeyID:      "key",
		IssuerID:   "issuer",
		PrivateKey: privateKey,
		// Audience and Scope intentionally omitted
	}

	auth := NewJWTAuth(config)

	// Should set defaults
	if auth.audience != DefaultJWTAudience {
		t.Errorf("Default audience = %v, want %v", auth.audience, DefaultJWTAudience)
	}

	if auth.scope != ScopeBusinessAPI {
		t.Errorf("Default scope = %v, want %v", auth.scope, ScopeBusinessAPI)
	}
}

func TestJWTAuth_ConcurrentAccess(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:       "test-key",
		issuerID:    "test-issuer",
		privateKey:  privateKey,
		accessToken: "cached-token",
		tokenExpiry: time.Now().Add(1 * time.Hour),
		httpClient:  resty.New(),
	}

	// Test concurrent ForceRefresh calls (should not panic)
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			auth.ForceRefresh()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if auth.accessToken != "" {
		t.Error("accessToken should be empty after concurrent ForceRefresh")
	}
}

func TestJWTAuth_GenerateClientAssertion_ClaimsStructure(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:      "test-key-id",
		issuerID:   "test-issuer-id",
		privateKey: privateKey,
		audience:   DefaultJWTAudience,
		scope:      ScopeBusinessAPI,
	}

	assertion, err := auth.generateClientAssertion()
	if err != nil {
		t.Fatalf("generateClientAssertion failed: %v", err)
	}

	// Parse without verification just to check structure
	token, _, err := jwt.NewParser().ParseUnverified(assertion, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	// Verify all required claims exist
	requiredClaims := []string{"iss", "sub", "aud", "iat", "exp", "jti"}
	for _, claim := range requiredClaims {
		if _, exists := claims[claim]; !exists {
			t.Errorf("Missing required claim: %s", claim)
		}
	}

	// Verify iat and exp are numeric
	if _, ok := claims["iat"].(float64); !ok {
		t.Error("iat claim is not numeric")
	}

	if _, ok := claims["exp"].(float64); !ok {
		t.Error("exp claim is not numeric")
	}
}

func TestAPIKeyAuth_ApplyAuth_BothHeaderTypes(t *testing.T) {
	tests := []struct {
		name   string
		header string
		apiKey string
	}{
		{
			name:   "Authorization header",
			header: "Authorization",
			apiKey: "key-12345",
		},
		{
			name:   "Custom X-API-Key header",
			header: "X-API-Key",
			apiKey: "custom-key-67890",
		},
		{
			name:   "Empty header defaults to Authorization",
			header: "",
			apiKey: "default-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAPIKeyAuth(tt.apiKey, tt.header)
			req := resty.New().R()

			err := auth.ApplyAuth(req)
			if err != nil {
				t.Errorf("ApplyAuth failed: %v", err)
			}
		})
	}
}

func TestAuthProviderInterface(t *testing.T) {
	// Test that both auth types implement the interface
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	var _ AuthProvider = NewJWTAuth(JWTAuthConfig{
		KeyID:      "key",
		IssuerID:   "issuer",
		PrivateKey: privateKey,
	})

	var _ AuthProvider = NewAPIKeyAuth("key", "header")
}

func TestJWTAuth_GenerateClientAssertion_ExpirationTime(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:      "key",
		issuerID:   "issuer",
		privateKey: privateKey,
	}

	beforeGeneration := time.Now()
	assertion, err := auth.generateClientAssertion()
	if err != nil {
		t.Fatalf("generateClientAssertion failed: %v", err)
	}

	// Parse to check expiration
	token, _, err := jwt.NewParser().ParseUnverified(assertion, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("exp claim is not numeric")
	}

	expTime := time.Unix(int64(exp), 0)
	
	// Should expire in approximately 180 days
	expectedExpiration := beforeGeneration.Add(180 * 24 * time.Hour)
	diff := expTime.Sub(expectedExpiration).Abs()
	
	// Allow 1 minute variance for test execution time
	if diff > time.Minute {
		t.Errorf("Expiration time differs from expected by %v", diff)
	}
}

func TestJWTAuth_GenerateClientAssertion_JTI_Uniqueness(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	auth := &JWTAuth{
		keyID:      "key",
		issuerID:   "issuer",
		privateKey: privateKey,
	}

	// Generate multiple assertions
	assertions := make(map[string]bool)
	for i := 0; i < 10; i++ {
		assertion, err := auth.generateClientAssertion()
		if err != nil {
			t.Fatalf("generateClientAssertion failed: %v", err)
		}

		if assertions[assertion] {
			t.Error("Generated duplicate JWT assertion")
		}
		assertions[assertion] = true

		// Small sleep to ensure different timestamps
		time.Sleep(time.Microsecond)
	}

	if len(assertions) != 10 {
		t.Errorf("Generated %d unique assertions, want 10", len(assertions))
	}
}

func TestJWTAuthConfig_AllFieldsSet(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	config := JWTAuthConfig{
		KeyID:      "key-123",
		IssuerID:   "issuer-456",
		PrivateKey: privateKey,
		Audience:   "custom-aud",
		Scope:      "custom.scope",
	}

	auth := NewJWTAuth(config)

	if auth.keyID != config.KeyID {
		t.Errorf("keyID not set correctly")
	}
	if auth.issuerID != config.IssuerID {
		t.Errorf("issuerID not set correctly")
	}
	if auth.audience != config.Audience {
		t.Errorf("audience not set correctly")
	}
	if auth.scope != config.Scope {
		t.Errorf("scope not set correctly")
	}
}
