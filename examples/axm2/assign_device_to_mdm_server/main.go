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
	exampleDeviceID    = "XABC123X0ABC123X0"                // Replace with actual device ID
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

	// Example: Assign a device to an MDM server
	log.Println("=== Assign Device to MDM Server ===")

	activity, err := client.OrgDeviceActivities().AssignDevice(ctx, exampleDeviceID, exampleMdmServerID)
	if err != nil {
		log.Fatalf("Failed to assign device: %v", err)
	}

	log.Printf("Assignment activity created:")
	log.Printf("  Activity ID: %s", activity.ID)
	log.Printf("  Type: %s", activity.Type)
	log.Printf("  Status: %s", activity.Attributes.Status)
	log.Printf("  Sub Status: %s", activity.Attributes.SubStatus)
	log.Printf("  Created: %s", activity.Attributes.CreatedDateTime)

	if activity.Attributes.CompletedDateTime != nil {
		log.Printf("  Completed: %s", activity.Attributes.CompletedDateTime)
	}

	if activity.Attributes.DownloadUrl != "" {
		log.Printf("  Download URL: %s", activity.Attributes.DownloadUrl)
	}

	log.Printf("Use the returned activity ID (%s) with get_org_device_activity to monitor progress", activity.ID)
}
