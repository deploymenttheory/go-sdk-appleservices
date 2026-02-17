package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestWithBaseURL(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customURL := "https://custom.api.example.com"
	client, err := NewTransport("key", "issuer", privateKey, WithBaseURL(customURL))

	if err != nil {
		t.Fatalf("NewTransport with WithBaseURL failed: %v", err)
	}

	if client.baseURL != customURL {
		t.Errorf("baseURL = %v, want %v", client.baseURL, customURL)
	}
}

func TestWithBaseURL_Empty(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithBaseURL(""))

	if err == nil {
		t.Error("Expected error for empty base URL, got nil")
	}
}

func TestWithLogger(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	logger := zap.NewExample()
	client, err := NewTransport("key", "issuer", privateKey, WithLogger(logger))

	if err != nil {
		t.Fatalf("NewTransport with WithLogger failed: %v", err)
	}

	if client.logger != logger {
		t.Error("Logger was not set correctly")
	}
}

func TestWithLogger_Nil(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithLogger(nil))

	if err == nil {
		t.Error("Expected error for nil logger, got nil")
	}
}

func TestWithTimeout(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	timeout := 60 * time.Second
	client, err := NewTransport("key", "issuer", privateKey, WithTimeout(timeout))

	if err != nil {
		t.Fatalf("NewTransport with WithTimeout failed: %v", err)
	}

	// Verify timeout was set (exact verification is internal to resty)
	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithTimeout_Negative(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithTimeout(-1*time.Second))

	if err == nil {
		t.Error("Expected error for negative timeout, got nil")
	}
}

