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

	// Example 1: Get MDM server with all fields
	log.Println("=== Get MDM Server (All Fields) ===")

	server, err := client.MdmServers().GetMdmServer(ctx, exampleMdmServerID)
	if err != nil {
		log.Fatalf("Failed to get MDM server: %v", err)
	}

	log.Printf("MDM Server Details:")
	log.Printf("  ID: %s", server.ID)
	log.Printf("  Type: %s", server.Type)
	log.Printf("  Server Name: %s", server.Attributes.ServerName)
	log.Printf("  Server Type: %s", server.Attributes.ServerType)

	if server.Attributes.CreatedDateTime != "" {
		log.Printf("  Created: %s", server.Attributes.CreatedDateTime)
	}
	if server.Attributes.UpdatedDateTime != "" {
		log.Printf("  Updated: %s", server.Attributes.UpdatedDateTime)
	}

	// Show relationships if available
	if server.Relationships != nil {
		log.Printf("  Relationships: Available")
		log.Printf("    This server has device relationships")
	}

	// Example 2: Get MDM server with specific fields only
	log.Println("\n=== Get MDM Server with Field Filtering ===")

	filteredServer, err := client.MdmServers().GetMdmServer(ctx, exampleMdmServerID,
		axm2.WithFieldsOption("mdmServers", []string{
			"serverName",
			"serverType",
			"createdDateTime",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered MDM server: %v", err)
	}

	log.Printf("Filtered MDM Server Details:")
	log.Printf("  Server Name: %s", filteredServer.Attributes.ServerName)
	log.Printf("  Server Type: %s", filteredServer.Attributes.ServerType)
	log.Printf("  Created: %s", filteredServer.Attributes.CreatedDateTime)

	// Resource links
	if server.Links.Self != "" {
		log.Printf("\n=== Resource Links ===")
		log.Printf("  Self Link: %s", server.Links.Self)
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("Server types you might see:")
	log.Printf("  - MDM: Standard Mobile Device Management server")
	log.Printf("  - APPLE_CONFIGURATOR: Apple Configurator server")
	log.Printf("")
	log.Printf("Use this server ID (%s) with other examples:", server.ID)
	log.Printf("  - get_mdm_server_devices: Get devices assigned to this server")
	log.Printf("  - assign_device_to_mdm_server: Assign devices to this server")
}
