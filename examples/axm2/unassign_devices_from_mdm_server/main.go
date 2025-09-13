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

// Example device IDs - Replace with actual device IDs from your organization
var exampleDeviceIDs = []string{
	"XABC123X0ABC123X0", // Replace with actual device ID
	"YABC123X0ABC123X1", // Replace with actual device ID
	"ZABC123X0ABC123X2", // Replace with actual device ID
}

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

	// Example: Unassign multiple devices from their current MDM servers (bulk unassignment)
	log.Println("=== Unassign Multiple Devices from MDM Servers (Bulk) ===")

	log.Printf("Unassigning %d devices from their current MDM servers:", len(exampleDeviceIDs))
	for i, deviceID := range exampleDeviceIDs {
		log.Printf("  Device %d: %s", i+1, deviceID)
	}

	activity, err := client.OrgDeviceActivities().UnassignDevices(ctx, exampleDeviceIDs)
	if err != nil {
		log.Fatalf("Failed to unassign devices: %v", err)
	}

	log.Printf("\nBulk unassignment activity created:")
	log.Printf("  Activity ID: %s", activity.ID)
	log.Printf("  Type: %s", activity.Type)
	log.Printf("  Status: %s", activity.Attributes.Status)
	log.Printf("  Sub Status: %s", activity.Attributes.SubStatus)
	log.Printf("  Created: %s", activity.Attributes.CreatedDateTime)

	if activity.Attributes.CompletedDateTime != nil {
		log.Printf("  Completed: %s", activity.Attributes.CompletedDateTime)
	} else {
		log.Printf("  Completed: In progress...")
	}

	if activity.Attributes.DownloadUrl != "" {
		log.Printf("  Download URL: %s", activity.Attributes.DownloadUrl)
		log.Printf("  ^ Download this CSV file for detailed unassignment results")
	}

	// Cross-reference: Check current assignments for first device
	log.Printf("\n=== Cross-Reference: Current Device Assignments ===")
	if len(exampleDeviceIDs) > 0 {
		firstDeviceID := exampleDeviceIDs[0]
		log.Printf("Checking current assignment for device: %s", firstDeviceID)

		mdmServerID, err := client.OrgDevices().GetAssignedMdmServer(ctx, firstDeviceID)
		if err != nil {
			log.Printf("Failed to get device assignment: %v", err)
		} else if mdmServerID == "" {
			log.Printf("  Device is currently UNASSIGNED")
		} else {
			log.Printf("  Device is currently assigned to MDM server: %s", mdmServerID)
			log.Printf("  After unassignment completes, this device will be UNASSIGNED")
		}
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("Bulk unassignment removes devices from their current MDM server assignments.")
	log.Printf("This is more efficient than individual unassignments for large device sets.")
	log.Printf("")
	log.Printf("Activity Status Meanings:")
	log.Printf("  - IN_PROGRESS: Unassignment is being processed")
	log.Printf("  - COMPLETED: Unassignment completed successfully")
	log.Printf("  - FAILED: Unassignment failed")
	log.Printf("")
	log.Printf("Sub-Status Details:")
	log.Printf("  - SUBMITTED: Request received and queued")
	log.Printf("  - COMPLETED_WITH_SUCCESS: All devices unassigned successfully")
	log.Printf("  - COMPLETED_WITH_ERRORS: Some devices failed to unassign")
	log.Printf("")
	log.Printf("After unassignment:")
	log.Printf("  - Devices will have status: UNASSIGNED")
	log.Printf("  - Devices can be reassigned to any MDM server")
	log.Printf("  - Use assign_devices_to_mdm_server to reassign")
	log.Printf("")
	log.Printf("Monitor progress using:")
	log.Printf("  - get_org_device_activity example with activity ID: %s", activity.ID)
	log.Printf("  - client.OrgDeviceActivities().GetActivity(ctx, \"%s\")", activity.ID)
	log.Printf("")
	log.Printf("For individual device unassignment, use:")
	log.Printf("  - unassign_device_from_mdm_server example")
	log.Printf("  - client.OrgDeviceActivities().UnassignDevice(ctx, deviceID)")
}
