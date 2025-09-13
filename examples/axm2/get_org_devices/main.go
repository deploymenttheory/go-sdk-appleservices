package main

import (
	"context"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm2"
)

// Configuration constants - externalized for easier maintenance
const (
	// API Configuration
	apiType  = axm2.APITypeABM                                    // or axm2.APITypeASM for Apple School Manager
	clientID = "BUSINESSAPI.3bb3a62b-xxxx-xxxx-xxxx-a69b86201c5a" // Replace with your client ID
	keyID    = "bb12ba87-xxxx-xxxx-xxxx-b4e6fbd5f9ba"             // Replace with your key ID

	// Private key - Replace with your private key
	privateKey = `-----BEGIN EC PRIVATE KEY-----
xxx
-----END EC PRIVATE KEY-----`

	// Debug settings
	enableDebug = true
)

func main() {
	// Configuration - in practice, load from JSON file or environment
	config := axm2.Config{
		APIType:    apiType,
		ClientID:   clientID,
		KeyID:      keyID,
		PrivateKey: privateKey,
		Debug:      enableDebug,
	}

	// Create client
	var client axm2.AXMClient
	var err error
	client, err = axm2.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Get all organization devices with default settings
	log.Println("=== Get All Organization Devices ===")

	devices, err := client.OrgDevices().GetOrgDevices(ctx)
	if err != nil {
		log.Fatalf("Failed to get organization devices: %v", err)
	}

	log.Printf("Found %d devices in organization:", len(devices))
	for i, device := range devices {
		if i >= 5 { // Limit output for readability
			log.Printf("  ... and %d more devices", len(devices)-5)
			break
		}
		log.Printf("  Device %d:", i+1)
		log.Printf("    ID: %s", device.ID)
		log.Printf("    Serial Number: %s", device.Attributes.SerialNumber)
		log.Printf("    Model: %s", device.Attributes.DeviceModel)
		log.Printf("    Product Family: %s", device.Attributes.ProductFamily)
		log.Printf("    Status: %s", device.Attributes.Status)
		if device.Attributes.AddedToOrgDateTime != "" {
			log.Printf("    Added to Org: %s", device.Attributes.AddedToOrgDateTime)
		}
	}

	// Example 2: Get devices with field filtering
	log.Println("\n=== Get Devices with Field Filtering ===")

	filteredDevices, err := client.OrgDevices().GetOrgDevices(ctx,
		axm2.WithFieldsOption("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"status",
			"productFamily",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered devices: %v", err)
	}

	log.Printf("Filtered devices (first 3):")
	for i, device := range filteredDevices {
		if i >= 3 {
			break
		}
		log.Printf("  Device %d:", i+1)
		log.Printf("    Serial: %s", device.Attributes.SerialNumber)
		log.Printf("    Model: %s", device.Attributes.DeviceModel)
		log.Printf("    Status: %s", device.Attributes.Status)
		log.Printf("    Family: %s", device.Attributes.ProductFamily)
	}

	// Example 3: Get devices with pagination limit
	log.Println("\n=== Get Devices with Pagination Limit ===")

	limitedDevices, err := client.OrgDevices().GetOrgDevices(ctx,
		axm2.WithLimitOption(10), // Limit to 10 devices per page
	)
	if err != nil {
		log.Fatalf("Failed to get limited devices: %v", err)
	}

	log.Printf("Retrieved %d devices with pagination limit", len(limitedDevices))

	// Note: The SDK automatically handles pagination, so you still get all devices,
	// but the API requests are made in smaller chunks of 10 devices per request
}
