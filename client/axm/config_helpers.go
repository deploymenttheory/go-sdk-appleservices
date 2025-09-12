package client

import (
	"fmt"
	"log"
	"os"
)

// LoadAndTestConfig loads configuration from a JSON file and tests it by making a simple API call
func LoadAndTestConfig(configPath string) (AXMConfig, *AXMClient, error) {
	var config AXMConfig

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil, fmt.Errorf("config file %s not found. Please create one based on config.example.json", configPath)
	}

	// Load configuration from file
	config, err := LoadConfigFromFile(configPath)
	if err != nil {
		return config, nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	// Create and test the AXM client
	axmClient, err := NewAXMClient(config)
	if err != nil {
		return config, nil, fmt.Errorf("failed to create AXM client: %w", err)
	}

	// Test authentication by forcing a token refresh
	if err := axmClient.ForceReauthenticate(); err != nil {
		axmClient.Close()
		return config, nil, fmt.Errorf("authentication test failed: %w", err)
	}

	return config, axmClient, nil
}

// LoadAndTestConfigWithEnvOverrides loads config from file with environment overrides and tests it
func LoadAndTestConfigWithEnvOverrides(configPath string) (AXMConfig, *AXMClient, error) {
	var config AXMConfig

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil, fmt.Errorf("config file %s not found. Please create one based on config.example.json", configPath)
	}

	// Load configuration with environment overrides
	config, err := LoadConfigFromFileWithEnvOverrides(configPath)
	if err != nil {
		return config, nil, fmt.Errorf("failed to load config with env overrides: %w", err)
	}

	// Create and test the AXM client
	axmClient, err := NewAXMClient(config)
	if err != nil {
		return config, nil, fmt.Errorf("failed to create AXM client: %w", err)
	}

	// Test authentication by forcing a token refresh
	if err := axmClient.ForceReauthenticate(); err != nil {
		axmClient.Close()
		return config, nil, fmt.Errorf("authentication test failed: %w", err)
	}

	return config, axmClient, nil
}

// ValidateConfig performs basic validation on an AXMConfig
func ValidateConfig(config AXMConfig) error {
	if config.ClientID == "" {
		return fmt.Errorf("clientID is required")
	}

	if config.KeyID == "" {
		return fmt.Errorf("keyID is required")
	}

	if config.PrivateKey == "" {
		return fmt.Errorf("privateKey is required")
	}

	if config.BaseURL == "" {
		return fmt.Errorf("baseURL is required")
	}

	// Test private key parsing - Apple AXM API requires ECDSA keys
	_, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return fmt.Errorf("invalid ECDSA private key: %w", err)
	}

	return nil
}

// CreateDefaultConfig creates an AXMConfig with sensible defaults
func CreateDefaultConfig(clientID, keyID, privateKey string) AXMConfig {
	return AXMConfig{
		BaseURL:    AppleSchoolManagerBaseURL,
		ClientID:   clientID,
		KeyID:      keyID,
		PrivateKey: privateKey,
		// Defaults will be set by NewAXMClient
	}
}

// LoadConfigFromEnv loads configuration from environment variables only
func LoadConfigFromEnv() (AXMConfig, error) {
	clientID := os.Getenv("APPLE_CLIENT_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	privateKey := os.Getenv("APPLE_PRIVATE_KEY")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	baseURL := os.Getenv("APPLE_BASE_URL")
	scope := os.Getenv("APPLE_SCOPE")

	// Validate required fields
	if clientID == "" {
		return AXMConfig{}, fmt.Errorf("APPLE_CLIENT_ID environment variable is required")
	}

	if keyID == "" {
		return AXMConfig{}, fmt.Errorf("APPLE_KEY_ID environment variable is required")
	}

	// Handle private key - either inline or from file path
	if privateKey == "" && privateKeyPath == "" {
		return AXMConfig{}, fmt.Errorf("either APPLE_PRIVATE_KEY or APPLE_PRIVATE_KEY_PATH environment variable is required")
	}

	if privateKey == "" && privateKeyPath != "" {
		var err error
		privateKey, err = LoadPrivateKeyFromFileWithValidation(privateKeyPath)
		if err != nil {
			return AXMConfig{}, fmt.Errorf("failed to load private key from env path %s: %w", privateKeyPath, err)
		}
	}

	config := CreateDefaultConfig(clientID, keyID, privateKey)

	// Override base URL if provided
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	// Set scope if provided
	if scope != "" {
		config.Scope = scope
	}

	return config, nil
}

// QuickStart provides a simple way to get started with minimal configuration
// It tries multiple configuration sources in order: env vars -> config file -> error
func QuickStart(configFilePath string) (*AXMClient, error) {
	// Try environment variables first
	if config, err := LoadConfigFromEnv(); err == nil {
		log.Printf("Loading configuration from environment variables")
		if err := ValidateConfig(config); err != nil {
			return nil, fmt.Errorf("environment config validation failed: %w", err)
		}

		client, err := NewAXMClient(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create client from env config: %w", err)
		}

		return client, nil
	}

	// Try config file if provided
	if configFilePath != "" {
		if _, err := os.Stat(configFilePath); err == nil {
			log.Printf("Loading configuration from file: %s", configFilePath)
			config, client, err := LoadAndTestConfigWithEnvOverrides(configFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from file %s: %w", configFilePath, err)
			}

			log.Printf("Successfully authenticated with Client ID: %s", config.ClientID)
			return client, nil
		}
	}

	// Try default config file location
	defaultConfigPath := "config.json"
	if _, err := os.Stat(defaultConfigPath); err == nil {
		log.Printf("Loading configuration from default file: %s", defaultConfigPath)
		config, client, err := LoadAndTestConfigWithEnvOverrides(defaultConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from default file %s: %w", defaultConfigPath, err)
		}

		log.Printf("Successfully authenticated with Client ID: %s", config.ClientID)
		return client, nil
	}

	return nil, fmt.Errorf("no valid configuration found. Please set environment variables or create a config.json file")
}

// PrintConfigSummary prints a summary of the loaded configuration (excluding sensitive data)
func PrintConfigSummary(config AXMConfig) {
	fmt.Printf("Configuration Summary:\n")
	fmt.Printf("  Base URL: %s\n", config.BaseURL)
	fmt.Printf("  Client ID: %s\n", config.ClientID)
	fmt.Printf("  Key ID: %s\n", config.KeyID)
	fmt.Printf("  Scope: %s\n", config.Scope)
	fmt.Printf("  Timeout: %v\n", config.Timeout)
	fmt.Printf("  Retry Count: %d\n", config.RetryCount)
	fmt.Printf("  Retry Delay: %v\n", config.RetryDelay)
	fmt.Printf("  User Agent: %s\n", config.UserAgent)
	fmt.Printf("  Debug: %t\n", config.Debug)
	fmt.Printf("  Private Key: %s\n", maskPrivateKey(config.PrivateKey))
}

// maskPrivateKey returns a masked version of the private key for logging
func maskPrivateKey(privateKey string) string {
	if privateKey == "" {
		return "Not set"
	}

	if len(privateKey) < 50 {
		return "***INVALID KEY***"
	}

	// Show first and last few characters
	return fmt.Sprintf("%s...%s (%d chars)",
		privateKey[:20],
		privateKey[len(privateKey)-20:],
		len(privateKey))
}
