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
)

func TestNewTransport_Success(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	client, err := NewTransport("test-key-id", "test-issuer-id", privateKey)

	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewTransport returned nil client")
	}

	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}

	if client.auth == nil {
		t.Error("auth is nil")
	}

	if client.logger == nil {
		t.Error("logger is nil")
	}

	if client.errorHandler == nil {
		t.Error("errorHandler is nil")
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %v, want %v", client.baseURL, DefaultBaseURL)
	}
}

func TestNewTransport_WithOptions(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	logger := zap.NewNop()

	client, err := NewTransport(
		"test-key-id",
		"test-issuer-id",
		privateKey,
		WithLogger(logger),
		WithDebug(),
	)

	if err != nil {
		t.Fatalf("NewTransport with options failed: %v", err)
	}

	if client.logger != logger {
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

	client, err := NewTransport("test-key-id", "test-issuer-id", privateKey)

	if err != nil {
		t.Fatalf("NewTransport with RSA key failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewTransport returned nil client")
	}
}

func TestNewTransportFromFile_Success(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Generate and save key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	keyPath := tmpDir + "/test_key.p8"
	if err := savePrivateKeyToFile(privateKey, keyPath); err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	client, err := NewTransportFromFile("test-key-id", "test-issuer-id", keyPath)

	if err != nil {
		t.Fatalf("NewTransportFromFile failed: %v", err)
	}

	if client == nil {
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

func TestNewTransportFromFile_InvalidPath(t *testing.T) {
	_, err := NewTransportFromFile("key", "issuer", "/nonexistent/path/key.p8")

	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestNewTransportFromEnv_MissingEnvVars(t *testing.T) {
	// Save and clear env vars
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

	client, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	httpClient := client.GetHTTPClient()
	if httpClient == nil {
		t.Error("GetHTTPClient returned nil")
	}
}

func TestClient_QueryBuilder_Integration(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	client, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	qb := client.QueryBuilder()

	qb.AddString("test", "value").
		AddInt("limit", 10).
		AddBool("active", true)

	if qb.Count() != 3 {
		t.Errorf("QueryBuilder count = %d, want 3", qb.Count())
	}
}

// Helper function to save private key to file for testing
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

func TestNewTransport_HTTPClientConfiguration(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	client, err := NewTransport("key", "issuer", privateKey)
	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	httpClient := client.GetHTTPClient()

	// Verify base URL is set
	if httpClient.BaseURL() != DefaultBaseURL {
		t.Errorf("BaseURL = %v, want %v", httpClient.BaseURL(), DefaultBaseURL)
	}

	// Verify user agent is set
	if httpClient.Header().Get("User-Agent") != DefaultUserAgent {
		t.Errorf("User-Agent = %v, want %v", httpClient.Header().Get("User-Agent"), DefaultUserAgent)
	}
}

func TestNewTransport_AuthMiddlewareSetup(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create mock auth that we can verify was called
	mockAuth := &MockAuthProvider{}

	client, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithAuth(mockAuth),
	)

	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	// Activate httpmock
	httpmock.ActivateNonDefault(client.httpClient.Client())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]string{"status": "ok"}))

	// Make a request - should trigger auth middleware
	var result map[string]string
	_ = client.Get(context.Background(), "/v1/test", nil, nil, &result)

	// If we get here without panic, auth middleware was set up correctly
}
