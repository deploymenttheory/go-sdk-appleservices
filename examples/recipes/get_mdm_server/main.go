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

	// First, get a list of MDM servers to find a server ID
	log.Println("\n=== Finding an MDM Server ID ===")
	servers, err := client.GetMdmServers(ctx, axm2.WithLimitOption(1))
	if err != nil {
		log.Fatalf("Failed to get MDM servers: %v", err)
	}

	if len(servers) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	serverID := servers[0].ID
	log.Printf("Using MDM server ID: %s", serverID)
	log.Printf("Server name: %s", servers[0].Attributes.Name)

	// Example: Get a specific MDM server
	log.Println("\n=== Get Specific MDM Server ===")

	// Method 1: Get server with all fields
	server, err := client.GetMdmServer(ctx, serverID)
	if err != nil {
		log.Fatalf("Failed to get MDM server: %v", err)
	}

	// Display detailed server information
	log.Printf("MDM Server Details:")
	log.Printf("  ID: %s", server.ID)
	log.Printf("  Type: %s", server.Type)
	log.Printf("  Name: %s", server.Attributes.Name)
	log.Printf("  Server URL: %s", server.Attributes.ServerURL)
	log.Printf("  Server UUID: %s", server.Attributes.ServerUUID)
	log.Printf("  Default MDM Server: %v", server.Attributes.DefaultMdmServer)
	log.Printf("  Organization ID: %s", server.Attributes.OrganizationID)
	log.Printf("  Organization Name: %s", server.Attributes.OrganizationName)
	log.Printf("  Organization UUID: %s", server.Attributes.OrganizationUUID)
	log.Printf("  Organization Apple ID: %s", server.Attributes.OrganizationAppleID)

	if server.Links.Self != "" {
		log.Printf("  Self Link: %s", server.Links.Self)
	}

	// Method 2: Get server with specific fields only
	log.Println("\n=== Get MDM Server with Field Filtering ===")
	filteredServer, err := client.GetMdmServer(ctx, serverID,
		axm2.WithFieldsOption("mdmServers", []string{
			"name",
			"serverUrl",
			"defaultMdmServer",
			"organizationName",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered MDM server: %v", err)
	}

	log.Printf("Filtered MDM Server:")
	log.Printf("  ID: %s", filteredServer.ID)
	log.Printf("  Name: %s", filteredServer.Attributes.Name)
	log.Printf("  Server URL: %s", filteredServer.Attributes.ServerURL)
	log.Printf("  Default: %v", filteredServer.Attributes.DefaultMdmServer)
	log.Printf("  Organization: %s", filteredServer.Attributes.OrganizationName)

	// Display organization information if this is the default server
	if server.Attributes.DefaultMdmServer {
		log.Println("\n=== Default MDM Server Organization Info ===")
		log.Printf("This is the default MDM server for:")
		log.Printf("  Organization: %s", server.Attributes.OrganizationName)
		log.Printf("  Organization ID: %s", server.Attributes.OrganizationID)
		log.Printf("  Organization UUID: %s", server.Attributes.OrganizationUUID)
		log.Printf("  Organization Apple ID: %s", server.Attributes.OrganizationAppleID)
	}

	log.Println("\n=== Example completed successfully ===")
}
