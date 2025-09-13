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

	// First, get a list of MDM servers to find a server ID
	log.Println("\n=== Finding an MDM Server ===")
	servers, err := client.GetMdmServers(ctx)
	if err != nil {
		log.Fatalf("Failed to get MDM servers: %v", err)
	}

	if len(servers) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	// Use the first server (preferably the default one)
	var targetServer *axm2.MdmServer
	for i := range servers {
		if servers[i].Attributes.DefaultMdmServer {
			targetServer = &servers[i]
			break
		}
	}
	if targetServer == nil {
		targetServer = &servers[0] // Use first server if no default found
	}

	log.Printf("Using MDM server: %s (ID: %s)", targetServer.Attributes.Name, targetServer.ID)

	// Example: Get device serial numbers assigned to an MDM server
	log.Println("\n=== Get MDM Server Devices ===")

	// Method 1: Get all devices assigned to the server
	deviceIDs, err := client.GetMdmServerDevices(ctx, targetServer.ID)
	if err != nil {
		log.Fatalf("Failed to get MDM server devices: %v", err)
	}

	log.Printf("Retrieved %d device IDs assigned to MDM server '%s'", len(deviceIDs), targetServer.Attributes.Name)

	if len(deviceIDs) == 0 {
		log.Println("No devices are currently assigned to this MDM server")
		log.Println("This could mean:")
		log.Println("  - The server is newly created and has no devices yet")
		log.Println("  - All devices are managed by a different MDM server")
		log.Println("  - Devices haven't been assigned to any MDM server yet")
	} else {
		// Display the first few device IDs
		log.Printf("First %d device IDs:", min(len(deviceIDs), 10))
		for i, deviceID := range deviceIDs {
			if i >= 10 {
				log.Printf("  ... and %d more", len(deviceIDs)-10)
				break
			}
			log.Printf("  %d. %s", i+1, deviceID)
		}

		// Get detailed information for the first few devices
		log.Println("\n=== Device Details for First Few Devices ===")
		for i, deviceID := range deviceIDs {
			if i >= 3 { // Only show first 3 devices in detail
				break
			}

			device, err := client.GetOrgDevice(ctx, deviceID,
				axm2.WithFieldsOption("orgDevices", []string{
					"serialNumber",
					"deviceModel",
					"productFamily",
					"status",
				}),
			)
			if err != nil {
				log.Printf("  Failed to get details for device %s: %v", deviceID, err)
				continue
			}

			log.Printf("  Device %d:", i+1)
			log.Printf("    ID: %s", device.ID)
			log.Printf("    Serial Number: %s", device.Attributes.SerialNumber)
			log.Printf("    Model: %s", device.Attributes.DeviceModel)
			log.Printf("    Product Family: %s", device.Attributes.ProductFamily)
			log.Printf("    Status: %s", device.Attributes.Status)
			log.Println()
		}
	}

	// Method 2: Get devices with pagination limit
	log.Println("=== Get MDM Server Devices with Limit ===")
	limitedDeviceIDs, err := client.GetMdmServerDevices(ctx, targetServer.ID,
		axm2.WithLimitOption(5),
	)
	if err != nil {
		log.Fatalf("Failed to get limited MDM server devices: %v", err)
	}

	log.Printf("Retrieved %d device IDs (limited to 5) for MDM server '%s'",
		len(limitedDeviceIDs), targetServer.Attributes.Name)

	for i, deviceID := range limitedDeviceIDs {
		log.Printf("  Limited Device %d: %s", i+1, deviceID)
	}

	// Summary information
	log.Println("\n=== Summary ===")
	log.Printf("MDM Server: %s", targetServer.Attributes.Name)
	log.Printf("Server URL: %s", targetServer.Attributes.ServerURL)
	log.Printf("Total Assigned Devices: %d", len(deviceIDs))
	log.Printf("Default Server: %v", targetServer.Attributes.DefaultMdmServer)

	log.Println("\n=== Example completed successfully ===")
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
