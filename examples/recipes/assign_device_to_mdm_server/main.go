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

	// Find an unassigned device
	log.Println("\n=== Finding an Unassigned Device ===")
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

	// Find an unassigned device
	var targetDevice *axm2.OrgDevice
	for i := range devices {
		assignedServerID, err := client.GetDeviceAssignedMdmServer(ctx, devices[i].ID)
		if err != nil {
			log.Printf("Failed to check assignment for device %s: %v", devices[i].ID, err)
			continue
		}
		if assignedServerID == "" {
			targetDevice = &devices[i]
			break
		}
	}

	if targetDevice == nil {
		log.Println("No unassigned devices found. Using the first device for demonstration.")
		log.Println("Note: This may reassign a device from one MDM server to another.")
		targetDevice = &devices[0]
	}

	log.Printf("Selected device:")
	log.Printf("  ID: %s", targetDevice.ID)
	log.Printf("  Serial Number: %s", targetDevice.Attributes.SerialNumber)
	log.Printf("  Model: %s", targetDevice.Attributes.DeviceModel)
	log.Printf("  Status: %s", targetDevice.Attributes.Status)

	// Get available MDM servers
	log.Println("\n=== Finding Target MDM Server ===")
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

	if len(servers) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	// Use the default server if available, otherwise use the first one
	var targetServer *axm2.MdmServer
	for i := range servers {
		if servers[i].Attributes.DefaultMdmServer {
			targetServer = &servers[i]
			break
		}
	}
	if targetServer == nil {
		targetServer = &servers[0]
	}

	log.Printf("Selected MDM server:")
	log.Printf("  ID: %s", targetServer.ID)
	log.Printf("  Name: %s", targetServer.Attributes.Name)
	log.Printf("  Server URL: %s", targetServer.Attributes.ServerURL)
	log.Printf("  Default Server: %v", targetServer.Attributes.DefaultMdmServer)

	// Example: Assign device to MDM server
	log.Println("\n=== Assign Device to MDM Server ===")

	log.Printf("Assigning device %s to MDM server %s...",
		targetDevice.Attributes.SerialNumber, targetServer.Attributes.Name)

	// Perform the assignment
	activity, err := client.AssignDeviceToMdmServer(ctx, targetDevice.ID, targetServer.ID)
	if err != nil {
		log.Fatalf("Failed to assign device: %v", err)
	}

	log.Printf("✅ Assignment activity created successfully!")
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

	// Monitor the assignment activity
	log.Println("\n=== Monitoring Assignment Progress ===")
	log.Println("Checking assignment status...")

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
			log.Printf("✅ Assignment completed at: %s",
				updatedActivity.Attributes.CompletedDateTime.Format(time.RFC3339))

			if updatedActivity.Attributes.SubStatus != "COMPLETED_WITH_SUCCESS" {
				log.Printf("❌ Assignment completed with issues - SubStatus: %s", updatedActivity.Attributes.SubStatus)
			} else {
				log.Printf("🎉 Assignment completed successfully!")
			}

			if updatedActivity.Attributes.DownloadUrl != "" {
				log.Printf("📄 Activity report available: %s", updatedActivity.Attributes.DownloadUrl)
			}
			break
		} else if updatedActivity.Attributes.Status == "FAILED" {
			log.Printf("❌ Assignment failed - SubStatus: %s", updatedActivity.Attributes.SubStatus)
			break
		} else {
			log.Printf("⏳ Assignment still in progress - SubStatus: %s", updatedActivity.Attributes.SubStatus)
		}

		if i == maxChecks-1 {
			log.Printf("⚠️ Stopped checking after %d attempts. Assignment may still be in progress.", maxChecks)
			log.Printf("   You can check the status later using activity ID: %s", activity.ID)
		}
	}

	// Verify the assignment
	log.Println("\n=== Verifying Assignment ===")
	assignedServerID, err := client.GetDeviceAssignedMdmServer(ctx, targetDevice.ID)
	if err != nil {
		log.Printf("❌ Failed to verify assignment: %v", err)
	} else if assignedServerID == targetServer.ID {
		log.Printf("✅ Verification successful! Device is now assigned to %s", targetServer.Attributes.Name)
	} else if assignedServerID == "" {
		log.Printf("⚠️ Device is still not assigned to any MDM server")
	} else {
		log.Printf("⚠️ Device is assigned to a different MDM server: %s", assignedServerID)
	}

	// Show bulk assignment example
	log.Println("\n=== Bulk Assignment Example ===")
	log.Println("For bulk assignment of multiple devices, use:")
	log.Println("  client.AssignDevicesToMdmServer(ctx, []string{deviceID1, deviceID2, ...}, serverID)")

	log.Println("\n=== Example completed successfully ===")
}
