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

	// Example: Get all MDM servers (device management services)
	log.Println("\n=== Get MDM Servers ===")

	// Method 1: Get all servers
	servers, err := client.GetMdmServers(ctx)
	if err != nil {
		log.Fatalf("Failed to get MDM servers: %v", err)
	}

	log.Printf("Retrieved %d MDM servers", len(servers))

	// Display server details
	for i, server := range servers {
		log.Printf("MDM Server %d:", i+1)
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
		log.Println()
	}

	// Method 2: Get servers with field filtering
	log.Println("=== Get MDM Servers with Field Filtering ===")
	filteredServers, err := client.GetMdmServers(ctx,
		axm2.WithLimitOption(5),
		axm2.WithFieldsOption("mdmServers", []string{
			"name",
			"serverUrl",
			"defaultMdmServer",
			"organizationName",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered MDM servers: %v", err)
	}

	log.Printf("Retrieved %d filtered MDM servers", len(filteredServers))
	for i, server := range filteredServers {
		log.Printf("Filtered Server %d:", i+1)
		log.Printf("  Name: %s", server.Attributes.Name)
		log.Printf("  Server URL: %s", server.Attributes.ServerURL)
		log.Printf("  Default: %v", server.Attributes.DefaultMdmServer)
		log.Printf("  Organization: %s", server.Attributes.OrganizationName)
		log.Println()
	}

	// Find the default MDM server
	log.Println("=== Default MDM Server ===")
	var defaultServer *axm2.MdmServer
	for i := range servers {
		if servers[i].Attributes.DefaultMdmServer {
			defaultServer = &servers[i]
			break
		}
	}

	if defaultServer != nil {
		log.Printf("Default MDM Server:")
		log.Printf("  ID: %s", defaultServer.ID)
		log.Printf("  Name: %s", defaultServer.Attributes.Name)
		log.Printf("  Server URL: %s", defaultServer.Attributes.ServerURL)
	} else {
		log.Println("No default MDM server found")
	}

	log.Println("\n=== Example completed successfully ===")
}
