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

	// Example 1: Get all MDM servers with default settings
	log.Println("=== Get All MDM Servers ===")

	servers, err := client.MdmServers().GetMdmServers(ctx)
	if err != nil {
		log.Fatalf("Failed to get MDM servers: %v", err)
	}

	log.Printf("Found %d MDM servers in organization:", len(servers))
	for i, server := range servers {
		log.Printf("  Server %d:", i+1)
		log.Printf("    ID: %s", server.ID)
		log.Printf("    Type: %s", server.Type)
		log.Printf("    Server Name: %s", server.Attributes.ServerName)
		log.Printf("    Server Type: %s", server.Attributes.ServerType)
		if server.Attributes.CreatedDateTime != "" {
			log.Printf("    Created: %s", server.Attributes.CreatedDateTime)
		}
		if server.Attributes.UpdatedDateTime != "" {
			log.Printf("    Updated: %s", server.Attributes.UpdatedDateTime)
		}

		// Show relationships if available
		if server.Relationships != nil {
			log.Printf("    Has device relationships: Yes")
		}

		log.Println() // Empty line for readability
	}

	// Example 2: Get MDM servers with field filtering
	log.Println("=== Get MDM Servers with Field Filtering ===")

	filteredServers, err := client.MdmServers().GetMdmServers(ctx,
		axm2.WithFieldsOption("mdmServers", []string{
			"serverName",
			"serverType",
			"createdDateTime",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered MDM servers: %v", err)
	}

	log.Printf("Filtered MDM servers:")
	for i, server := range filteredServers {
		log.Printf("  Server %d:", i+1)
		log.Printf("    Name: %s", server.Attributes.ServerName)
		log.Printf("    Type: %s", server.Attributes.ServerType)
		log.Printf("    Created: %s", server.Attributes.CreatedDateTime)
	}

	// Example 3: Get MDM servers with pagination limit
	log.Println("\n=== Get MDM Servers with Pagination Limit ===")

	limitedServers, err := client.MdmServers().GetMdmServers(ctx,
		axm2.WithLimitOption(5), // Limit to 5 servers per page
	)
	if err != nil {
		log.Fatalf("Failed to get limited MDM servers: %v", err)
	}

	log.Printf("Retrieved %d MDM servers with pagination limit", len(limitedServers))

	// Note: The SDK automatically handles pagination, so you still get all servers,
	// but the API requests are made in smaller chunks of 5 servers per request

	if len(servers) > 0 {
		log.Printf("\nFirst server details:")
		firstServer := servers[0]
		log.Printf("  Use server ID '%s' in other examples like get_mdm_server", firstServer.ID)
		log.Printf("  Use server ID '%s' in get_mdm_server_devices example", firstServer.ID)
	}
}
