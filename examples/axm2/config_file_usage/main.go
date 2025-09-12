package main

import (
	"context"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm2"
)

func main() {
	// Load configuration from JSON file
	configPath := "config.json" // Update path as needed
	config, err := axm2.LoadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create client (Resty v3)
	client, err := axm2.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}
	defer client.Close() // Resty v3 requires explicit close

	ctx := context.Background()

	log.Printf("Connected to %s API at %s", client.GetAPIType(), client.GetBaseURL())

	// Test authentication
	if !client.IsAuthenticated() {
		log.Println("Authenticating...")
		if err := client.ForceReauthenticate(); err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}
	}
	log.Println("Authentication successful!")

	// Get devices with comprehensive filtering
	queryBuilder := client.NewQueryBuilder().
		Limit(25).
		Fields("orgDevices", []string{ // Correct field name: "orgDevices" (plural)
			"serialNumber",
			"deviceModel",
			"productFamily",
			"productType",
			"status",
			"addedToOrgDateTime",
			"color",
			"deviceCapacity",
		})
		// Note: Apple AXM API doesn't support sort parameter

	devices, err := client.GetOrgDevices(ctx, queryBuilder)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	log.Printf("Retrieved %d devices", len(devices))

	// Analyze device statistics
	familyStats := make(map[string]int)
	statusStats := make(map[string]int)

	for _, device := range devices {
		if device.Attributes.ProductFamily != "" {
			familyStats[device.Attributes.ProductFamily]++
		}
		if device.Attributes.Status != "" {
			statusStats[device.Attributes.Status]++
		}
	}

	// Display statistics
	log.Println("\n=== Device Statistics ===")
	log.Println("By Product Family:")
	for family, count := range familyStats {
		log.Printf("  %s: %d devices", family, count)
	}

	log.Println("\nBy Status:")
	for status, count := range statusStats {
		log.Printf("  %s: %d devices", status, count)
	}

	// Show recent devices
	log.Println("\n=== Recent Devices ===")
	for i, device := range devices {
		log.Printf("%d. %s (%s) - %s - Added: %s",
			i+1,
			device.Attributes.SerialNumber,
			device.Attributes.DeviceModel,
			device.Attributes.Status,
			device.Attributes.AddedToOrgDateTime)

		if i >= 9 { // Show first 10
			break
		}
	}
}
