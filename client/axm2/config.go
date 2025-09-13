package axm2

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(configPath string) (Config, error) {
	var config Config

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate required fields
	if err := validateConfig(config); err != nil {
		return config, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validateConfig checks that all required fields are present with enhanced validation
func validateConfig(config Config) error {
	validators := []func(Config) error{
		validateRequiredFields,
		validateTimeouts,
		validateRetrySettings,
		validateAPIEndpoints,
	}

	for _, validator := range validators {
		if err := validator(config); err != nil {
			return err
		}
	}
	return nil
}

// validateRequiredFields checks that all required fields are present
func validateRequiredFields(config Config) error {
	if config.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if config.KeyID == "" {
		return fmt.Errorf("key_id is required")
	}
	if config.PrivateKey == "" {
		return fmt.Errorf("private_key is required")
	}
	if config.APIType == "" {
		return fmt.Errorf("api_type is required (must be 'abm' or 'asm')")
	}
	if config.APIType != APITypeABM && config.APIType != APITypeASM {
		return fmt.Errorf("api_type must be 'abm' or 'asm', got: %s", config.APIType)
	}

	return nil
}

// validateTimeouts checks timeout configuration
func validateTimeouts(config Config) error {
	if config.Timeout > 0 && config.Timeout < time.Second {
		return fmt.Errorf("timeout must be at least 1 second")
	}
	return nil
}

// validateRetrySettings checks retry configuration
func validateRetrySettings(config Config) error {
	if config.RetryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}
	if config.RetryMaxWait > 0 && config.RetryMinWait > 0 && config.RetryMaxWait < config.RetryMinWait {
		return fmt.Errorf("retry max wait must be >= min wait")
	}
	return nil
}

// validateAPIEndpoints checks API endpoint configuration
func validateAPIEndpoints(config Config) error {
	// Validate custom base URL if provided
	if config.BaseURL != "" {
		if !strings.HasPrefix(config.BaseURL, "https://") {
			return fmt.Errorf("base URL must use HTTPS")
		}
	}
	return nil
}
