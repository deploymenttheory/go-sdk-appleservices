package client

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKeyFromFile loads a private key (RSA or ECDSA) from a PEM file
func LoadPrivateKeyFromFile(filePath string) (any, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	return ParsePrivateKey(keyData)
}

// ParsePrivateKey parses a private key (RSA or ECDSA) from PEM-encoded data
func ParsePrivateKey(keyData []byte) (any, error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try different parsing methods
	var key any
	var err error

	// Try PKCS8 first (most common for .p8 files)
	key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// Try EC private key format
			key, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key (tried PKCS8, PKCS1, and EC formats): %w", err)
			}
		}
	}

	// Check if it's an RSA or ECDSA key
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return k, nil
	case *ecdsa.PrivateKey:
		return k, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %T (expected RSA or ECDSA)", key)
	}
}

// LoadPrivateKeyFromEnv loads a private key from the environment variable APPLE_PRIVATE_KEY_PATH
func LoadPrivateKeyFromEnv() (any, error) {
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		return nil, fmt.Errorf("APPLE_PRIVATE_KEY_PATH environment variable is not set")
	}

	return LoadPrivateKeyFromFile(privateKeyPath)
}

// ValidatePrivateKey validates that the private key is suitable for JWT signing
func ValidatePrivateKey(privateKey any) error {
	if privateKey == nil {
		return fmt.Errorf("private key is nil")
	}

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		// Check key size (Apple requires at least 2048 bits for RSA)
		keySize := key.Size() * 8
		if keySize < 2048 {
			return fmt.Errorf("RSA private key size (%d bits) is too small, minimum 2048 bits required", keySize)
		}
	case *ecdsa.PrivateKey:
		// ECDSA keys are generally acceptable for Apple APIs
		// P-256 curve is commonly used and supported
		return nil
	default:
		return fmt.Errorf("unsupported private key type: %T", privateKey)
	}

	return nil
}
