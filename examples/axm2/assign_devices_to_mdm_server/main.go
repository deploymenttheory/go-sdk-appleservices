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

	// Example: Assign multiple devices to an MDM server (bulk assignment)
	log.Println("=== Assign Multiple Devices to MDM Server (Bulk) ===")

	log.Printf("Assigning %d devices to MDM server '%s':", len(exampleDeviceIDs), exampleMdmServerID)
	for i, deviceID := range exampleDeviceIDs {
		log.Printf("  Device %d: %s", i+1, deviceID)
	}

	activity, err := client.OrgDeviceActivities().AssignDevices(ctx, exampleDeviceIDs, exampleMdmServerID)
	if err != nil {
		log.Fatalf("Failed to assign devices: %v", err)
	}

	log.Printf("\nBulk assignment activity created:")
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
		log.Printf("  ^ Download this CSV file for detailed assignment results")
	}

	// Cross-reference: Get MDM server details
	log.Printf("\n=== Cross-Reference: MDM Server Details ===")
	mdmServer, err := client.MdmServers().GetMdmServer(ctx, exampleMdmServerID,
		axm2.WithFieldsOption("mdmServers", []string{
			"serverName",
			"serverType",
		}),
	)
	if err != nil {
		log.Printf("Failed to get MDM server details: %v", err)
	} else {
		log.Printf("Target MDM Server:")
		log.Printf("  Server Name: %s", mdmServer.Attributes.ServerName)
		log.Printf("  Server Type: %s", mdmServer.Attributes.ServerType)
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("Bulk assignment allows you to assign multiple devices in a single API call.")
	log.Printf("This is more efficient than individual assignments for large device sets.")
	log.Printf("")
	log.Printf("Activity Status Meanings:")
	log.Printf("  - IN_PROGRESS: Assignment is being processed")
	log.Printf("  - COMPLETED: Assignment completed successfully")
	log.Printf("  - FAILED: Assignment failed")
	log.Printf("")
	log.Printf("Sub-Status Details:")
	log.Printf("  - SUBMITTED: Request received and queued")
	log.Printf("  - COMPLETED_WITH_SUCCESS: All devices assigned successfully")
	log.Printf("  - COMPLETED_WITH_ERRORS: Some devices failed to assign")
	log.Printf("")
	log.Printf("Monitor progress using:")
	log.Printf("  - get_org_device_activity example with activity ID: %s", activity.ID)
	log.Printf("  - client.OrgDeviceActivities().GetActivity(ctx, \"%s\")", activity.ID)
	log.Printf("")
	log.Printf("For individual device assignment, use:")
	log.Printf("  - assign_device_to_mdm_server example")
	log.Printf("  - client.OrgDeviceActivities().AssignDevice(ctx, deviceID, mdmServerID)")
}
