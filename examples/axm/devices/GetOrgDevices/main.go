package main

import (
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm"
	"go.uber.org/zap"
)

func main() {
	// Load configuration from JSON file
	configPath := "/Users/dafyddwatkins/localtesting/axm/config.example.json"
	config, err := client.LoadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create AXM client (this initializes the logger based on config.Debug)
	axmClient, err := client.NewAXMClient(config)
	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}
	defer axmClient.Close()

	// Use the client's logger
	logger := axmClient.Logger

	logger.Info("=== Basic GetOrgDevices Example ===")
	logger.Info("This example demonstrates loading config from JSON and calling GetOrgDevices")

	// Test authentication
	if err := axmClient.ForceReauthenticate(); err != nil {
		logger.Fatal("Authentication test failed", zap.Error(err))
	}
	logger.Info("Successfully authenticated", zap.String("client_id", config.ClientID))

	// Create AXM service client
	axmService := axm.NewClient(axmClient)
	defer axmService.Close()

	// Build query with optional filtering and field selection
	queryBuilder := axmService.NewQueryBuilder().
		Limit(50).                    // Limit results to 50 devices
		Fields("orgDevice", []string{ // Select specific fields to reduce response size
			"serialNumber",
			"deviceModel",
			"productFamily",
			"status",
			"addedToOrgDateTime",
		}).
		Filter("status", "ASSIGNED"). // Filter for assigned devices only
		Sort("addedToOrgDateTime")    // Sort by when devices were added

	logger.Info("Fetching organization devices with filters",
		zap.String("client_id", config.ClientID),
		zap.String("base_url", config.BaseURL))

	// Get organization devices
	devices, err := axmService.GetOrgDevices(queryBuilder)
	if err != nil {
		logger.Fatal("Failed to get organization devices", zap.Error(err))
	}

	logger.Info("Successfully retrieved organization devices",
		zap.Int("device_count", len(devices)))

	if len(devices) == 0 {
		logger.Info("No devices found matching the criteria")
		return
	}

	// Count devices by product family and status
	familyCount := make(map[string]int)
	statusCount := make(map[string]int)

	for _, device := range devices {
		if device.Attributes.ProductFamily != "" {
			familyCount[device.Attributes.ProductFamily]++
		}
		if device.Attributes.Status != "" {
			statusCount[device.Attributes.Status]++
		}
	}

	// Log summary statistics
	logger.Info("Device summary by product family")
	for family, count := range familyCount {
		logger.Info("Product family count",
			zap.String("family", family),
			zap.Int("count", count))
	}

	logger.Info("Device summary by status")
	for status, count := range statusCount {
		logger.Info("Status count",
			zap.String("status", status),
			zap.Int("count", count))
	}

	// Process the first few devices
	logger.Info("=== Device Details ===")
	for i, device := range devices {
		logger.Info("Device details",
			zap.Int("index", i+1),
			zap.String("id", device.ID),
			zap.String("serial_number", device.Attributes.SerialNumber),
			zap.String("device_model", device.Attributes.DeviceModel),
			zap.String("product_family", device.Attributes.ProductFamily),
			zap.String("status", device.Attributes.Status),
			zap.String("added_to_org", device.Attributes.AddedToOrgDateTime),
		)

		// Stop after showing first 10 devices for demo purposes
		if i >= 9 {
			logger.Info("Showing first 10 devices only for demo",
				zap.Int("total_devices", len(devices)))
			break
		}
	}

	logger.Info("Example completed successfully",
		zap.Int("total_devices_processed", len(devices)))
}
