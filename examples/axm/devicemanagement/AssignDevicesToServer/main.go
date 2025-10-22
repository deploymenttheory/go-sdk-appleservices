package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devicemanagement"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Assign Devices to Server Example ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	// Parse the private key
	privateKey, err := client.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Get MDM servers
	fmt.Println("\nStep 1: Getting available MDM servers...")

	serversResponse, err := client.
		DeviceManagement.
		GetDeviceManagementServices(ctx, &devicemanagement.RequestQueryOptions{
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

	devicesResponse, err := client.
		Devices.
		GetOrganizationDevices(ctx, &devices.RequestQueryOptions{
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

	// Step 3: Find unassigned devices
	fmt.Println("\nStep 3: Finding unassigned devices...")

	var unassignedDevices []devices.OrgDevice

	for _, device := range devicesResponse.Data {
		serverLinkage, err := client.
			DeviceManagement.
			GetAssignedDeviceManagementServiceIDForADevice(ctx, device.ID)
		if err != nil {
			unassignedDevices = append(unassignedDevices, device)
			continue
		}

		if serverLinkage.Data.ID == "" {
			unassignedDevices = append(unassignedDevices, device)
		}
	}

	fmt.Printf("Found %d potentially unassigned devices\n", len(unassignedDevices))

	if len(unassignedDevices) == 0 {
		fmt.Println("No unassigned devices found. Using first device for demonstration...")
		unassignedDevices = append(unassignedDevices, devicesResponse.Data[0])
	}

	// Example 1: Assign a single device to server
	fmt.Println("\n=== Example 1: Assign Single Device to Server ===")

	if len(unassignedDevices) > 0 {
		singleDevice := unassignedDevices[0]
		fmt.Printf("Assigning device %s (Serial: %s) to server %s...\n",
			singleDevice.ID, singleDevice.Attributes.SerialNumber, targetServerName)

		singleDeviceIDs := []string{singleDevice.ID}

		assignResponse, err := client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, singleDeviceIDs)
		if err != nil {
			log.Printf("Error assigning single device: %v", err)
		} else {
			fmt.Printf("Assignment successful!\n")
			fmt.Printf("  Activity ID: %s\n", assignResponse.Data.ID)
			fmt.Printf("  Activity Type: %s\n", assignResponse.Data.Type)

			if assignResponse.Data.Attributes != nil {
				fmt.Printf("  Activity: %s\n", assignResponse.Data.Attributes.ActivityType)
				if assignResponse.Data.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", assignResponse.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
			}

			if assignResponse.Links != nil {
				fmt.Printf("  Self Link: %s\n", assignResponse.Links.Self)
			}
		}
	}

	// Example 2: Assign multiple devices to server
	fmt.Println("\n=== Example 2: Assign Multiple Devices to Server ===")

	maxDevicesToAssign := min(3, len(unassignedDevices))
	if maxDevicesToAssign > 1 {
		var multipleDeviceIDs []string
		fmt.Printf("Assigning %d devices to server %s:\n", maxDevicesToAssign, targetServerName)

		for i := 0; i < maxDevicesToAssign; i++ {
			device := unassignedDevices[i]
			multipleDeviceIDs = append(multipleDeviceIDs, device.ID)
			fmt.Printf("  - %s (Serial: %s)\n", device.ID, device.Attributes.SerialNumber)
		}

		multipleAssignResponse, err := client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, multipleDeviceIDs)
		if err != nil {
			log.Printf("Error assigning multiple devices: %v", err)
		} else {
			fmt.Printf("Multiple device assignment successful!\n")
			fmt.Printf("  Activity ID: %s\n", multipleAssignResponse.Data.ID)

			if multipleAssignResponse.Data.Attributes != nil {
				fmt.Printf("  Activity Type: %s\n", multipleAssignResponse.Data.Attributes.ActivityType)
				if multipleAssignResponse.Data.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", multipleAssignResponse.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
			}
		}
	} else {
		fmt.Println("Not enough unassigned devices for multiple assignment demo")
	}

	// Example 3: Verify assignment
	fmt.Println("\n=== Example 3: Verify Device Assignment ===")

	if len(unassignedDevices) > 0 {
		deviceToCheck := unassignedDevices[0]
		fmt.Printf("Verifying assignment for device %s (Serial: %s)...\n",
			deviceToCheck.ID, deviceToCheck.Attributes.SerialNumber)

		// Wait a moment for the assignment to process
		time.Sleep(2 * time.Second)

		serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, deviceToCheck.ID)
		if err != nil {
			log.Printf("Error checking device assignment: %v", err)
		} else {
			if serverLinkage.Data.ID == targetMDMServerID {
				fmt.Printf("✓ Device successfully assigned to server %s\n", targetMDMServerID)
			} else if serverLinkage.Data.ID != "" {
				fmt.Printf("⚠ Device assigned to different server: %s\n", serverLinkage.Data.ID)
			} else {
				fmt.Printf("✗ Device appears to still be unassigned\n")
			}
		}
	}

	// Example 4: Error handling - invalid MDM server ID
	fmt.Println("\n=== Example 4: Error Handling (Invalid MDM Server ID) ===")

	if len(unassignedDevices) > 0 {
		invalidServerID := "invalid-mdm-server-id"
		testDeviceIDs := []string{unassignedDevices[0].ID}

		_, err = client.DeviceManagement.AssignDevicesToServer(ctx, invalidServerID, testDeviceIDs)
		if err != nil {
			fmt.Printf("Expected error for invalid MDM server ID '%s': %v\n", invalidServerID, err)
		}
	}

	// Example 5: Error handling - empty MDM server ID
	fmt.Println("\n=== Example 5: Error Handling (Empty MDM Server ID) ===")

	if len(unassignedDevices) > 0 {
		testDeviceIDs := []string{unassignedDevices[0].ID}

		_, err = client.DeviceManagement.AssignDevicesToServer(ctx, "", testDeviceIDs)
		if err != nil {
			fmt.Printf("Expected error for empty MDM server ID: %v\n", err)
		}
	}

	// Example 6: Error handling - empty device IDs
	fmt.Println("\n=== Example 6: Error Handling (Empty Device IDs) ===")

	emptyDeviceIDs := []string{}
	_, err = client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, emptyDeviceIDs)
	if err != nil {
		fmt.Printf("Expected error for empty device IDs: %v\n", err)
	}

	// Example 7: Error handling - invalid device IDs
	fmt.Println("\n=== Example 7: Error Handling (Invalid Device IDs) ===")

	invalidDeviceIDs := []string{"invalid-device-id-1", "invalid-device-id-2"}
	_, err = client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, invalidDeviceIDs)
	if err != nil {
		fmt.Printf("Expected error for invalid device IDs: %v\n", err)
	}

	// Example 8: Assignment status tracking
	fmt.Println("\n=== Example 8: Assignment Status Tracking ===")

	if len(unassignedDevices) > 0 {
		trackingDevice := unassignedDevices[0]
		fmt.Printf("Tracking assignment status for device %s...\n", trackingDevice.Attributes.SerialNumber)

		// Check initial status
		fmt.Printf("Initial status check...\n")
		initialLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, trackingDevice.ID)
		if err != nil {
			fmt.Printf("  Error checking initial status: %v\n", err)
		} else {
			fmt.Printf("  Initial assignment: %s\n", initialLinkage.Data.ID)
		}

		// Perform assignment
		fmt.Printf("Performing assignment...\n")
		trackingDeviceIDs := []string{trackingDevice.ID}
		assignResponse, err := client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, trackingDeviceIDs)
		if err != nil {
			fmt.Printf("  Assignment error: %v\n", err)
		} else {
			fmt.Printf("  Assignment activity created: %s\n", assignResponse.Data.ID)
		}

		// Check final status
		fmt.Printf("Final status check (after 3 seconds)...\n")
		time.Sleep(3 * time.Second)
		finalLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, trackingDevice.ID)
		if err != nil {
			fmt.Printf("  Error checking final status: %v\n", err)
		} else {
			fmt.Printf("  Final assignment: %s\n", finalLinkage.Data.ID)
			if finalLinkage.Data.ID == targetMDMServerID {
				fmt.Printf("  ✓ Assignment successful!\n")
			} else {
				fmt.Printf("  ⚠ Assignment may still be processing\n")
			}
		}
	}

	// Example 9: Pretty print JSON response
	fmt.Println("\n=== Example 9: Full JSON Response ===")
	if len(unassignedDevices) > 0 {
		// Perform one more assignment for JSON demo
		jsonDemoDeviceIDs := []string{unassignedDevices[0].ID}
		jsonResponse, err := client.DeviceManagement.AssignDevicesToServer(ctx, targetMDMServerID, jsonDemoDeviceIDs)
		if err != nil {
			log.Printf("Error in JSON demo assignment: %v", err)
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
	fmt.Println("Note: Device assignments may take some time to process in Apple's system.")
	fmt.Println("Check the Apple Business Manager portal to verify final assignment status.")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
