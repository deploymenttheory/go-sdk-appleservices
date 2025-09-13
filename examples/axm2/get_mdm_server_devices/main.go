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

	// Example data - Replace with actual IDs from your organization
	exampleMdmServerID = "1F97349736CF4614A94F624E705841AD" // Replace with actual MDM server ID

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

	// Example 1: Get all devices assigned to MDM server
	log.Println("=== Get Devices Assigned to MDM Server ===")

	deviceIDs, err := client.MdmServers().GetDevices(ctx, exampleMdmServerID)
	if err != nil {
		log.Fatalf("Failed to get MDM server devices: %v", err)
	}

	log.Printf("Found %d devices assigned to MDM server '%s':", len(deviceIDs), exampleMdmServerID)

	if len(deviceIDs) == 0 {
		log.Printf("  No devices are currently assigned to this MDM server")
		log.Printf("  Use the assign_device_to_mdm_server example to assign devices")
	} else {
		for i, deviceID := range deviceIDs {
			if i >= 10 { // Limit output for readability
				log.Printf("  ... and %d more devices", len(deviceIDs)-10)
				break
			}
			log.Printf("  Device %d: %s", i+1, deviceID)
		}
	}

	// Example 2: Get devices with pagination limit
	log.Println("\n=== Get MDM Server Devices with Pagination Limit ===")

	limitedDeviceIDs, err := client.MdmServers().GetDevices(ctx, exampleMdmServerID,
		axm2.WithLimitOption(5), // Limit to 5 devices per page
	)
	if err != nil {
		log.Fatalf("Failed to get limited MDM server devices: %v", err)
	}

	log.Printf("Retrieved %d device IDs with pagination limit", len(limitedDeviceIDs))

	// Note: The SDK automatically handles pagination, so you still get all device IDs,
	// but the API requests are made in smaller chunks of 5 devices per request

	// Example 3: Cross-reference with device details
	if len(deviceIDs) > 0 {
		log.Println("\n=== Cross-Reference with Device Details ===")
		log.Printf("Getting detailed information for first assigned device...")

		firstDeviceID := deviceIDs[0]
		device, err := client.OrgDevices().GetOrgDevice(ctx, firstDeviceID,
			axm2.WithFieldsOption("orgDevices", []string{
				"serialNumber",
				"deviceModel",
				"status",
				"productFamily",
			}),
		)
		if err != nil {
			log.Printf("Failed to get device details for %s: %v", firstDeviceID, err)
		} else {
			log.Printf("Device Details:")
			log.Printf("  ID: %s", device.ID)
			log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
			log.Printf("  Model: %s", device.Attributes.DeviceModel)
			log.Printf("  Status: %s", device.Attributes.Status)
			log.Printf("  Product Family: %s", device.Attributes.ProductFamily)
		}
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("This endpoint returns device serial numbers (IDs) assigned to the MDM server.")
	log.Printf("To get full device details, use the device IDs with:")
	log.Printf("  - client.OrgDevices().GetOrgDevice(ctx, deviceID) for individual devices")
	log.Printf("  - client.OrgDevices().GetOrgDevices(ctx) to get all devices and filter by assignment")
	log.Printf("")
	log.Printf("Device assignment status can be:")
	log.Printf("  - ASSIGNED: Device is assigned to an MDM server")
	log.Printf("  - UNASSIGNED: Device is not assigned to any MDM server")
}
