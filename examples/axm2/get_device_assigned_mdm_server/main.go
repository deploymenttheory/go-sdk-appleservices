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

	// Example: Get assigned MDM server ID for a device
	log.Println("=== Get Device Assigned MDM Server ID ===")

	mdmServerID, err := client.OrgDevices().GetAssignedMdmServer(ctx, exampleDeviceID)
	if err != nil {
		log.Fatalf("Failed to get device assigned MDM server: %v", err)
	}

	if mdmServerID == "" {
		log.Printf("Device '%s' is not assigned to any MDM server", exampleDeviceID)
		log.Printf("Device status: UNASSIGNED")
		log.Printf("")
		log.Printf("To assign this device to an MDM server, use:")
		log.Printf("  - assign_device_to_mdm_server example")
		log.Printf("  - client.AssignDeviceToMdmServer(ctx, deviceID, mdmServerID)")
	} else {
		log.Printf("Device '%s' is assigned to MDM server:", exampleDeviceID)
		log.Printf("  MDM Server ID: %s", mdmServerID)
		log.Printf("  Device status: ASSIGNED")

		// Cross-reference: Get device details
		log.Printf("\n=== Cross-Reference: Device Details ===")
		device, err := client.OrgDevices().GetOrgDevice(ctx, exampleDeviceID,
			axm2.WithFieldsOption("orgDevices", []string{
				"serialNumber",
				"deviceModel",
				"status",
			}),
		)
		if err != nil {
			log.Printf("Failed to get device details: %v", err)
		} else {
			log.Printf("Device Information:")
			log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
			log.Printf("  Model: %s", device.Attributes.DeviceModel)
			log.Printf("  Status: %s", device.Attributes.Status)
		}

		// Cross-reference: Get MDM server details
		log.Printf("\n=== Cross-Reference: MDM Server Details ===")
		mdmServer, err := client.MdmServers().GetMdmServer(ctx, mdmServerID,
			axm2.WithFieldsOption("mdmServers", []string{
				"serverName",
				"serverType",
			}),
		)
		if err != nil {
			log.Printf("Failed to get MDM server details: %v", err)
		} else {
			log.Printf("MDM Server Information:")
			log.Printf("  Server Name: %s", mdmServer.Attributes.ServerName)
			log.Printf("  Server Type: %s", mdmServer.Attributes.ServerType)
		}
	}

	// Usage notes
	log.Printf("\n=== Usage Notes ===")
	log.Printf("This endpoint returns only the MDM server ID (relationship linkage).")
	log.Printf("For full MDM server details, use:")
	log.Printf("  - get_device_assigned_mdm_server_info example")
	log.Printf("  - client.OrgDevices().GetAssignedMdmServerInfo(ctx, deviceID)")
	log.Printf("")
	log.Printf("To modify device assignments:")
	log.Printf("  - assign_device_to_mdm_server: Assign to a specific server")
	log.Printf("  - unassign_device_from_mdm_server: Remove current assignment")
	log.Printf("")
	log.Printf("Apple API endpoint: GET /v1/orgDevices/{id}/relationships/assignedServer")
}
