package main

import (
	"context"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm2"
)

func main() {
	// Configuration - in practice, load from JSON file or environment
	config := axm2.Config{
		APIType:  axm2.APITypeABM,                                    // or axm2.APITypeASM for Apple School Manager
		ClientID: "BUSINESSAPI.3bb3a62b-xxxx-xxxx-xxxx-a69b86201c5a", // Replace with your client ID
		KeyID:    "bb12ba87-147b-4e0d-9808-b4e6fbd5f9ba",             // Replace with your key ID
		PrivateKey: `-----BEGIN EC PRIVATE KEY-----
xxx
-----END EC PRIVATE KEY-----`, // Replace with your private key
		Debug: true, // Enable debug logging
	}

	// Create client
	client, err := axm2.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Test authentication
	if err := client.ForceReauthenticate(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	log.Printf("Successfully authenticated with client ID: %s", client.GetClientID())

	// Example: Get all organization devices with RequestOption pattern
	log.Println("\n=== Get Organization Devices ===")

	// Method 1: Using RequestOption pattern (recommended)
	devices, err := client.GetOrgDevices(ctx,
		axm2.WithLimitOption(10),
		axm2.WithFieldsOption("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"productFamily",
			"status",
			"addedToOrgDateTime",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	log.Printf("Retrieved %d devices", len(devices))

	// Display device details
	for i, device := range devices {
		log.Printf("Device %d:", i+1)
		log.Printf("  ID: %s", device.ID)
		log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
		log.Printf("  Model: %s", device.Attributes.DeviceModel)
		log.Printf("  Product Family: %s", device.Attributes.ProductFamily)
		log.Printf("  Status: %s", device.Attributes.Status)
		log.Printf("  Added: %s", device.Attributes.AddedToOrgDateTime)
		log.Println()
	}

	// Method 2: Using legacy QueryBuilder (for backward compatibility)
	log.Println("=== Using Legacy QueryBuilder ===")
	queryBuilder := client.NewQueryBuilder().
		Limit(5).
		Fields("orgDevices", []string{"serialNumber", "deviceModel", "status"})

	legacyDevices, err := client.GetOrgDevicesWithQuery(ctx, queryBuilder)
	if err != nil {
		log.Fatalf("Failed to get devices with QueryBuilder: %v", err)
	}

	log.Printf("Retrieved %d devices using legacy QueryBuilder", len(legacyDevices))
	for i, device := range legacyDevices {
		log.Printf("Legacy Device %d: %s (%s) - %s",
			i+1,
			device.Attributes.SerialNumber,
			device.Attributes.DeviceModel,
			device.Attributes.Status)
	}

	log.Println("\n=== Example completed successfully ===")
}
