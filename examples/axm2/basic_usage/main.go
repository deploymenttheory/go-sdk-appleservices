package main

import (
	"context"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm2"
)

func main() {
	// Example configuration - in practice, load from JSON file or environment
	config := axm2.Config{
		APIType:  axm2.APITypeABM,                                    // or axm2.APITypeASM for Apple School Manager
		ClientID: "BUSINESSAPI.00000000-0000-0000-0000-000000000000", // Replace with your client ID
		KeyID:    "00000000-0000-0000-0000-000000000000",             // Replace with your key ID
		PrivateKey: `-----BEGIN EC PRIVATE KEY-----
xxxxxx
-----END EC PRIVATE KEY-----`, // Replace with your private key
		Debug: true, // Enable debug logging
	}

	// Create client using direct pattern (Resty v3)
	client, err := axm2.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}
	defer client.Close() // Resty v3 requires explicit close

	ctx := context.Background()

	// Test authentication
	if err := client.ForceReauthenticate(); err != nil {
		log.Fatalf("Authentication test failed: %v", err)
	}
	log.Printf("Successfully authenticated with client ID: %s", client.GetClientID())

	// Example 1: Get all devices with basic query
	log.Println("\n=== Example 1: Get All Devices ===")
	devices, err := client.GetOrgDevices(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}
	log.Printf("Retrieved %d devices", len(devices))

	// Example 2: Get devices with field selection (Apple AXM API specific)
	log.Println("\n=== Example 2: Device Query with Field Selection ===")
	queryBuilder := client.NewQueryBuilder().
		Limit(10).                     // Limit to 10 devices
		Fields("orgDevices", []string{ // Correct field name: "orgDevices" (plural)
			"serialNumber",
			"deviceModel",
			"productFamily",
			"status",
			"addedToOrgDateTime",
		})
		// Note: Apple AXM API doesn't support filter[status] or sort parameters
		// These would cause 400 Bad Request errors

	filteredDevices, err := client.GetOrgDevices(ctx, queryBuilder)
	if err != nil {
		log.Fatalf("Failed to get filtered devices: %v", err)
	}

	log.Printf("Retrieved %d filtered devices", len(filteredDevices))

	// Display device details
	for i, device := range filteredDevices {
		log.Printf("Device %d: %s (%s) - %s - Status: %s",
			i+1,
			device.Attributes.SerialNumber,
			device.Attributes.DeviceModel,
			device.Attributes.ProductFamily,
			device.Attributes.Status)

		if i >= 4 { // Show first 5 devices
			break
		}
	}

	// Example 3: Get a specific device by ID
	if len(devices) > 0 {
		log.Println("\n=== Example 3: Get Specific Device ===")
		deviceID := devices[0].ID

		specificDevice, err := client.GetOrgDevice(ctx, deviceID, nil)
		if err != nil {
			log.Printf("Failed to get specific device: %v", err)
		} else {
			log.Printf("Retrieved device: %s (%s)",
				specificDevice.Attributes.SerialNumber,
				specificDevice.Attributes.DeviceModel)
		}
	}

	log.Println("\n=== Examples completed successfully ===")
}
