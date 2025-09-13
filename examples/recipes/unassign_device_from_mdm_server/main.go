package main

import (
	"context"
	"log"
	"time"

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

	// Find an assigned device
	log.Println("\n=== Finding an Assigned Device ===")
	devices, err := client.GetOrgDevices(ctx,
		axm2.WithLimitOption(10),
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

	// Find an assigned device
	var targetDevice *axm2.OrgDevice
	var assignedServerID string
	for i := range devices {
		serverID, err := client.GetDeviceAssignedMdmServer(ctx, devices[i].ID)
		if err != nil {
			log.Printf("Failed to check assignment for device %s: %v", devices[i].ID, err)
			continue
		}
		if serverID != "" {
			targetDevice = &devices[i]
			assignedServerID = serverID
			break
		}
	}

	if targetDevice == nil {
		log.Fatalf("No assigned devices found. Please assign a device first using the assign_device_to_mdm_server example.")
	}

	log.Printf("Selected assigned device:")
	log.Printf("  ID: %s", targetDevice.ID)
	log.Printf("  Serial Number: %s", targetDevice.Attributes.SerialNumber)
	log.Printf("  Model: %s", targetDevice.Attributes.DeviceModel)
	log.Printf("  Status: %s", targetDevice.Attributes.Status)

	// Get details of the currently assigned MDM server
	log.Println("\n=== Current MDM Server Assignment ===")
	currentServer, err := client.GetMdmServer(ctx, assignedServerID,
		axm2.WithFieldsOption("mdmServers", []string{
			"name",
			"serverUrl",
			"defaultMdmServer",
		}),
	)
	if err != nil {
		log.Printf("Failed to get current MDM server details: %v", err)
		log.Printf("Current assigned server ID: %s", assignedServerID)
	} else {
		log.Printf("Currently assigned to:")
		log.Printf("  Server ID: %s", currentServer.ID)
		log.Printf("  Server Name: %s", currentServer.Attributes.Name)
		log.Printf("  Server URL: %s", currentServer.Attributes.ServerURL)
		log.Printf("  Default Server: %v", currentServer.Attributes.DefaultMdmServer)
	}

	// Example: Unassign device from MDM server
	log.Println("\n=== Unassign Device from MDM Server ===")

	log.Printf("Unassigning device %s from its current MDM server...",
		targetDevice.Attributes.SerialNumber)

	// Perform the unassignment
	activity, err := client.UnassignDeviceFromMdmServer(ctx, targetDevice.ID)
	if err != nil {
		log.Fatalf("Failed to unassign device: %v", err)
	}

	log.Printf("✅ Unassignment activity created successfully!")
	log.Printf("Activity Details:")
	log.Printf("  Activity ID: %s", activity.ID)
	log.Printf("  Type: %s", activity.Type)
	log.Printf("  Status: %s", activity.Attributes.Status)
	log.Printf("  Sub Status: %s", activity.Attributes.SubStatus)
	log.Printf("  Created: %s", activity.Attributes.CreatedDateTime.Format(time.RFC3339))

	if activity.Attributes.CompletedDateTime != nil {
		log.Printf("  Completed: %s", activity.Attributes.CompletedDateTime.Format(time.RFC3339))
	} else {
		log.Printf("  Completed: Not yet completed")
	}

	if activity.Attributes.DownloadUrl != "" {
		log.Printf("  Download URL: %s", activity.Attributes.DownloadUrl)
	}

	// Monitor the unassignment activity
	log.Println("\n=== Monitoring Unassignment Progress ===")
	log.Println("Checking unassignment status...")

	maxChecks := 10
	checkInterval := 3 * time.Second

	for i := 0; i < maxChecks; i++ {
		// Wait before checking (except for first check)
		if i > 0 {
			time.Sleep(checkInterval)
		}

		// Get updated activity status
		updatedActivity, err := client.GetOrgDeviceActivity(ctx, activity.ID)
		if err != nil {
			log.Printf("❌ Failed to get activity status (attempt %d): %v", i+1, err)
			continue
		}

		log.Printf("Check %d - Status: %s, SubStatus: %s", i+1, updatedActivity.Attributes.Status, updatedActivity.Attributes.SubStatus)

		// Check if completed
		if updatedActivity.Attributes.CompletedDateTime != nil {
			log.Printf("✅ Unassignment completed at: %s",
				updatedActivity.Attributes.CompletedDateTime.Format(time.RFC3339))

			if updatedActivity.Attributes.SubStatus != "COMPLETED_WITH_SUCCESS" {
				log.Printf("❌ Unassignment completed with issues - SubStatus: %s", updatedActivity.Attributes.SubStatus)
			} else {
				log.Printf("🎉 Unassignment completed successfully!")
			}

			if updatedActivity.Attributes.DownloadUrl != "" {
				log.Printf("📄 Activity report available: %s", updatedActivity.Attributes.DownloadUrl)
			}
			break
		} else if updatedActivity.Attributes.Status == "FAILED" {
			log.Printf("❌ Unassignment failed - SubStatus: %s", updatedActivity.Attributes.SubStatus)
			break
		} else {
			log.Printf("⏳ Unassignment still in progress - SubStatus: %s", updatedActivity.Attributes.SubStatus)
		}

		if i == maxChecks-1 {
			log.Printf("⚠️ Stopped checking after %d attempts. Unassignment may still be in progress.", maxChecks)
			log.Printf("   You can check the status later using activity ID: %s", activity.ID)
		}
	}

	// Verify the unassignment
	log.Println("\n=== Verifying Unassignment ===")
	newAssignedServerID, err := client.GetDeviceAssignedMdmServer(ctx, targetDevice.ID)
	if err != nil {
		log.Printf("❌ Failed to verify unassignment: %v", err)
	} else if newAssignedServerID == "" {
		log.Printf("✅ Verification successful! Device is now unassigned from all MDM servers")
		log.Printf("   The device is available for assignment to any MDM server")
	} else if newAssignedServerID == assignedServerID {
		log.Printf("⚠️ Device is still assigned to the same MDM server: %s", newAssignedServerID)
		log.Printf("   The unassignment may not have completed yet")
	} else {
		log.Printf("⚠️ Device is now assigned to a different MDM server: %s", newAssignedServerID)
	}

	// Show bulk unassignment example
	log.Println("\n=== Bulk Unassignment Example ===")
	log.Println("For bulk unassignment of multiple devices, use:")
	log.Println("  client.UnassignDevicesFromMdmServer(ctx, []string{deviceID1, deviceID2, ...})")

	// Show what happens next
	log.Println("\n=== Next Steps ===")
	if newAssignedServerID == "" {
		log.Println("✅ Device is now unassigned and can be:")
		log.Println("   - Assigned to a different MDM server")
		log.Println("   - Left unmanaged (not recommended for production)")
		log.Println("   - Assigned to the default MDM server")
		log.Println("\n💡 To reassign this device, use the assign_device_to_mdm_server example")
	}

	log.Println("\n=== Example completed successfully ===")
}
