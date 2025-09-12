package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AXMConfigFile represents the JSON configuration file structure
type AXMConfigFile struct {
	BaseURL         string `json:"baseUrl,omitempty"`
	OrgID           string `json:"orgId"`
	KeyID           string `json:"keyId"`
	PrivateKey      string `json:"privateKey,omitempty"`
	PrivateKeyPath  string `json:"privateKeyPath,omitempty"`
	TimeoutSeconds  int    `json:"timeoutSeconds,omitempty"`
	RetryCount      int    `json:"retryCount,omitempty"`
	RetryDelayMs    int    `json:"retryDelayMs,omitempty"`
	UserAgent       string `json:"userAgent,omitempty"`
	Debug           bool   `json:"debug,omitempty"`
}

// LoadConfigFromFile loads AXM configuration from a JSON file
func LoadConfigFromFile(filePath string) (AXMConfig, error) {
	var config AXMConfig
	
	if filePath == "" {
		return config, fmt.Errorf("config file path cannot be empty")
	}

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	// Parse JSON into config file structure
	var configFile AXMConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return config, fmt.Errorf("failed to parse JSON config file %s: %w", filePath, err)
	}

	// Convert to AXMConfig with defaults
	config = AXMConfig{
		OrgID:  configFile.OrgID,
		KeyID:  configFile.KeyID,
		Debug:  configFile.Debug,
	}

	// Set BaseURL with default
	if configFile.BaseURL != "" {
		config.BaseURL = configFile.BaseURL
	} else {
		config.BaseURL = AppleSchoolManagerBaseURL
	}

	// Handle private key - either inline or from file path
	if configFile.PrivateKey != "" {
		config.PrivateKey = configFile.PrivateKey
	} else if configFile.PrivateKeyPath != "" {
		privateKey, err := LoadPrivateKeyFromFileWithValidation(configFile.PrivateKeyPath)
		if err != nil {
			return config, fmt.Errorf("failed to load private key from path %s: %w", configFile.PrivateKeyPath, err)
		}
		config.PrivateKey = privateKey
	} else {
		return config, fmt.Errorf("either privateKey or privateKeyPath must be specified in config file")
	}

	// Set timeout with default
	if configFile.TimeoutSeconds > 0 {
		config.Timeout = time.Duration(configFile.TimeoutSeconds) * time.Second
	} else {
		config.Timeout = 30 * time.Second
	}

	// Set retry count with default
	if configFile.RetryCount > 0 {
		config.RetryCount = configFile.RetryCount
	} else {
		config.RetryCount = 3
	}

	// Set retry delay with default
	if configFile.RetryDelayMs > 0 {
		config.RetryDelay = time.Duration(configFile.RetryDelayMs) * time.Millisecond
	} else {
		config.RetryDelay = 1 * time.Second
	}

	// Set user agent with default
	if configFile.UserAgent != "" {
		config.UserAgent = configFile.UserAgent
	} else {
		config.UserAgent = "go-api-sdk-apple/1.0.0"
	}

	return config, nil
}

// LoadConfigFromFileWithEnvOverrides loads config from JSON file and allows environment variable overrides
func LoadConfigFromFileWithEnvOverrides(filePath string) (AXMConfig, error) {
	config, err := LoadConfigFromFile(filePath)
	if err != nil {
		return config, err
	}

	// Allow environment variables to override config file values
	if orgID := os.Getenv("APPLE_ORG_ID"); orgID != "" {
		config.OrgID = orgID
	}

	if keyID := os.Getenv("APPLE_KEY_ID"); keyID != "" {
		config.KeyID = keyID
	}

	if privateKey := os.Getenv("APPLE_PRIVATE_KEY"); privateKey != "" {
		config.PrivateKey = privateKey
	}

	if privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH"); privateKeyPath != "" {
		privateKey, err := LoadPrivateKeyFromFileWithValidation(privateKeyPath)
		if err != nil {
			return config, fmt.Errorf("failed to load private key from env path %s: %w", privateKeyPath, err)
		}
		config.PrivateKey = privateKey
	}

	if baseURL := os.Getenv("APPLE_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	// Validate required fields after env overrides
	if config.OrgID == "" {
		return config, fmt.Errorf("orgId is required (set in config file or APPLE_ORG_ID env var)")
	}

	if config.KeyID == "" {
		return config, fmt.Errorf("keyId is required (set in config file or APPLE_KEY_ID env var)")
	}

	if config.PrivateKey == "" {
		return config, fmt.Errorf("private key is required (set privateKey/privateKeyPath in config file or APPLE_PRIVATE_KEY/APPLE_PRIVATE_KEY_PATH env vars)")
	}

	return config, nil
}

// SaveConfigToFile saves an AXMConfig to a JSON file (excludes sensitive private key data)
func SaveConfigToFile(config AXMConfig, filePath string, includePrivateKeyPath string) error {
	if filePath == "" {
		return fmt.Errorf("config file path cannot be empty")
	}

	configFile := AXMConfigFile{
		BaseURL:        config.BaseURL,
		OrgID:          config.OrgID,
		KeyID:          config.KeyID,
		TimeoutSeconds: int(config.Timeout.Seconds()),
		RetryCount:     config.RetryCount,
		RetryDelayMs:   int(config.RetryDelay.Milliseconds()),
		UserAgent:      config.UserAgent,
		Debug:          config.Debug,
	}

	// Include private key path if provided (safer than including the actual key)
	if includePrivateKeyPath != "" {
		configFile.PrivateKeyPath = includePrivateKeyPath
	}

	data, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, 0600) // Restrictive permissions for config files
	if err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filePath, err)
	}

	return nil
}