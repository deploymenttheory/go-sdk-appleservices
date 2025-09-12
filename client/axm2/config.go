package axm2

import (
	"encoding/json"
	"fmt"
	"os"
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

// validateConfig checks that all required fields are present
func validateConfig(config Config) error {
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
