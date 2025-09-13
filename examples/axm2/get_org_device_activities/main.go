package main

import (
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

	// Example: Demonstrate device activity limitations
	log.Println("=== Device Activities Information ===")

	log.Println("\n=== Important Note ===")
	log.Println("Apple's API does not support listing all device activities.")
	log.Println("You can only get individual activities by ID using:")
	log.Println("  client.GetOrgDeviceActivity(ctx, activityID)")
	log.Println("")
	log.Println("Activity IDs are returned when you:")
	log.Println("  - Assign devices: client.AssignDeviceToMdmServer()")
	log.Println("  - Unassign devices: client.UnassignDeviceFromMdmServer()")
	log.Println("")
	log.Println("Activities are only available for the past 30 days.")
	log.Println("")
	log.Println("To get a specific activity, use the get_org_device_activity example with a valid activity ID.")
}
