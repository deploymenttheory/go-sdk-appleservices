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

	// First, get a list of devices to check their assignments
	log.Println("\n=== Finding Devices to Check ===")
	devices, err := client.GetOrgDevices(ctx,
		axm2.WithLimitOption(5),
		axm2.WithFieldsOption("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"status",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	if len(devices) == 0 {
		log.Fatalf("No devices found in organization")
	}

	log.Printf("Found %d devices to check", len(devices))

	// Also get the list of available MDM servers for reference
	log.Println("\n=== Available MDM Servers ===")
	servers, err := client.GetMdmServers(ctx,
		axm2.WithFieldsOption("mdmServers", []string{
			"name",
			"serverUrl",
			"defaultMdmServer",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get MDM servers: %v", err)
	}

	log.Printf("Available MDM servers in organization:")
	serverMap := make(map[string]*axm2.MdmServer)
	for i, server := range servers {
		log.Printf("  %d. %s (ID: %s) - Default: %v",
			i+1, server.Attributes.Name, server.ID, server.Attributes.DefaultMdmServer)
		serverMap[server.ID] = &server
	}

	// Example: Check assigned MDM server for each device
	log.Println("\n=== Get Device Assigned MDM Server ===")

	for i, device := range devices {
		log.Printf("\nChecking Device %d:", i+1)
		log.Printf("  Device ID: %s", device.ID)
		log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
		log.Printf("  Model: %s", device.Attributes.DeviceModel)
		log.Printf("  Status: %s", device.Attributes.Status)

		// Get the assigned MDM server for this device
		assignedServerID, err := client.GetDeviceAssignedMdmServer(ctx, device.ID)
		if err != nil {
			log.Printf("  ❌ Failed to get assigned MDM server: %v", err)
			continue
		}

		if assignedServerID == "" {
			log.Printf("  📝 Assignment Status: Not assigned to any MDM server")
			log.Printf("  💡 This device is available for assignment")
		} else {
			log.Printf("  ✅ Assignment Status: Assigned to MDM server")
			log.Printf("  📋 Assigned Server ID: %s", assignedServerID)

			// Look up server details if we have them
			if server, exists := serverMap[assignedServerID]; exists {
				log.Printf("  🖥️  Server Name: %s", server.Attributes.Name)
				log.Printf("  🌐 Server URL: %s", server.Attributes.ServerURL)
				log.Printf("  ⭐ Default Server: %v", server.Attributes.DefaultMdmServer)
			} else {
				// If we don't have the server details, fetch them
				serverDetails, err := client.GetMdmServer(ctx, assignedServerID,
					axm2.WithFieldsOption("mdmServers", []string{"name", "serverUrl", "defaultMdmServer"}),
				)
				if err != nil {
					log.Printf("  ⚠️  Could not fetch server details: %v", err)
				} else {
					log.Printf("  🖥️  Server Name: %s", serverDetails.Attributes.Name)
					log.Printf("  🌐 Server URL: %s", serverDetails.Attributes.ServerURL)
					log.Printf("  ⭐ Default Server: %v", serverDetails.Attributes.DefaultMdmServer)
				}
			}
		}
	}

	// Summary statistics
	log.Println("\n=== Assignment Summary ===")
	assignedCount := 0
	unassignedCount := 0
	errorCount := 0

	for _, device := range devices {
		assignedServerID, err := client.GetDeviceAssignedMdmServer(ctx, device.ID)
		if err != nil {
			errorCount++
		} else if assignedServerID == "" {
			unassignedCount++
		} else {
			assignedCount++
		}
	}

	log.Printf("Total devices checked: %d", len(devices))
	log.Printf("✅ Assigned to MDM server: %d", assignedCount)
	log.Printf("📝 Not assigned: %d", unassignedCount)
	log.Printf("❌ Errors checking: %d", errorCount)

	if unassignedCount > 0 {
		log.Println("\n💡 Tip: Unassigned devices can be assigned to an MDM server using the")
		log.Println("   assign_device_to_mdm_server example")
	}

	log.Println("\n=== Example completed successfully ===")
}
