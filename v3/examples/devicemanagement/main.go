package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	axm "github.com/deploymenttheory/go-api-sdk-apple/v3"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager V3 - Direct Service Access Example ===")

	// Use credentials directly for testing
	keyID := "bd4bd60b-6ddf-4fee-1111-3ed8f6dc4bd1"
	issuerID := "BUSINESSAPI.3bb3a62b-6f21-4802-1111-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgSVST2uwXoc9Gc87H
uqq7jYhn+PlsrtxPQebp0LeDXp2hRANCAASBtSEWU1075awq69clg4ZPNdPiAX77
mdH5iVYM8fVK1mAAk1ewo3YWlhz2GEGuox04Ng5xVrpotMQXo2WQ1111
-----END EC PRIVATE KEY-----`

	// Parse the private key
	privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Create AXM client with embedded services - NO WRAPPERS!
	builder := axm.NewClientWithBuilder().
		WithJWTAuth(keyID, issuerID, privateKey).
		WithDebug(true)

	client, err := axm.BuildClient(builder)

	if err != nil {
		log.Fatalf("Failed to create AXM client: %v", err)
	}

	// Create context
	ctx := context.Background()

	// V3 PATTERN: Direct access via client.Service.Method()
	fmt.Println("\n=== V3 Pattern: client.Devices.GetOrganizationDevices() ===")

	devicesResponse, err := client.Devices.GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
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

	// V3 PATTERN: Direct access to device management service
	fmt.Println("\n=== V3 Pattern: client.DeviceManagement.GetDeviceManagementServices() ===")

	serversResponse, err := client.DeviceManagement.GetDeviceManagementServices(ctx, nil)
	if err != nil {
		log.Printf("Error getting MDM servers: %v", err)
	} else {
		fmt.Printf("Found %d MDM servers\n", len(serversResponse.Data))
	}

	// V3 PATTERN: Chaining service calls with direct access
	fmt.Println("\n=== V3 Pattern: Direct service method chaining ===")

	firstDevice := devicesResponse.Data[0]
	deviceID := firstDevice.ID

	fmt.Printf("Checking device: %s (Serial: %s, Model: %s)\n",
		deviceID, firstDevice.Attributes.SerialNumber, firstDevice.Attributes.DeviceModel)

	// Direct service access - no wrapper clients needed!
	assignedServerResponse, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, deviceID)
	if err != nil {
		log.Printf("Error getting assigned server ID for device %s: %v", deviceID, err)
	} else {
		fmt.Printf("Assigned Server Linkage:\n")
		fmt.Printf("  Type: %s\n", assignedServerResponse.Data.Type)
		fmt.Printf("  Server ID: %s\n", assignedServerResponse.Data.ID)
	}

	// V3 PATTERN: Multiple service access in same flow
	fmt.Println("\n=== V3 Pattern: Mixed service operations ===")

	// Get device details using devices service
	deviceDetails, err := client.Devices.GetDeviceInformationByDeviceID(ctx, deviceID, &devices.GetDeviceInformationOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
			devices.FieldAssignedServer,
		},
	})
	if err != nil {
		log.Printf("Error getting device details: %v", err)
	} else {
		fmt.Printf("Device Details:\n")
		fmt.Printf("  Serial: %s\n", deviceDetails.Data.Attributes.SerialNumber)
		fmt.Printf("  Model: %s\n", deviceDetails.Data.Attributes.DeviceModel)
		fmt.Printf("  Status: %s\n", deviceDetails.Data.Attributes.Status)
	}

	// Get MDM server device linkages using device management service
	if len(serversResponse.Data) > 0 {
		firstServer := serversResponse.Data[0]
		linkages, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, firstServer.ID, nil)
		if err != nil {
			log.Printf("Error getting server linkages: %v", err)
		} else {
			fmt.Printf("Server '%s' has %d linked devices\n", firstServer.Attributes.ServerName, len(linkages.Data))
		}
	}

	// Pretty print JSON response
	fmt.Println("\n=== V3 Architecture Benefits ===")
	fmt.Println("✓ No wrapper clients needed")
	fmt.Println("✓ Direct access: client.Service.Method()")
	fmt.Println("✓ Type-safe service interfaces")
	fmt.Println("✓ Preserved functionality from original CRUD methods")
	fmt.Println("✓ Integrated pagination support")

	jsonData, err := json.MarshalIndent(devicesResponse, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Printf("\nSample JSON Response:\n%s\n", string(jsonData)[:500])
		fmt.Println("...")
	}

	fmt.Println("\n=== V3 Example Complete ===")
}
