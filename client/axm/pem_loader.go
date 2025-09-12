package client

import (
	"fmt"
	"os"
	"strings"
)

// LoadPrivateKeyFromFile reads a PEM-encoded private key from a file path
func LoadPrivateKeyFromFile(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}

	// Read the PEM file
	pemData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read PEM file %s: %w", filePath, err)
	}

	// Convert bytes to string and return
	return string(pemData), nil
}

// LoadPrivateKeyFromFileWithValidation reads and validates a PEM-encoded private key from a file path
func LoadPrivateKeyFromFileWithValidation(filePath string) (string, error) {
	pemContent, err := LoadPrivateKeyFromFile(filePath)
	if err != nil {
		return "", err
	}

	// Basic validation - check if it contains PEM markers
	if !containsPEMMarkers(pemContent) {
		return "", fmt.Errorf("file %s does not appear to contain a valid PEM-encoded key", filePath)
	}

	return pemContent, nil
}

// containsPEMMarkers checks if the content contains basic PEM markers
func containsPEMMarkers(content string) bool {
	hasBegin := false
	hasEnd := false

	// Check for common PEM begin markers
	beginMarkers := []string{
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----",
	}

	// Check for common PEM end markers
	endMarkers := []string{
		"-----END RSA PRIVATE KEY-----",
		"-----END PRIVATE KEY-----",
		"-----END EC PRIVATE KEY-----",
	}

	for _, marker := range beginMarkers {
		if strings.Contains(content, marker) {
			hasBegin = true
			break
		}
	}

	for _, marker := range endMarkers {
		if strings.Contains(content, marker) {
			hasEnd = true
			break
		}
	}

	return hasBegin && hasEnd
}

