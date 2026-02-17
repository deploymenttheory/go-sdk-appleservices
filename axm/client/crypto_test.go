package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePrivateKey_ECDSA(t *testing.T) {
	// Generate ECDSA key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	// Marshal to PKCS8 format
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Parse key
	parsedKey, err := ParsePrivateKey(pemData)
	if err != nil {
		t.Fatalf("ParsePrivateKey failed: %v", err)
	}

	// Verify it's an ECDSA key
	if _, ok := parsedKey.(*ecdsa.PrivateKey); !ok {
		t.Errorf("Expected *ecdsa.PrivateKey, got %T", parsedKey)
	}
}

func TestParsePrivateKey_RSA(t *testing.T) {
	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Marshal to PKCS8 format
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Parse key
	parsedKey, err := ParsePrivateKey(pemData)
	if err != nil {
		t.Fatalf("ParsePrivateKey failed: %v", err)
	}

	// Verify it's an RSA key
	if _, ok := parsedKey.(*rsa.PrivateKey); !ok {
		t.Errorf("Expected *rsa.PrivateKey, got %T", parsedKey)
	}
}

func TestParsePrivateKey_PKCS1_RSA(t *testing.T) {
	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Marshal to PKCS1 format
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Parse key
	parsedKey, err := ParsePrivateKey(pemData)
	if err != nil {
		t.Fatalf("ParsePrivateKey failed: %v", err)
	}

	// Verify it's an RSA key
	if _, ok := parsedKey.(*rsa.PrivateKey); !ok {
		t.Errorf("Expected *rsa.PrivateKey, got %T", parsedKey)
	}
}

func TestParsePrivateKey_EC_Format(t *testing.T) {
	// Generate ECDSA key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	// Marshal to EC format
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Parse key
	parsedKey, err := ParsePrivateKey(pemData)
	if err != nil {
		t.Fatalf("ParsePrivateKey failed: %v", err)
	}

	// Verify it's an ECDSA key
	if _, ok := parsedKey.(*ecdsa.PrivateKey); !ok {
		t.Errorf("Expected *ecdsa.PrivateKey, got %T", parsedKey)
	}
}

func TestParsePrivateKey_InvalidPEM(t *testing.T) {
	invalidPEM := []byte("not a valid PEM")

	_, err := ParsePrivateKey(invalidPEM)
	if err == nil {
		t.Error("Expected error for invalid PEM, got nil")
	}
}

func TestParsePrivateKey_InvalidKey(t *testing.T) {
	// Create PEM with invalid key data
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: []byte("invalid key data"),
	})

	_, err := ParsePrivateKey(pemData)
	if err == nil {
		t.Error("Expected error for invalid key data, got nil")
	}
}

func TestLoadPrivateKeyFromFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Generate ECDSA key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	// Marshal to PKCS8 format
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Write to temp file
	keyPath := filepath.Join(tmpDir, "test_key.p8")
	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	// Load key from file
	loadedKey, err := LoadPrivateKeyFromFile(keyPath)
	if err != nil {
		t.Fatalf("LoadPrivateKeyFromFile failed: %v", err)
	}

	// Verify it's an ECDSA key
	if _, ok := loadedKey.(*ecdsa.PrivateKey); !ok {
		t.Errorf("Expected *ecdsa.PrivateKey, got %T", loadedKey)
	}
}

func TestLoadPrivateKeyFromFile_FileNotFound(t *testing.T) {
	_, err := LoadPrivateKeyFromFile("/nonexistent/key.p8")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestValidatePrivateKey_ECDSA(t *testing.T) {
	// Generate ECDSA key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	err = ValidatePrivateKey(privateKey)
	if err != nil {
		t.Errorf("ValidatePrivateKey failed for valid ECDSA key: %v", err)
	}
}

func TestValidatePrivateKey_RSA_Valid(t *testing.T) {
	// Generate 2048-bit RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	err = ValidatePrivateKey(privateKey)
	if err != nil {
		t.Errorf("ValidatePrivateKey failed for valid RSA key: %v", err)
	}
}

func TestValidatePrivateKey_RSA_TooSmall(t *testing.T) {
	// Generate 1024-bit RSA key (too small)
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	err = ValidatePrivateKey(privateKey)
	if err == nil {
		t.Error("Expected error for RSA key too small, got nil")
	}
}

func TestValidatePrivateKey_Nil(t *testing.T) {
	err := ValidatePrivateKey(nil)
	if err == nil {
		t.Error("Expected error for nil key, got nil")
	}
}

func TestValidatePrivateKey_UnsupportedType(t *testing.T) {
	err := ValidatePrivateKey("not a key")
	if err == nil {
		t.Error("Expected error for unsupported key type, got nil")
	}
}

func TestLoadPrivateKeyFromEnv_NotSet(t *testing.T) {
	// Ensure env var is not set
	oldValue := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	os.Unsetenv("APPLE_PRIVATE_KEY_PATH")
	defer func() {
		if oldValue != "" {
			os.Setenv("APPLE_PRIVATE_KEY_PATH", oldValue)
		}
	}()

	_, err := LoadPrivateKeyFromEnv()
	if err == nil {
		t.Error("Expected error when APPLE_PRIVATE_KEY_PATH not set, got nil")
	}
}

func TestLoadPrivateKeyFromEnv_Success(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Generate ECDSA key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	// Marshal to PKCS8 format
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal key: %v", err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	// Write to temp file
	keyPath := filepath.Join(tmpDir, "test_key.p8")
	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	// Set env var
	oldValue := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	os.Setenv("APPLE_PRIVATE_KEY_PATH", keyPath)
	defer func() {
		if oldValue != "" {
			os.Setenv("APPLE_PRIVATE_KEY_PATH", oldValue)
		} else {
			os.Unsetenv("APPLE_PRIVATE_KEY_PATH")
		}
	}()

	// Load key from env
	loadedKey, err := LoadPrivateKeyFromEnv()
	if err != nil {
		t.Fatalf("LoadPrivateKeyFromEnv failed: %v", err)
	}

	// Verify it's an ECDSA key
	if _, ok := loadedKey.(*ecdsa.PrivateKey); !ok {
		t.Errorf("Expected *ecdsa.PrivateKey, got %T", loadedKey)
	}
}
