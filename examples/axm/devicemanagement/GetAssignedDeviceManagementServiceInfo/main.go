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
	fmt.Println("=== Apple Business Manager - Get Assigned Device Management Service Information Example ===")

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

	devicesResponse, err := client.Devices.GetOrganizationDevices(ctx, &devices.RequestQueryOptions{
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

	// Step 2: Find a device that has an assigned server
	fmt.Println("\nStep 2: Finding a device with an assigned server...")

	var assignedDeviceID string
	var assignedDeviceSerial string

	for _, device := range devicesResponse.Data {
		// Check if this device has an assigned server
		serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, device.ID)
		if err != nil {
			continue // Skip devices with errors
		}

		if serverLinkage.Data.ID != "" {
			assignedDeviceID = device.ID
			assignedDeviceSerial = device.Attributes.SerialNumber
			fmt.Printf("Found device with assigned server: %s (Serial: %s, Server ID: %s)\n",
				assignedDeviceID, assignedDeviceSerial, serverLinkage.Data.ID)
			break
		}
	}

	if assignedDeviceID == "" {
		fmt.Println("No devices with assigned servers found. Using first device for demonstration...")
		assignedDeviceID = devicesResponse.Data[0].ID
		assignedDeviceSerial = devicesResponse.Data[0].Attributes.SerialNumber
	}

	// Example 1: Get assigned server information with all fields
	fmt.Println("\n=== Example 1: Get Assigned Server Information (All Fields) ===")

	fmt.Printf("Getting server information for device: %s (Serial: %s)\n", assignedDeviceID, assignedDeviceSerial)

	serverInfo, err := client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, assignedDeviceID, nil)
	if err != nil {
		log.Printf("Error getting assigned server information: %v", err)
	} else {
		fmt.Printf("Assigned Server Information:\n")
		fmt.Printf("  Server ID: %s\n", serverInfo.Data.ID)
		fmt.Printf("  Type: %s\n", serverInfo.Data.Type)

		if serverInfo.Data.Attributes != nil {
			fmt.Printf("  Server Name: %s\n", serverInfo.Data.Attributes.ServerName)
			fmt.Printf("  Server Type: %s\n", serverInfo.Data.Attributes.ServerType)

			if serverInfo.Data.Attributes.CreatedDateTime != nil {
				fmt.Printf("  Created: %s\n", serverInfo.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
			}

			if serverInfo.Data.Attributes.UpdatedDateTime != nil {
				fmt.Printf("  Updated: %s\n", serverInfo.Data.Attributes.UpdatedDateTime.Format(time.RFC3339))
			}
		}

		if serverInfo.Data.Relationships != nil && serverInfo.Data.Relationships.Devices != nil {
			fmt.Printf("  Devices Relationship: %+v\n", serverInfo.Data.Relationships.Devices)
		}

		if serverInfo.Links != nil {
			fmt.Printf("  Self Link: %s\n", serverInfo.Links.Self)
		}
	}

	// Example 2: Get assigned server information with specific fields
	fmt.Println("\n=== Example 2: Get Assigned Server Information (Specific Fields) ===")

	specificFieldsOptions := &devicemanagement.GetAssignedServerInfoOptions{
		Fields: []string{
			"serverName",
			"serverType",
			"createdDateTime",
		},
	}

	specificServerInfo, err := client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, assignedDeviceID, specificFieldsOptions)
	if err != nil {
		log.Printf("Error getting specific server information: %v", err)
	} else {
		fmt.Printf("Specific Server Information:\n")
		fmt.Printf("  Server ID: %s\n", specificServerInfo.Data.ID)

		if specificServerInfo.Data.Attributes != nil {
			fmt.Printf("  Server Name: %s\n", specificServerInfo.Data.Attributes.ServerName)
			fmt.Printf("  Server Type: %s\n", specificServerInfo.Data.Attributes.ServerType)

			if specificServerInfo.Data.Attributes.CreatedDateTime != nil {
				fmt.Printf("  Created: %s\n", specificServerInfo.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
			}
		}
	}

	// Example 3: Get server information for multiple devices
	fmt.Println("\n=== Example 3: Get Server Information for Multiple Devices ===")

	maxDevicesToCheck := min(5, len(devicesResponse.Data))
	fmt.Printf("Checking server information for first %d devices:\n", maxDevicesToCheck)

	for i := 0; i < maxDevicesToCheck; i++ {
		device := devicesResponse.Data[i]
		fmt.Printf("\nDevice %d: %s (Serial: %s)\n", i+1, device.ID, device.Attributes.SerialNumber)

		serverInfo, err := client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, device.ID, &devicemanagement.GetAssignedServerInfoOptions{
			Fields: []string{
				"serverName",
				"serverType",
			},
		})

		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		if serverInfo.Data.Attributes != nil {
			fmt.Printf("  ✓ Assigned Server: %s (Type: %s)\n",
				serverInfo.Data.Attributes.ServerName,
				serverInfo.Data.Attributes.ServerType)
		} else {
			fmt.Printf("  ✗ No server information available\n")
		}
	}

	// Example 4: Get server information with all available fields
	fmt.Println("\n=== Example 4: Get Server Information (All Available Fields) ===")

	allFieldsOptions := &devicemanagement.GetAssignedServerInfoOptions{
		Fields: []string{
			"serverName",
			"serverType",
			"createdDateTime",
			"updatedDateTime",
			"devices",
		},
	}

	allFieldsServerInfo, err := client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, assignedDeviceID, allFieldsOptions)
	if err != nil {
		log.Printf("Error getting complete server information: %v", err)
	} else {
		fmt.Printf("Complete Server Information:\n")
		fmt.Printf("  Server ID: %s\n", allFieldsServerInfo.Data.ID)
		fmt.Printf("  Type: %s\n", allFieldsServerInfo.Data.Type)

		if allFieldsServerInfo.Data.Attributes != nil {
			attrs := allFieldsServerInfo.Data.Attributes
			fmt.Printf("  Server Name: %s\n", attrs.ServerName)
			fmt.Printf("  Server Type: %s\n", attrs.ServerType)

			if attrs.CreatedDateTime != nil {
				fmt.Printf("  Created: %s\n", attrs.CreatedDateTime.Format(time.RFC3339))
			}

			if attrs.UpdatedDateTime != nil {
				fmt.Printf("  Updated: %s\n", attrs.UpdatedDateTime.Format(time.RFC3339))
			}
		}

		if allFieldsServerInfo.Data.Relationships != nil {
			fmt.Printf("  Relationships: %+v\n", allFieldsServerInfo.Data.Relationships)
		}
	}

	// Example 5: Error handling - invalid device ID
	fmt.Println("\n=== Example 5: Error Handling (Invalid Device ID) ===")

	invalidDeviceID := "invalid-device-id-12345"
	_, err = client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, invalidDeviceID, nil)
	if err != nil {
		fmt.Printf("Expected error for invalid device ID '%s': %v\n", invalidDeviceID, err)
	}

	// Example 6: Error handling - empty device ID
	fmt.Println("\n=== Example 6: Error Handling (Empty Device ID) ===")

	_, err = client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, "", nil)
	if err != nil {
		fmt.Printf("Expected error for empty device ID: %v\n", err)
	}

	// Example 7: Compare server linkage vs server information
	fmt.Println("\n=== Example 7: Compare Server Linkage vs Server Information ===")

	fmt.Printf("Comparing linkage and information for device: %s\n", assignedDeviceSerial)

	// Get server linkage (just the ID)
	serverLinkage, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, assignedDeviceID)
	if err != nil {
		log.Printf("Error getting server linkage: %v", err)
	} else {
		fmt.Printf("Server Linkage - Server ID: %s\n", serverLinkage.Data.ID)
	}

	// Get server information (full details)
	serverInfo, err = client.DeviceManagement.GetAssignedDeviceManagementServiceInformationByDeviceID(ctx, assignedDeviceID, nil)
	if err != nil {
		log.Printf("Error getting server information: %v", err)
	} else {
		fmt.Printf("Server Information - Server ID: %s", serverInfo.Data.ID)
		if serverInfo.Data.Attributes != nil {
			fmt.Printf(" (Name: %s)", serverInfo.Data.Attributes.ServerName)
		}
		fmt.Println()
	}

	// Example 8: Pretty print JSON response
	fmt.Println("\n=== Example 8: Full JSON Response ===")
	if serverInfo != nil {
		jsonData, err := json.MarshalIndent(serverInfo, "", "  ")
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
