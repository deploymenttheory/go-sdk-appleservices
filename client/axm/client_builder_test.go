package axm

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientBuilder(t *testing.T) {
	builder := NewClientBuilder()

	assert.Equal(t, "appstoreconnect-v1", builder.audience)
	assert.Equal(t, "https://api-business.apple.com/v1", builder.baseURL)
	assert.Equal(t, 30*time.Second, builder.timeout)
	assert.Equal(t, 3, builder.retryCount)
	assert.Equal(t, 1*time.Second, builder.retryWait)
	assert.Equal(t, "go-api-sdk-apple/1.0.0", builder.userAgent)
	assert.False(t, builder.debug)
}

func TestClientBuilder_WithMethods(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	builder := NewClientBuilder().
		WithJWTAuth("test-key-id", "test-issuer-id", privateKey).
		WithBaseURL("https://test.example.com").
		WithTimeout(60*time.Second).
		WithRetry(5, 2*time.Second).
		WithUserAgent("test-agent/1.0.0").
		WithDebug(true).
		WithAudience("test-audience")

	assert.Equal(t, "test-key-id", builder.keyID)
	assert.Equal(t, "test-issuer-id", builder.issuerID)
	assert.Equal(t, privateKey, builder.privateKey)
	assert.Equal(t, "https://test.example.com", builder.baseURL)
	assert.Equal(t, 60*time.Second, builder.timeout)
	assert.Equal(t, 5, builder.retryCount)
	assert.Equal(t, 2*time.Second, builder.retryWait)
	assert.Equal(t, "test-agent/1.0.0", builder.userAgent)
	assert.True(t, builder.debug)
	assert.Equal(t, "test-audience", builder.audience)
}

func TestClientBuilder_Validate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tests := []struct {
		name        string
		builder     *ClientBuilder
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			builder: NewClientBuilder().
				WithJWTAuth("key-id", "issuer-id", privateKey),
			expectError: false,
		},
		{
			name:        "missing key ID",
			builder:     NewClientBuilder().WithJWTAuth("", "issuer-id", privateKey),
			expectError: true,
			errorMsg:    "key ID is required",
		},
		{
			name:        "missing issuer ID",
			builder:     NewClientBuilder().WithJWTAuth("key-id", "", privateKey),
			expectError: true,
			errorMsg:    "issuer ID is required",
		},
		{
			name:        "missing private key",
			builder:     NewClientBuilder().WithJWTAuth("key-id", "issuer-id", nil),
			expectError: true,
			errorMsg:    "private key is required",
		},
		{
			name: "invalid timeout",
			builder: NewClientBuilder().
				WithJWTAuth("key-id", "issuer-id", privateKey).
				WithTimeout(-1 * time.Second),
			expectError: true,
			errorMsg:    "timeout must be positive",
		},
		{
			name: "invalid retry count",
			builder: NewClientBuilder().
				WithJWTAuth("key-id", "issuer-id", privateKey).
				WithRetry(-1, time.Second),
			expectError: true,
			errorMsg:    "retry count cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.builder.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClientBuilder_Clone(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	original := NewClientBuilder().
		WithJWTAuth("key-id", "issuer-id", privateKey).
		WithDebug(true)

	clone := original.Clone()

	// Verify clone has same values
	assert.Equal(t, original.keyID, clone.keyID)
	assert.Equal(t, original.issuerID, clone.issuerID)
	assert.Equal(t, original.privateKey, clone.privateKey)
	assert.Equal(t, original.debug, clone.debug)

	// Verify they are separate instances
	clone.WithDebug(false)
	assert.True(t, original.debug)
	assert.False(t, clone.debug)
}

func TestClientBuilder_WithJWTAuthFromEnv_MissingVars(t *testing.T) {
	// Clear environment variables
	originalKeyID := os.Getenv("APPLE_KEY_ID")
	originalIssuerID := os.Getenv("APPLE_ISSUER_ID")
	originalKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	defer func() {
		os.Setenv("APPLE_KEY_ID", originalKeyID)
		os.Setenv("APPLE_ISSUER_ID", originalIssuerID)
		os.Setenv("APPLE_PRIVATE_KEY_PATH", originalKeyPath)
	}()

	os.Unsetenv("APPLE_KEY_ID")
	os.Unsetenv("APPLE_ISSUER_ID")
	os.Unsetenv("APPLE_PRIVATE_KEY_PATH")

	builder := NewClientBuilder().WithJWTAuthFromEnv()

	// Should result in empty credentials
	assert.Empty(t, builder.keyID)
	assert.Empty(t, builder.issuerID)
	assert.Nil(t, builder.privateKey)

	// Validation should fail
	err := builder.Validate()
	assert.Error(t, err)
}

func TestNewClientFromEnv_MissingVars(t *testing.T) {
	// Clear environment variables
	originalKeyID := os.Getenv("APPLE_KEY_ID")
	originalIssuerID := os.Getenv("APPLE_ISSUER_ID")
	originalKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	defer func() {
		os.Setenv("APPLE_KEY_ID", originalKeyID)
		os.Setenv("APPLE_ISSUER_ID", originalIssuerID)
		os.Setenv("APPLE_PRIVATE_KEY_PATH", originalKeyPath)
	}()

	os.Unsetenv("APPLE_KEY_ID")
	os.Unsetenv("APPLE_ISSUER_ID")
	os.Unsetenv("APPLE_PRIVATE_KEY_PATH")

	client, err := NewClientFromEnv()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientFromFile_InvalidPath(t *testing.T) {
	client, err := NewClientFromFile("key-id", "issuer-id", "/nonexistent/path.p8")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientFromEnvWithOptions(t *testing.T) {
	// This test will fail without proper environment variables, but we can test the function exists
	client, err := NewClientFromEnvWithOptions(true, 30*time.Second, "test-agent")
	// We expect an error due to missing env vars, but the function should exist
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientFromFileWithOptions(t *testing.T) {
	client, err := NewClientFromFileWithOptions(
		"key-id", "issuer-id", "/nonexistent/path.p8",
		true, 30*time.Second, "test-agent",
	)
	assert.Error(t, err)
	assert.Nil(t, client)
}
