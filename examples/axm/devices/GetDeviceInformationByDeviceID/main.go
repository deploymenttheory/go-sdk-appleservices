package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Device Information Example ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	// Parse the private key (supports both RSA and ECDSA)
	privateKey, err := client.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// First, get a list of devices to find a device ID to query
	fmt.Println("\nStep 1: Getting organization devices to find a device ID...")

	listOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
		Limit: 5,
	}

	devicesResponse, err := client.Devices.GetOrganizationDevices(ctx, listOptions)
	if err != nil {
		log.Fatalf("Error getting organization devices: %v", err)
	}

	if len(devicesResponse.Data) == 0 {
		log.Fatalf("No devices found in organization")
	}

	// Use the first device for detailed information
	deviceID := devicesResponse.Data[0].ID
	fmt.Printf("Found device ID: %s (Serial: %s)\n", deviceID, devicesResponse.Data[0].Attributes.SerialNumber)

	// Step 2: Get detailed information for the specific device
	fmt.Printf("\nStep 2: Getting detailed information for device %s...\n", deviceID)

	// Example 1: Get all available fields
	fmt.Println("\n=== Example 1: Get All Available Device Information ===")

	allFieldsOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldProductFamily,
			devices.FieldStatus,
			devices.FieldColor,
			devices.FieldAddedToOrgDateTime,
			devices.FieldUpdatedDateTime,
			devices.FieldIMEI,
			devices.FieldMEID,
			devices.FieldPurchaseSourceUid,
		},
	}

	deviceInfo, err := client.Devices.GetDeviceInformationByDeviceID(ctx, deviceID, allFieldsOptions)
	if err != nil {
		log.Fatalf("Error getting device information: %v", err)
	}

	// Display detailed device information
	device := deviceInfo.Data
	fmt.Printf("Device Information:\n")
	fmt.Printf("  ID: %s\n", device.ID)
	fmt.Printf("  Type: %s\n", device.Type)
	fmt.Printf("  Serial Number: %s\n", device.Attributes.SerialNumber)
	fmt.Printf("  Device Model: %s\n", device.Attributes.DeviceModel)
	fmt.Printf("  Product Family: %s\n", device.Attributes.ProductFamily)
	fmt.Printf("  Status: %s\n", device.Attributes.Status)
	fmt.Printf("  Color: %s\n", device.Attributes.Color)

	if device.Attributes.AddedToOrgDateTime != nil {
		fmt.Printf("  Added to Org: %s\n", device.Attributes.AddedToOrgDateTime.Format(time.RFC3339))
	}

	if device.Attributes.UpdatedDateTime != nil {
		fmt.Printf("  Last Updated: %s\n", device.Attributes.UpdatedDateTime.Format(time.RFC3339))
	}

	if len(device.Attributes.IMEI) > 0 {
		fmt.Printf("  IMEI: %v\n", device.Attributes.IMEI)
	}

	if len(device.Attributes.MEID) > 0 {
		fmt.Printf("  MEID: %v\n", device.Attributes.MEID)
	}

	if device.Attributes.PurchaseSourceUid != "" {
		fmt.Printf("  Purchase Source UID: %s\n", device.Attributes.PurchaseSourceUid)
	}

	// Example 2: Get only specific fields
	fmt.Println("\n=== Example 2: Get Only Specific Fields ===")

	specificFieldsOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
	}

	specificDeviceInfo, err := client.Devices.GetDeviceInformationByDeviceID(ctx, deviceID, specificFieldsOptions)
	if err != nil {
		log.Printf("Error getting specific device information: %v", err)
	} else {
		specificDevice := specificDeviceInfo.Data
		fmt.Printf("Specific Fields Only:\n")
		fmt.Printf("  Serial: %s\n", specificDevice.Attributes.SerialNumber)
		fmt.Printf("  Model: %s\n", specificDevice.Attributes.DeviceModel)
		fmt.Printf("  Status: %s\n", specificDevice.Attributes.Status)
	}

	// Example 3: Get device information with no field filtering (all fields)
	fmt.Println("\n=== Example 3: Get Device Information (No Field Filtering) ===")

	noFilterDeviceInfo, err := client.Devices.GetDeviceInformationByDeviceID(ctx, deviceID, nil)
	if err != nil {
		log.Printf("Error getting unfiltered device information: %v", err)
	} else {
		fmt.Printf("Unfiltered device information retrieved successfully\n")
		fmt.Printf("Device ID: %s\n", noFilterDeviceInfo.Data.ID)
	}

	// Example 4: Try to get information for a non-existent device (error handling)
	fmt.Println("\n=== Example 4: Error Handling (Non-existent Device) ===")

	fakeDeviceID := "non-existent-device-id"
	_, err = client.Devices.GetDeviceInformationByDeviceID(ctx, fakeDeviceID, nil)
	if err != nil {
		fmt.Printf("Expected error for non-existent device: %v\n", err)
	}

	// Example 5: Pretty print full JSON response
	fmt.Println("\n=== Example 5: Full JSON Response ===")
	jsonData, err := json.MarshalIndent(deviceInfo, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	// Example 6: Multiple devices information
	fmt.Println("\n=== Example 6: Get Information for Multiple Devices ===")

	if len(devicesResponse.Data) > 1 {
		fmt.Printf("Getting information for first %d devices:\n", min(3, len(devicesResponse.Data)))

		for i, dev := range devicesResponse.Data[:min(3, len(devicesResponse.Data))] {
			fmt.Printf("\nDevice %d (ID: %s):\n", i+1, dev.ID)

			info, err := client.Devices.GetDeviceInformationByDeviceID(ctx, dev.ID, &devices.RequestQueryOptions{
				Fields: []string{
					devices.FieldSerialNumber,
					devices.FieldDeviceModel,
					devices.FieldStatus,
					devices.FieldColor,
				},
			})

			if err != nil {
				fmt.Printf("  Error: %v\n", err)
				continue
			}

			d := info.Data
			fmt.Printf("  Serial: %s\n", d.Attributes.SerialNumber)
			fmt.Printf("  Model: %s\n", d.Attributes.DeviceModel)
			fmt.Printf("  Status: %s\n", d.Attributes.Status)
			fmt.Printf("  Color: %s\n", d.Attributes.Color)
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
