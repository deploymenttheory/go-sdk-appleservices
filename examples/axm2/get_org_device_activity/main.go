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
	exampleActivityID = "b1481656-b267-480d-b284-a809eed8b041" // Replace with actual activity ID

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

	// Example: Get organization device activity by ID
	log.Println("=== Get Organization Device Activity ===")

	// Method 1: Get activity with all fields
	activity, err := client.OrgDeviceActivities().GetActivity(ctx, exampleActivityID)
	if err != nil {
		log.Fatalf("Failed to get device activity: %v", err)
	}

	log.Printf("Device Activity:")
	log.Printf("  ID: %s", activity.ID)
	log.Printf("  Type: %s", activity.Type)
	log.Printf("  Status: %s", activity.Attributes.Status)
	log.Printf("  Sub Status: %s", activity.Attributes.SubStatus)
	log.Printf("  Created: %s", activity.Attributes.CreatedDateTime)

	if activity.Attributes.CompletedDateTime != nil {
		log.Printf("  Completed: %s", activity.Attributes.CompletedDateTime)
	} else {
		log.Printf("  Completed: Not yet completed")
	}

	if activity.Attributes.DownloadUrl != "" {
		log.Printf("  Download URL: %s", activity.Attributes.DownloadUrl)
		log.Printf("  ^ Download this CSV file for detailed activity information")
	}

	if activity.Links.Self != "" {
		log.Printf("  Self Link: %s", activity.Links.Self)
	}

	// Method 2: Get activity with specific fields only
	log.Println("\n=== Get Activity with Field Filtering ===")
	filteredActivity, err := client.OrgDeviceActivities().GetActivity(ctx, exampleActivityID,
		axm2.WithFieldsOption("orgDeviceActivities", []string{
			"status",
			"subStatus",
			"createdDateTime",
			"completedDateTime",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered activity: %v", err)
	}

	log.Printf("Filtered Activity:")
	log.Printf("  Status: %s", filteredActivity.Attributes.Status)
	log.Printf("  Sub Status: %s", filteredActivity.Attributes.SubStatus)
	log.Printf("  Created: %s", filteredActivity.Attributes.CreatedDateTime)
	if filteredActivity.Attributes.CompletedDateTime != nil {
		log.Printf("  Completed: %s", filteredActivity.Attributes.CompletedDateTime)
	}
}
