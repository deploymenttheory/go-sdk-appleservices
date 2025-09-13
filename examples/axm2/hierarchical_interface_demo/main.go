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
	exampleDeviceID    = "XABC123X0ABC123X0"                    // Replace with actual device ID
	exampleMdmServerID = "1F97349736CF4614A94F624E705841AD"     // Replace with actual MDM server ID
	exampleActivityID  = "b1481656-b267-480d-b284-a809eed8b041" // Replace with actual activity ID

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

	// Demonstrate the new hierarchical interface structure
	log.Println("=== Hierarchical Interface Demo ===")
	log.Printf("Apple AXM SDK now uses endpoint-based service organization:")
	log.Printf("")

	// 1. OrgDevices Service - /v1/orgDevices endpoint operations
	log.Println("1. OrgDevices Service (/v1/orgDevices):")
	log.Printf("   client.OrgDevices().GetOrgDevices(ctx, opts...)")
	log.Printf("   client.OrgDevices().GetOrgDevice(ctx, deviceID, opts...)")
	log.Printf("   client.OrgDevices().GetAssignedMdmServer(ctx, deviceID, opts...)")
	log.Printf("   client.OrgDevices().GetAssignedMdmServerInfo(ctx, deviceID, opts...)")
	log.Printf("")

	// 2. MdmServers Service - /v1/mdmServers endpoint operations
	log.Println("2. MdmServers Service (/v1/mdmServers):")
	log.Printf("   client.MdmServers().GetMdmServers(ctx, opts...)")
	log.Printf("   client.MdmServers().GetMdmServer(ctx, serverID, opts...)")
	log.Printf("   client.MdmServers().GetDevices(ctx, serverID, opts...)")
	log.Printf("")

	// 3. OrgDeviceActivities Service - /v1/orgDeviceActivities endpoint operations
	log.Println("3. OrgDeviceActivities Service (/v1/orgDeviceActivities):")
	log.Printf("   client.OrgDeviceActivities().GetActivity(ctx, activityID, opts...)")
	log.Printf("   client.OrgDeviceActivities().AssignDevice(ctx, deviceID, mdmServerID)")
	log.Printf("   client.OrgDeviceActivities().UnassignDevice(ctx, deviceID)")
	log.Printf("   client.OrgDeviceActivities().AssignDevices(ctx, deviceIDs, mdmServerID)")
	log.Printf("   client.OrgDeviceActivities().UnassignDevices(ctx, deviceIDs)")
	log.Printf("")

	// Demonstrate actual usage (if you have valid credentials)
	log.Println("=== Live Demo (requires valid credentials) ===")

	// Example 1: Get organization devices
	log.Printf("Getting organization devices...")
	devices, err := client.OrgDevices().GetOrgDevices(ctx,
		axm2.WithLimitOption(3), // Limit for demo
	)
	if err != nil {
		log.Printf("  Error: %v (expected if credentials are example values)", err)
	} else {
		log.Printf("  Found %d devices", len(devices))
	}

	// Example 2: Get MDM servers
	log.Printf("Getting MDM servers...")
	servers, err := client.MdmServers().GetMdmServers(ctx,
		axm2.WithLimitOption(3), // Limit for demo
	)
	if err != nil {
		log.Printf("  Error: %v (expected if credentials are example values)", err)
	} else {
		log.Printf("  Found %d MDM servers", len(servers))
	}

	// Example 3: Get device activity
	log.Printf("Getting device activity...")
	activity, err := client.OrgDeviceActivities().GetActivity(ctx, exampleActivityID)
	if err != nil {
		log.Printf("  Error: %v (expected if credentials/IDs are example values)", err)
	} else {
		log.Printf("  Activity status: %s", activity.Attributes.Status)
	}

	// Benefits of hierarchical interface
	log.Printf("\n=== Benefits of Hierarchical Interface ===")
	log.Printf("✅ Clear separation by Apple API endpoints")
	log.Printf("✅ Better organization and discoverability")
	log.Printf("✅ Matches Apple's API documentation structure")
	log.Printf("✅ Easier to understand which operations belong together")
	log.Printf("✅ Future-proof for new endpoints")
	log.Printf("")

	// Migration guide
	log.Printf("=== Migration from Flat Interface ===")
	log.Printf("Old: client.GetOrgDevices(ctx, opts...)")
	log.Printf("New: client.OrgDevices().GetOrgDevices(ctx, opts...)")
	log.Printf("")
	log.Printf("Old: client.AssignDeviceToMdmServer(ctx, deviceID, mdmServerID)")
	log.Printf("New: client.OrgDeviceActivities().AssignDevice(ctx, deviceID, mdmServerID)")
	log.Printf("")
	log.Printf("Old: client.GetMdmServerDevices(ctx, serverID, opts...)")
	log.Printf("New: client.MdmServers().GetDevices(ctx, serverID, opts...)")
}