func TestWithRetryCount(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	client, err := NewTransport("key", "issuer", privateKey, WithRetryCount(5))

	if err != nil {
		t.Fatalf("NewTransport with WithRetryCount failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithRetryCount_Negative(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithRetryCount(-1))

	if err == nil {
		t.Error("Expected error for negative retry count, got nil")
	}
}

func TestWithRetryWaitTime(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	waitTime := 2 * time.Second
	client, err := NewTransport("key", "issuer", privateKey, WithRetryWaitTime(waitTime))

	if err != nil {
		t.Fatalf("NewTransport with WithRetryWaitTime failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithRetryWaitTime_Negative(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithRetryWaitTime(-1*time.Second))

	if err == nil {
		t.Error("Expected error for negative retry wait time, got nil")
	}
}

func TestWithRetryMaxWaitTime(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	maxWaitTime := 30 * time.Second
	client, err := NewTransport("key", "issuer", privateKey, WithRetryMaxWaitTime(maxWaitTime))

	if err != nil {
		t.Fatalf("NewTransport with WithRetryMaxWaitTime failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithRetryMaxWaitTime_Negative(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithRetryMaxWaitTime(-1*time.Second))

	if err == nil {
		t.Error("Expected error for negative retry max wait time, got nil")
	}
}

func TestWithUserAgent(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customUA := "CustomApp/1.0.0"
	client, err := NewTransport("key", "issuer", privateKey, WithUserAgent(customUA))

	if err != nil {
		t.Fatalf("NewTransport with WithUserAgent failed: %v", err)
	}

	userAgent := client.httpClient.Header().Get("User-Agent")
	if userAgent != customUA {
		t.Errorf("User-Agent = %v, want %v", userAgent, customUA)
	}
}

func TestWithUserAgent_Empty(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithUserAgent(""))

	if err == nil {
		t.Error("Expected error for empty user agent, got nil")
	}
}

func TestWithCustomAgent(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customAgent := "MyApp/2.0"
	client, err := NewTransport("key", "issuer", privateKey, WithCustomAgent(customAgent))

	if err != nil {
		t.Fatalf("NewTransport with WithCustomAgent failed: %v", err)
	}

	userAgent := client.httpClient.Header().Get("User-Agent")
	expectedUA := DefaultUserAgent + "; " + customAgent
	if userAgent != expectedUA {
		t.Errorf("User-Agent = %v, want %v", userAgent, expectedUA)
	}
}

func TestWithDebug(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	client, err := NewTransport("key", "issuer", privateKey, WithDebug())

	if err != nil {
		t.Fatalf("NewTransport with WithDebug failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithAuth(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customAuth := &MockAuthProvider{}
	client, err := NewTransport("key", "issuer", privateKey, WithAuth(customAuth))

	if err != nil {
		t.Fatalf("NewTransport with WithAuth failed: %v", err)
	}

	if client.auth != customAuth {
		t.Error("Custom auth provider was not set")
	}
}

func TestWithAuth_Nil(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithAuth(nil))

	if err == nil {
		t.Error("Expected error for nil auth provider, got nil")
	}
}

func TestWithGlobalHeader(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	client, err := NewTransport("key", "issuer", privateKey, WithGlobalHeader("X-Custom", "value"))

	if err != nil {
		t.Fatalf("NewTransport with WithGlobalHeader failed: %v", err)
	}

	headerValue := client.httpClient.Header().Get("X-Custom")
	if headerValue != "value" {
		t.Errorf("Global header = %v, want 'value'", headerValue)
	}
}

func TestWithGlobalHeaders(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	headers := map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}

	client, err := NewTransport("key", "issuer", privateKey, WithGlobalHeaders(headers))

	if err != nil {
		t.Fatalf("NewTransport with WithGlobalHeaders failed: %v", err)
	}

	for key, expectedValue := range headers {
		gotValue := client.httpClient.Header().Get(key)
		if gotValue != expectedValue {
			t.Errorf("Header %q = %v, want %v", key, gotValue, expectedValue)
		}
	}
}

func TestWithProxy(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	proxyURL := "http://proxy.example.com:8080"
	client, err := NewTransport("key", "issuer", privateKey, WithProxy(proxyURL))

	if err != nil {
		t.Fatalf("NewTransport with WithProxy failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithProxy_Empty(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithProxy(""))

	if err == nil {
		t.Error("Expected error for empty proxy URL, got nil")
	}
}

func TestWithTLSClientConfig(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	client, err := NewTransport("key", "issuer", privateKey, WithTLSClientConfig(tlsConfig))

	if err != nil {
		t.Fatalf("NewTransport with WithTLSClientConfig failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithMinTLSVersion(t *testing.T) {
	tests := []struct {
		name       string
		minVersion uint16
	}{
		{"TLS 1.2", tls.VersionTLS12},
		{"TLS 1.3", tls.VersionTLS13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

			client, err := NewTransport("key", "issuer", privateKey, WithMinTLSVersion(tt.minVersion))

			if err != nil {
				t.Fatalf("NewTransport with WithMinTLSVersion failed: %v", err)
			}

			if client == nil {
				t.Error("Client is nil")
			}
		})
	}
}

func TestWithInsecureSkipVerify(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	client, err := NewTransport("key", "issuer", privateKey, WithInsecureSkipVerify())

	if err != nil {
		t.Fatalf("NewTransport with WithInsecureSkipVerify failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithClientCertificate(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	tmpDir := t.TempDir()

	// Generate test certificate and key
	cert, key := generateTestCert(t)
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certFile, cert, 0600); err != nil {
		t.Fatalf("Failed to write cert file: %v", err)
	}
	if err := os.WriteFile(keyFile, key, 0600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	client, err := NewTransport("key", "issuer", privateKey, WithClientCertificate(certFile, keyFile))

	if err != nil {
		t.Fatalf("NewTransport with WithClientCertificate failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithClientCertificateFromString(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	cert, key := generateTestCert(t)

	client, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithClientCertificateFromString(string(cert), string(key)),
	)

	if err != nil {
		t.Fatalf("NewTransport with WithClientCertificateFromString failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithRootCertificates(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	tmpDir := t.TempDir()
	caCert, _ := generateTestCert(t)
	caFile := filepath.Join(tmpDir, "ca.pem")

	if err := os.WriteFile(caFile, caCert, 0600); err != nil {
		t.Fatalf("Failed to write CA file: %v", err)
	}

	client, err := NewTransport("key", "issuer", privateKey, WithRootCertificates(caFile))

	if err != nil {
		t.Fatalf("NewTransport with WithRootCertificates failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithRootCertificateFromString(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	caCert, _ := generateTestCert(t)

	client, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithRootCertificateFromString(string(caCert)),
	)

	if err != nil {
		t.Fatalf("NewTransport with WithRootCertificateFromString failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestWithErrorHandler(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customHandler := NewErrorHandler(zap.NewNop())
	client, err := NewTransport("key", "issuer", privateKey, WithErrorHandler(customHandler))

	if err != nil {
		t.Fatalf("NewTransport with WithErrorHandler failed: %v", err)
	}

	if client.errorHandler != customHandler {
		t.Error("Error handler was not set correctly")
	}
}

func TestWithErrorHandler_Nil(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	_, err := NewTransport("key", "issuer", privateKey, WithErrorHandler(nil))

	if err == nil {
		t.Error("Expected error for nil error handler, got nil")
	}
}

func TestWithAudience(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customAudience := "custom-audience-v1"
	client, err := NewTransport("key", "issuer", privateKey, WithAudience(customAudience))

	if err != nil {
		t.Fatalf("NewTransport with WithAudience failed: %v", err)
	}

	// Verify audience was set on JWTAuth
	if jwtAuth, ok := client.auth.(*JWTAuth); ok {
		if jwtAuth.audience != customAudience {
			t.Errorf("JWT audience = %v, want %v", jwtAuth.audience, customAudience)
		}
	} else {
		t.Error("Auth is not JWTAuth type")
	}
}

func TestWithScope(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customScope := ScopeSchoolAPI
	client, err := NewTransport("key", "issuer", privateKey, WithScope(customScope))

	if err != nil {
		t.Fatalf("NewTransport with WithScope failed: %v", err)
	}

	// Verify scope was set on JWTAuth
	if jwtAuth, ok := client.auth.(*JWTAuth); ok {
		if jwtAuth.scope != customScope {
			t.Errorf("JWT scope = %v, want %v", jwtAuth.scope, customScope)
		}
	} else {
		t.Error("Auth is not JWTAuth type")
	}
}

func TestMultipleOptions(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	logger := zap.NewNop()
	customUA := "TestApp/1.0"

	client, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithLogger(logger),
		WithUserAgent(customUA),
		WithRetryCount(5),
		WithTimeout(60*time.Second),
		WithDebug(),
	)

	if err != nil {
		t.Fatalf("NewTransport with multiple options failed: %v", err)
	}

	if client.logger != logger {
		t.Error("Logger not set")
	}

	userAgent := client.httpClient.Header().Get("User-Agent")
	if userAgent != customUA {
		t.Errorf("User-Agent = %v, want %v", userAgent, customUA)
	}
}

func TestWithAPIVersion(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// WithAPIVersion is a no-op for Apple (version in URL path)
	client, err := NewTransport("key", "issuer", privateKey, WithAPIVersion("v2"))

	if err != nil {
		t.Fatalf("NewTransport with WithAPIVersion failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestOptionsAppliedInOrder(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// Apply same option multiple times - last one should win
	client, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithUserAgent("UA1"),
		WithUserAgent("UA2"),
		WithUserAgent("UA3"),
	)

	if err != nil {
		t.Fatalf("NewTransport failed: %v", err)
	}

	userAgent := client.httpClient.Header().Get("User-Agent")
	if userAgent != "UA3" {
		t.Errorf("User-Agent = %v, want 'UA3' (last applied)", userAgent)
	}
}

// Helper function to generate test certificate and key
func generateTestCert(t *testing.T) ([]byte, []byte) {
	// Generate a test key
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	// Marshal key
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// For testing, we'll use the same PEM as cert (not a real cert, but good enough for testing)
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: keyBytes, // Not a real cert, but sufficient for testing certificate loading
	})

	return certPEM, keyPEM
}

func TestWithTransport(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	customTransport := &http.Transport{
		MaxIdleConns: 100,
	}

	client, err := NewTransport("key", "issuer", privateKey, WithTransport(customTransport))

	if err != nil {
		t.Fatalf("NewTransport with WithTransport failed: %v", err)
	}

	if client == nil {
		t.Error("Client is nil")
	}
}

func TestAllOptionsDoNotError(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	logger := zap.NewNop()

	// Test that all options can be applied without error
	_, err := NewTransport(
		"key",
		"issuer",
		privateKey,
		WithLogger(logger),
		WithBaseURL("https://custom.api.example.com"),
		WithTimeout(30*time.Second),
		WithRetryCount(3),
		WithRetryWaitTime(1*time.Second),
		WithRetryMaxWaitTime(10*time.Second),
		WithUserAgent("TestApp/1.0"),
		WithDebug(),
		WithGlobalHeader("X-Test", "value"),
		WithMinTLSVersion(tls.VersionTLS12),
		WithAudience("custom-audience"),
		WithScope(ScopeBusinessAPI),
	)

	if err != nil {
		t.Fatalf("NewTransport with all options failed: %v", err)
	}
}
