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

	// First, get a list of devices to find a device ID
	log.Println("\n=== Finding a Device ID ===")
	devices, err := client.GetOrgDevices(ctx, axm2.WithLimitOption(1))
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	if len(devices) == 0 {
		log.Fatalf("No devices found in organization")
	}

	deviceID := devices[0].ID
	log.Printf("Using device ID: %s", deviceID)

	// Example: Get a specific organization device
	log.Println("\n=== Get Specific Organization Device ===")

	// Method 1: Using RequestOption pattern (recommended)
	device, err := client.GetOrgDevice(ctx, deviceID,
		axm2.WithFieldsOption("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"productFamily",
			"status",
			"addedToOrgDateTime",
			"deviceCapacity",
			"color",
			"partNumber",
			"orderNumber",
			"imei",
			"meid",
			"eid",
			"wifiMacAddress",
			"bluetoothMacAddress",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get device: %v", err)
	}

	// Display detailed device information
	log.Printf("Device Details:")
	log.Printf("  ID: %s", device.ID)
	log.Printf("  Type: %s", device.Type)
	log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
	log.Printf("  Device Model: %s", device.Attributes.DeviceModel)
	log.Printf("  Product Family: %s", device.Attributes.ProductFamily)
	log.Printf("  Product Type: %s", device.Attributes.ProductType)
	log.Printf("  Status: %s", device.Attributes.Status)
	log.Printf("  Device Capacity: %s", device.Attributes.DeviceCapacity)
	log.Printf("  Color: %s", device.Attributes.Color)
	log.Printf("  Part Number: %s", device.Attributes.PartNumber)
	log.Printf("  Order Number: %s", device.Attributes.OrderNumber)
	log.Printf("  Added to Org: %s", device.Attributes.AddedToOrgDateTime)
	log.Printf("  Updated: %s", device.Attributes.UpdatedDateTime)
	log.Printf("  Order Date: %s", device.Attributes.OrderDateTime)

	if len(device.Attributes.IMEI) > 0 {
		log.Printf("  IMEI: %v", device.Attributes.IMEI)
	}
	if len(device.Attributes.MEID) > 0 {
		log.Printf("  MEID: %v", device.Attributes.MEID)
	}
	if device.Attributes.EID != "" {
		log.Printf("  EID: %s", device.Attributes.EID)
	}
	if device.Attributes.WifiMacAddress != "" {
		log.Printf("  WiFi MAC: %s", device.Attributes.WifiMacAddress)
	}
	if device.Attributes.BluetoothMacAddress != "" {
		log.Printf("  Bluetooth MAC: %s", device.Attributes.BluetoothMacAddress)
	}
	if device.Attributes.PurchaseSourceUid != "" {
		log.Printf("  Purchase Source UID: %s", device.Attributes.PurchaseSourceUid)
	}
	if device.Attributes.PurchaseSourceType != "" {
		log.Printf("  Purchase Source Type: %s", device.Attributes.PurchaseSourceType)
	}

	// Method 2: Using legacy QueryBuilder (for backward compatibility)
	log.Println("\n=== Using Legacy QueryBuilder ===")
	queryBuilder := client.NewQueryBuilder().
		Fields("orgDevices", []string{"serialNumber", "deviceModel", "status"})

	legacyDevice, err := client.GetOrgDeviceWithQuery(ctx, deviceID, queryBuilder)
	if err != nil {
		log.Fatalf("Failed to get device with QueryBuilder: %v", err)
	}

	log.Printf("Legacy Device: %s (%s) - %s",
		legacyDevice.Attributes.SerialNumber,
		legacyDevice.Attributes.DeviceModel,
		legacyDevice.Attributes.Status)

	log.Println("\n=== Example completed successfully ===")
}
