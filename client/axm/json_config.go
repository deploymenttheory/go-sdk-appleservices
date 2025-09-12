package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AXMConfigFile represents the JSON configuration file structure
type AXMConfigFile struct {
	BaseURL        string `json:"baseUrl,omitempty"`
	ClientID       string `json:"clientId"`
	KeyID          string `json:"keyId"`
	PrivateKey     string `json:"privateKey,omitempty"`
	PrivateKeyPath string `json:"privateKeyPath,omitempty"`
	Scope          string `json:"scope,omitempty"`
	TimeoutSeconds int    `json:"timeoutSeconds,omitempty"`
	RetryCount     int    `json:"retryCount,omitempty"`
	RetryDelayMs   int    `json:"retryDelayMs,omitempty"`
	UserAgent      string `json:"userAgent,omitempty"`
	Debug          bool   `json:"debug,omitempty"`
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

	// Convert to AXMConfig (defaults will be set by NewAXMClient)
	config = AXMConfig{
		BaseURL:  configFile.BaseURL,
		ClientID: configFile.ClientID,
		KeyID:    configFile.KeyID,
		Scope:    configFile.Scope,
		Debug:    configFile.Debug,
	}

	// Set optional timeout values if provided
	if configFile.TimeoutSeconds > 0 {
		config.Timeout = time.Duration(configFile.TimeoutSeconds) * time.Second
	}
	if configFile.RetryCount > 0 {
		config.RetryCount = configFile.RetryCount
	}
	if configFile.RetryDelayMs > 0 {
		config.RetryDelay = time.Duration(configFile.RetryDelayMs) * time.Millisecond
	}
	if configFile.UserAgent != "" {
		config.UserAgent = configFile.UserAgent
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

	return config, nil
}
