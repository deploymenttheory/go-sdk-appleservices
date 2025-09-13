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
	exampleDeviceID = "XABC123X0ABC123X0" // Replace with actual device ID

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

	// Example 1: Get full MDM server information for a device
	log.Println("=== Get Device Assigned MDM Server Information ===")

	mdmServer, err := client.OrgDevices().GetAssignedMdmServerInfo(ctx, exampleDeviceID)
	if err != nil {
		log.Fatalf("Failed to get device assigned MDM server info: %v", err)
	}

	if mdmServer == nil {
		log.Printf("Device '%s' is not assigned to any MDM server", exampleDeviceID)
		log.Printf("Device status: UNASSIGNED")
	} else {
		log.Printf("Device '%s' is assigned to MDM server:", exampleDeviceID)
		log.Printf("  MDM Server ID: %s", mdmServer.ID)
		log.Printf("  Type: %s", mdmServer.Type)
		log.Printf("  Server Name: %s", mdmServer.Attributes.ServerName)
		log.Printf("  Server Type: %s", mdmServer.Attributes.ServerType)

		if mdmServer.Attributes.CreatedDateTime != "" {
			log.Printf("  Created: %s", mdmServer.Attributes.CreatedDateTime)
		}
		if mdmServer.Attributes.UpdatedDateTime != "" {
			log.Printf("  Updated: %s", mdmServer.Attributes.UpdatedDateTime)
		}

		// Show relationships if available
		if mdmServer.Relationships != nil {
			log.Printf("  Has device relationships: Yes")
		}
	}

	// Example 2: Get MDM server info with field filtering
	log.Println("\n=== Get MDM Server Info with Field Filtering ===")

	filteredMdmServer, err := client.OrgDevices().GetAssignedMdmServerInfo(ctx, exampleDeviceID,
		axm2.WithFieldsOption("mdmServers", []string{
			"serverName",
			"serverType",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered MDM server info: %v", err)
	}

	if filteredMdmServer != nil {
		log.Printf("Filtered MDM Server Information:")
		log.Printf("  Server Name: %s", filteredMdmServer.Attributes.ServerName)
		log.Printf("  Server Type: %s", filteredMdmServer.Attributes.ServerType)
	}

	// Cross-reference: Get device details
	if mdmServer != nil {
		log.Printf("\n=== Cross-Reference: Device Details ===")
		device, err := client.OrgDevices().GetOrgDevice(ctx, exampleDeviceID,
			axm2.WithFieldsOption("orgDevices", []string{
				"serialNumber",
				"deviceModel",
				"status",
				"productFamily",
			}),
		)
		if err != nil {
			log.Printf("Failed to get device details: %v", err)
		} else {
			log.Printf("Device Information:")
			log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
			log.Printf("  Model: %s", device.Attributes.DeviceModel)
			log.Printf("  Status: %s", device.Attributes.Status)
			log.Printf("  Product Family: %s", device.Attributes.ProductFamily)
		}
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("This endpoint returns complete MDM server information for the assigned server.")
	log.Printf("Compare with get_device_assigned_mdm_server which returns only the server ID.")
	log.Printf("")
	log.Printf("Server types you might see:")
	log.Printf("  - MDM: Standard Mobile Device Management server")
	log.Printf("  - APPLE_CONFIGURATOR: Apple Configurator server")
	log.Printf("")
	log.Printf("To modify device assignments:")
	log.Printf("  - assign_device_to_mdm_server: Assign to a different server")
	log.Printf("  - unassign_device_from_mdm_server: Remove current assignment")
	log.Printf("")
	log.Printf("Apple API endpoint: GET /v1/orgDevices/{id}/assignedServer")
}
