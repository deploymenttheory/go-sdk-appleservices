package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestNewTransport_Success(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	c, err := NewTransport("test-key-id", "test-issuer-id", privateKey)

	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	if c == nil {
		t.Fatal("NewTransport returned nil client")
	}

	if c.httpClient == nil {
		t.Error("httpClient is nil")
	}

	if c.auth == nil {
		t.Error("auth is nil")
	}

	if c.logger == nil {
		t.Error("logger is nil")
	}

	if c.errorHandler == nil {
		t.Error("errorHandler is nil")
	}

	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %v, want %v", c.baseURL, DefaultBaseURL)
	}
}

func TestNewTransport_WithOptions(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	logger := zap.NewNop()

	c, err := NewTransport(
		"test-key-id",
		"test-issuer-id",
		privateKey,
		WithLogger(logger),
		WithDebug(),
	)

	if err != nil {
		t.Fatalf("NewTransport with options failed: %v", err)
	}

	if c.logger != logger {
		t.Error("Logger option not applied")
	}
}

func TestNewTransport_MissingKeyID(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	_, err = NewTransport("", "test-issuer-id", privateKey)

	if err == nil {
		t.Error("Expected error for missing keyID, got nil")
	}
}

func TestNewTransport_MissingIssuerID(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	_, err = NewTransport("test-key-id", "", privateKey)

	if err == nil {
		t.Error("Expected error for missing issuerID, got nil")
	}
}

func TestNewTransport_NilPrivateKey(t *testing.T) {
	_, err := NewTransport("test-key-id", "test-issuer-id", nil)

	if err == nil {
		t.Error("Expected error for nil privateKey, got nil")
	}
}

func TestNewTransport_RSAKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	c, err := NewTransport("test-key-id", "test-issuer-id", privateKey)

	if err != nil {
		t.Fatalf("NewTransport with RSA key failed: %v", err)
	}

	if c == nil {
		t.Fatal("NewTransport returned nil client")
	}
}

func TestNewTransportFromFile_Success(t *testing.T) {
	tmpDir := t.TempDir()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	keyPath := tmpDir + "/test_key.p8"
	if err := savePrivateKeyToFile(privateKey, keyPath); err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	c, err := NewTransportFromFile("test-key-id", "test-issuer-id", keyPath)

	if err != nil {
		t.Fatalf("NewTransportFromFile failed: %v", err)
	}

	if c == nil {
		t.Fatal("NewTransportFromFile returned nil")
	}
}

func TestNewTransportFromFile_MissingParameters(t *testing.T) {
	tests := []struct {
		name           string
		keyID          string
		issuerID       string
		privateKeyPath string
	}{
		{"Missing keyID", "", "issuer", "/path/to/key"},
		{"Missing issuerID", "key", "", "/path/to/key"},
		{"Missing privateKeyPath", "key", "issuer", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransportFromFile(tt.keyID, tt.issuerID, tt.privateKeyPath)

			if err == nil {
				t.Error("Expected error for missing parameter, got nil")
			}
		})
	}
}

func TestNewTransportFromEnv_MissingEnvVars(t *testing.T) {
	oldKeyID := os.Getenv("APPLE_KEY_ID")
	oldIssuerID := os.Getenv("APPLE_ISSUER_ID")
	oldKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	os.Unsetenv("APPLE_KEY_ID")
	os.Unsetenv("APPLE_ISSUER_ID")
	os.Unsetenv("APPLE_PRIVATE_KEY_PATH")

	defer func() {
		if oldKeyID != "" {
			os.Setenv("APPLE_KEY_ID", oldKeyID)
		}
		if oldIssuerID != "" {
			os.Setenv("APPLE_ISSUER_ID", oldIssuerID)
		}
		if oldKeyPath != "" {
			os.Setenv("APPLE_PRIVATE_KEY_PATH", oldKeyPath)
		}
	}()

	_, err := NewTransportFromEnv()

	if err == nil {
		t.Error("Expected error for missing env vars, got nil")
	}
}

func TestClient_GetHTTPClient_NotNil(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	c, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	httpClient := c.GetHTTPClient()
	if httpClient == nil {
		t.Error("GetHTTPClient returned nil")
	}
}

func TestClient_QueryBuilder_Integration(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	c, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	qb := c.QueryBuilder()

	qb.AddString("test", "value").
		AddInt("limit", 10).
		AddBool("active", true)

	if qb.Count() != 3 {
		t.Errorf("QueryBuilder count = %d, want 3", qb.Count())
	}
}

func TestNewTransport_HTTPClientConfiguration(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	c, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	httpClient := c.GetHTTPClient()

	if httpClient.BaseURL() != DefaultBaseURL {
		t.Errorf("BaseURL = %v, want %v", httpClient.BaseURL(), DefaultBaseURL)
	}

	if httpClient.Header().Get("User-Agent") != DefaultUserAgent {
		t.Errorf("User-Agent = %v, want %v", httpClient.Header().Get("User-Agent"), DefaultUserAgent)
	}
}

// MockAuthProvider implements AuthProvider for testing.
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	return nil
}

func TestNewTransport_AuthMiddlewareSetup(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	mockAuth := &MockAuthProvider{}

	c, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithAuth(mockAuth),
	)

	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	httpmock.ActivateNonDefault(c.httpClient.Client())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://appstoreconnect.apple.com/notary/v2/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	var result map[string]string
	ctx := context.Background()
	_, _ = c.NewRequest(ctx).SetResult(&result).Get("/notary/v2/test")
}

// savePrivateKeyToFile is a test helper that serializes an ECDSA key to a PEM file.
func savePrivateKeyToFile(key *ecdsa.PrivateKey, path string) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	return os.WriteFile(path, pemData, 0600)
}
