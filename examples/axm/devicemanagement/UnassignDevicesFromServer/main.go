package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/devicemanagement"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Unassign Devices from Server Example ===")

	// Use credentials directly for testing
	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	// Parse the private key
	privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Create AXM client
	axmClient, err := axm.NewClientBuilder().
		WithJWTAuth(keyID, issuerID, privateKey).
		WithDebug(true).
		Build()

	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}

	// Create service clients
	dmClient := devicemanagement.NewClient(axmClient)
	devicesClient := devices.NewClient(axmClient)

	// Create context
	ctx := context.Background()

	// Step 1: Get MDM servers
	fmt.Println("\nStep 1: Getting available MDM servers...")

	serversResponse, err := dmClient.GetDeviceManagementServices(ctx, &devicemanagement.GetMDMServersOptions{
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Error getting MDM servers: %v", err)
	}

	if len(serversResponse.Data) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	fmt.Printf("Found %d MDM servers:\n", len(serversResponse.Data))
	for i, server := range serversResponse.Data {
		serverName := "Unknown"
		if server.Attributes != nil {
			serverName = server.Attributes.ServerName
		}
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, serverName, server.ID)
	}

	// Use the first MDM server for examples
	targetMDMServerID := serversResponse.Data[0].ID
	targetServerName := "Unknown"
	if serversResponse.Data[0].Attributes != nil {
		targetServerName = serversResponse.Data[0].Attributes.ServerName
	}

	fmt.Printf("\nUsing MDM Server: %s (ID: %s)\n", targetServerName, targetMDMServerID)

	// Step 2: Get organization devices
	fmt.Println("\nStep 2: Getting organization devices...")

	devicesResponse, err := devicesClient.GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Error getting organization devices: %v", err)
	}

	if len(devicesResponse.Data) == 0 {
		log.Fatalf("No devices found in organization")
	}

	fmt.Printf("Found %d devices in organization\n", len(devicesResponse.Data))

	// Step 3: Find assigned devices
	fmt.Println("\nStep 3: Finding devices assigned to the target server...")

	var assignedDevices []devices.OrgDevice

	for _, device := range devicesResponse.Data {
		// Check if device is assigned to our target server
		serverLinkage, err := dmClient.GetAssignedDeviceManagementServiceIDForADevice(ctx, device.ID)
		if err != nil {
			continue // Skip devices with errors
		}

		if serverLinkage.Data.ID == targetMDMServerID {
			assignedDevices = append(assignedDevices, device)
		}
	}

	fmt.Printf("Found %d devices assigned to server %s\n", len(assignedDevices), targetServerName)

	// If no devices are assigned, assign some first for demonstration
	if len(assignedDevices) == 0 {
		fmt.Println("\nNo devices assigned to target server. Assigning some devices first for demonstration...")

		// Assign first 2 devices to the server
		maxToAssign := min(2, len(devicesResponse.Data))
		var deviceIDsToAssign []string

		for i := 0; i < maxToAssign; i++ {
			deviceIDsToAssign = append(deviceIDsToAssign, devicesResponse.Data[i].ID)
		}

		_, err = dmClient.AssignDevicesToServer(ctx, targetMDMServerID, deviceIDsToAssign)
		if err != nil {
			log.Printf("Error assigning devices for demo: %v", err)
		} else {
			fmt.Printf("Assigned %d devices to server for demonstration\n", len(deviceIDsToAssign))
			// Wait for assignment to process
			time.Sleep(3 * time.Second)

			// Update our assigned devices list
			for i := 0; i < maxToAssign; i++ {
				assignedDevices = append(assignedDevices, devicesResponse.Data[i])
			}
		}
	}

	// Example 1: Unassign a single device from server
	fmt.Println("\n=== Example 1: Unassign Single Device from Server ===")

	if len(assignedDevices) > 0 {
		singleDevice := assignedDevices[0]
		fmt.Printf("Unassigning device %s (Serial: %s) from server %s...\n",
			singleDevice.ID, singleDevice.Attributes.SerialNumber, targetServerName)

		singleDeviceIDs := []string{singleDevice.ID}

		unassignResponse, err := dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, singleDeviceIDs)
		if err != nil {
			log.Printf("Error unassigning single device: %v", err)
		} else {
			fmt.Printf("Unassignment successful!\n")
			fmt.Printf("  Activity ID: %s\n", unassignResponse.Data.ID)
			fmt.Printf("  Activity Type: %s\n", unassignResponse.Data.Type)

			if unassignResponse.Data.Attributes != nil {
				fmt.Printf("  Activity: %s\n", unassignResponse.Data.Attributes.ActivityType)
				if unassignResponse.Data.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", unassignResponse.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
			}

			if unassignResponse.Links != nil {
				fmt.Printf("  Self Link: %s\n", unassignResponse.Links.Self)
			}
		}
	}

	// Example 2: Unassign multiple devices from server
	fmt.Println("\n=== Example 2: Unassign Multiple Devices from Server ===")

	maxDevicesToUnassign := min(3, len(assignedDevices))
	if maxDevicesToUnassign > 1 {
		var multipleDeviceIDs []string
		fmt.Printf("Unassigning %d devices from server %s:\n", maxDevicesToUnassign, targetServerName)

		for i := 0; i < maxDevicesToUnassign; i++ {
			device := assignedDevices[i]
			multipleDeviceIDs = append(multipleDeviceIDs, device.ID)
			fmt.Printf("  - %s (Serial: %s)\n", device.ID, device.Attributes.SerialNumber)
		}

		multipleUnassignResponse, err := dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, multipleDeviceIDs)
		if err != nil {
			log.Printf("Error unassigning multiple devices: %v", err)
		} else {
			fmt.Printf("Multiple device unassignment successful!\n")
			fmt.Printf("  Activity ID: %s\n", multipleUnassignResponse.Data.ID)

			if multipleUnassignResponse.Data.Attributes != nil {
				fmt.Printf("  Activity Type: %s\n", multipleUnassignResponse.Data.Attributes.ActivityType)
				if multipleUnassignResponse.Data.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", multipleUnassignResponse.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
			}
		}
	} else {
		fmt.Println("Not enough assigned devices for multiple unassignment demo")
	}

	// Example 3: Verify unassignment
	fmt.Println("\n=== Example 3: Verify Device Unassignment ===")

	if len(assignedDevices) > 0 {
		deviceToCheck := assignedDevices[0]
		fmt.Printf("Verifying unassignment for device %s (Serial: %s)...\n",
			deviceToCheck.ID, deviceToCheck.Attributes.SerialNumber)

		// Wait a moment for the unassignment to process
		time.Sleep(2 * time.Second)

		serverLinkage, err := dmClient.GetAssignedDeviceManagementServiceIDForADevice(ctx, deviceToCheck.ID)
		if err != nil {
			fmt.Printf("✓ Device appears to be unassigned (no server linkage found)\n")
		} else {
			if serverLinkage.Data.ID == "" {
				fmt.Printf("✓ Device successfully unassigned from server\n")
			} else if serverLinkage.Data.ID != targetMDMServerID {
				fmt.Printf("⚠ Device assigned to different server: %s\n", serverLinkage.Data.ID)
			} else {
				fmt.Printf("✗ Device still appears to be assigned to original server\n")
			}
		}
	}

	// Example 4: Error handling - invalid MDM server ID
	fmt.Println("\n=== Example 4: Error Handling (Invalid MDM Server ID) ===")

	if len(assignedDevices) > 0 {
		invalidServerID := "invalid-mdm-server-id"
		testDeviceIDs := []string{assignedDevices[0].ID}

		_, err = dmClient.UnassignDevicesFromServer(ctx, invalidServerID, testDeviceIDs)
		if err != nil {
			fmt.Printf("Expected error for invalid MDM server ID '%s': %v\n", invalidServerID, err)
		}
	}

	// Example 5: Error handling - empty MDM server ID
	fmt.Println("\n=== Example 5: Error Handling (Empty MDM Server ID) ===")

	if len(assignedDevices) > 0 {
		testDeviceIDs := []string{assignedDevices[0].ID}

		_, err = dmClient.UnassignDevicesFromServer(ctx, "", testDeviceIDs)
		if err != nil {
			fmt.Printf("Expected error for empty MDM server ID: %v\n", err)
		}
	}

	// Example 6: Error handling - empty device IDs
	fmt.Println("\n=== Example 6: Error Handling (Empty Device IDs) ===")

	emptyDeviceIDs := []string{}
	_, err = dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, emptyDeviceIDs)
	if err != nil {
		fmt.Printf("Expected error for empty device IDs: %v\n", err)
	}

	// Example 7: Error handling - invalid device IDs
	fmt.Println("\n=== Example 7: Error Handling (Invalid Device IDs) ===")

	invalidDeviceIDs := []string{"invalid-device-id-1", "invalid-device-id-2"}
	_, err = dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, invalidDeviceIDs)
	if err != nil {
		fmt.Printf("Expected error for invalid device IDs: %v\n", err)
	}

	// Example 8: Unassignment status tracking
	fmt.Println("\n=== Example 8: Unassignment Status Tracking ===")

	if len(assignedDevices) > 0 {
		trackingDevice := assignedDevices[0]
		fmt.Printf("Tracking unassignment status for device %s...\n", trackingDevice.Attributes.SerialNumber)

		// Check initial status
		fmt.Printf("Initial status check...\n")
		initialLinkage, err := dmClient.GetAssignedDeviceManagementServiceIDForADevice(ctx, trackingDevice.ID)
		if err != nil {
			fmt.Printf("  Error checking initial status: %v\n", err)
		} else {
			fmt.Printf("  Initial assignment: %s\n", initialLinkage.Data.ID)
		}

		// Perform unassignment
		fmt.Printf("Performing unassignment...\n")
		trackingDeviceIDs := []string{trackingDevice.ID}
		unassignResponse, err := dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, trackingDeviceIDs)
		if err != nil {
			fmt.Printf("  Unassignment error: %v\n", err)
		} else {
			fmt.Printf("  Unassignment activity created: %s\n", unassignResponse.Data.ID)
		}

		// Check final status
		fmt.Printf("Final status check (after 3 seconds)...\n")
		time.Sleep(3 * time.Second)
		finalLinkage, err := dmClient.GetAssignedDeviceManagementServiceIDForADevice(ctx, trackingDevice.ID)
		if err != nil {
			fmt.Printf("  ✓ Device appears to be unassigned (no linkage found)\n")
		} else {
			if finalLinkage.Data.ID == "" {
				fmt.Printf("  ✓ Unassignment successful!\n")
			} else {
				fmt.Printf("  ⚠ Device still assigned to: %s (unassignment may still be processing)\n", finalLinkage.Data.ID)
			}
		}
	}

	// Example 9: Bulk unassignment from server
	fmt.Println("\n=== Example 9: Bulk Unassignment from Server ===")

	// Get all devices assigned to the target server
	fmt.Printf("Getting all devices assigned to server %s for bulk unassignment...\n", targetServerName)

	serverDeviceLinkages, err := dmClient.GetMDMServerDeviceLinkages(ctx, targetMDMServerID, &devicemanagement.GetMDMServerDeviceLinkagesOptions{
		Limit: 100,
	})
	if err != nil {
		log.Printf("Error getting server device linkages: %v", err)
	} else {
		if len(serverDeviceLinkages.Data) > 0 {
			fmt.Printf("Found %d devices assigned to server\n", len(serverDeviceLinkages.Data))

			// Unassign up to 5 devices for demonstration
			maxBulkUnassign := min(5, len(serverDeviceLinkages.Data))
			var bulkDeviceIDs []string

			for i := 0; i < maxBulkUnassign; i++ {
				bulkDeviceIDs = append(bulkDeviceIDs, serverDeviceLinkages.Data[i].ID)
			}

			fmt.Printf("Performing bulk unassignment of %d devices...\n", len(bulkDeviceIDs))
			bulkUnassignResponse, err := dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, bulkDeviceIDs)
			if err != nil {
				log.Printf("Error in bulk unassignment: %v", err)
			} else {
				fmt.Printf("Bulk unassignment successful! Activity ID: %s\n", bulkUnassignResponse.Data.ID)
			}
		} else {
			fmt.Printf("No devices currently assigned to server %s\n", targetServerName)
		}
	}

	// Example 10: Pretty print JSON response
	fmt.Println("\n=== Example 10: Full JSON Response ===")
	if len(assignedDevices) > 0 {
		// Perform one more unassignment for JSON demo
		jsonDemoDeviceIDs := []string{assignedDevices[0].ID}
		jsonResponse, err := dmClient.UnassignDevicesFromServer(ctx, targetMDMServerID, jsonDemoDeviceIDs)
		if err != nil {
			log.Printf("Error in JSON demo unassignment: %v", err)
		} else {
			jsonData, err := json.MarshalIndent(jsonResponse, "", "  ")
			if err != nil {
				log.Printf("Error marshaling response to JSON: %v", err)
			} else {
				fmt.Println(string(jsonData))
			}
		}
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("Note: Device unassignments may take some time to process in Apple's system.")
	fmt.Println("Check the Apple Business Manager portal to verify final unassignment status.")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
