package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Assigned Device Management Service ID Example ===")

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

	// Step 1: Get organization devices to find device IDs
	fmt.Println("\nStep 1: Getting organization devices to find device IDs...")

	devicesResponse, err := client.
		Devices.
		GetOrganizationDevicesV1(ctx, &devices.RequestQueryOptions{
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

	// Example 1: Get assigned server ID for the first device
	fmt.Println("\n=== Example 1: Get Assigned Server ID for First Device ===")

	firstDevice := devicesResponse.Data[0]
	deviceID := firstDevice.ID

	fmt.Printf("Checking device: %s (Serial: %s, Model: %s)\n",
		deviceID, firstDevice.Attributes.SerialNumber, firstDevice.Attributes.DeviceModel)

	assignedServerResponse, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, deviceID)
	if err != nil {
		log.Printf("Error getting assigned server ID for device %s: %v", deviceID, err)
	} else {
		fmt.Printf("Assigned Server Linkage:\n")
		fmt.Printf("  Type: %s\n", assignedServerResponse.Data.Type)
		fmt.Printf("  Server ID: %s\n", assignedServerResponse.Data.ID)

		if assignedServerResponse.Links != nil {
			fmt.Printf("  Self Link: %s\n", assignedServerResponse.Links.Self)
		}
	}

	// Example 2: Check multiple devices for assigned servers
	fmt.Println("\n=== Example 2: Check Multiple Devices for Assigned Servers ===")

	maxDevicesToCheck := min(5, len(devicesResponse.Data))
	fmt.Printf("Checking assigned servers for first %d devices:\n", maxDevicesToCheck)

	for i := 0; i < maxDevicesToCheck; i++ {
		device := devicesResponse.Data[i]
		fmt.Printf("\nDevice %d: %s (Serial: %s)\n", i+1, device.ID, device.Attributes.SerialNumber)

		serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, device.ID)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		if serverLinkage.Data.ID != "" {
			fmt.Printf("  ✓ Assigned to Server ID: %s\n", serverLinkage.Data.ID)
		} else {
			fmt.Printf("  ✗ No server assigned\n")
		}
	}

	// Example 3: Error handling - invalid device ID
	fmt.Println("\n=== Example 3: Error Handling (Invalid Device ID) ===")

	invalidDeviceID := "invalid-device-id-12345"
	_, err = client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, invalidDeviceID)
	if err != nil {
		fmt.Printf("Expected error for invalid device ID '%s': %v\n", invalidDeviceID, err)
	}

	// Example 4: Error handling - empty device ID
	fmt.Println("\n=== Example 4: Error Handling (Empty Device ID) ===")

	_, err = client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, "")
	if err != nil {
		fmt.Printf("Expected error for empty device ID: %v\n", err)
	}

	// Example 5: Check devices that might not have assigned servers
	fmt.Println("\n=== Example 5: Comprehensive Device Assignment Check ===")

	assignedCount := 0
	unassignedCount := 0
	errorCount := 0

	fmt.Printf("Checking all %d devices for server assignments...\n", len(devicesResponse.Data))

	for i, device := range devicesResponse.Data {
		serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, device.ID)
		if err != nil {
			errorCount++
			fmt.Printf("  Device %d (%s): ERROR - %v\n", i+1, device.Attributes.SerialNumber, err)
			continue
		}

		if serverLinkage.Data.ID != "" {
			assignedCount++
			fmt.Printf("  Device %d (%s): ASSIGNED to %s\n", i+1, device.Attributes.SerialNumber, serverLinkage.Data.ID)
		} else {
			unassignedCount++
			fmt.Printf("  Device %d (%s): UNASSIGNED\n", i+1, device.Attributes.SerialNumber)
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Assigned devices: %d\n", assignedCount)
	fmt.Printf("  Unassigned devices: %d\n", unassignedCount)
	fmt.Printf("  Errors: %d\n", errorCount)
	fmt.Printf("  Total checked: %d\n", len(devicesResponse.Data))

	// Example 6: Get detailed information about assigned servers
	fmt.Println("\n=== Example 6: Get Server Details for Assigned Devices ===")

	// Get MDM servers for reference
	serversResponse, err := client.DeviceManagement.GetDeviceManagementServicesV1(ctx, nil)
	if err != nil {
		log.Printf("Error getting MDM servers: %v", err)
	} else {
		// Create a map of server ID to server name for quick lookup
		serverNames := make(map[string]string)
		for _, server := range serversResponse.Data {
			if server.Attributes != nil {
				serverNames[server.ID] = server.Attributes.ServerName
			}
		}

		fmt.Printf("Checking first 3 devices with server name resolution:\n")
		for i := 0; i < min(3, len(devicesResponse.Data)); i++ {
			device := devicesResponse.Data[i]

			serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADeviceV1(ctx, device.ID)
			if err != nil {
				fmt.Printf("  Device %s: Error - %v\n", device.Attributes.SerialNumber, err)
				continue
			}

			serverName := "Unknown"
			if name, exists := serverNames[serverLinkage.Data.ID]; exists {
				serverName = name
			}

			fmt.Printf("  Device %s: Assigned to '%s' (ID: %s)\n",
				device.Attributes.SerialNumber, serverName, serverLinkage.Data.ID)
		}
	}

	// Example 7: Pretty print JSON response
	fmt.Println("\n=== Example 7: Full JSON Response ===")
	if assignedServerResponse != nil {
		jsonData, err := json.MarshalIndent(assignedServerResponse, "", "  ")
		if err != nil {
			log.Printf("Error marshaling response to JSON: %v", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}

	fmt.Println("\n=== Example Complete ===")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
