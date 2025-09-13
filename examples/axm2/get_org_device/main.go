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

	// Example 1: Get device with all fields
	log.Println("=== Get Organization Device (All Fields) ===")

	device, err := client.OrgDevices().GetOrgDevice(ctx, exampleDeviceID)
	if err != nil {
		log.Fatalf("Failed to get organization device: %v", err)
	}

	log.Printf("Device Details:")
	log.Printf("  ID: %s", device.ID)
	log.Printf("  Type: %s", device.Type)
	log.Printf("  Serial Number: %s", device.Attributes.SerialNumber)
	log.Printf("  Device Model: %s", device.Attributes.DeviceModel)
	log.Printf("  Product Family: %s", device.Attributes.ProductFamily)
	log.Printf("  Product Type: %s", device.Attributes.ProductType)
	log.Printf("  Device Capacity: %s", device.Attributes.DeviceCapacity)
	log.Printf("  Part Number: %s", device.Attributes.PartNumber)
	log.Printf("  Order Number: %s", device.Attributes.OrderNumber)
	log.Printf("  Color: %s", device.Attributes.Color)
	log.Printf("  Status: %s", device.Attributes.Status)

	if device.Attributes.AddedToOrgDateTime != "" {
		log.Printf("  Added to Org: %s", device.Attributes.AddedToOrgDateTime)
	}
	if device.Attributes.UpdatedDateTime != "" {
		log.Printf("  Updated: %s", device.Attributes.UpdatedDateTime)
	}
	if device.Attributes.OrderDateTime != "" {
		log.Printf("  Order Date: %s", device.Attributes.OrderDateTime)
	}

	// Network identifiers
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

	// Purchase information
	if device.Attributes.PurchaseSourceUid != "" {
		log.Printf("  Purchase Source UID: %s", device.Attributes.PurchaseSourceUid)
	}
	if device.Attributes.PurchaseSourceType != "" {
		log.Printf("  Purchase Source Type: %s", device.Attributes.PurchaseSourceType)
	}

	// Example 2: Get device with specific fields only
	log.Println("\n=== Get Device with Field Filtering ===")

	filteredDevice, err := client.OrgDevices().GetOrgDevice(ctx, exampleDeviceID,
		axm2.WithFieldsOption("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"status",
			"productFamily",
			"addedToOrgDateTime",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get filtered device: %v", err)
	}

	log.Printf("Filtered Device Details:")
	log.Printf("  Serial Number: %s", filteredDevice.Attributes.SerialNumber)
	log.Printf("  Device Model: %s", filteredDevice.Attributes.DeviceModel)
	log.Printf("  Status: %s", filteredDevice.Attributes.Status)
	log.Printf("  Product Family: %s", filteredDevice.Attributes.ProductFamily)
	log.Printf("  Added to Org: %s", filteredDevice.Attributes.AddedToOrgDateTime)

	// Relationship information
	if device.Relationships.AssignedServer != nil && len(device.Relationships.AssignedServer.Links) > 0 {
		log.Printf("\n=== Device Relationships ===")
		log.Printf("Device has assigned server relationship:")
		for key, url := range device.Relationships.AssignedServer.Links {
			log.Printf("  %s: %s", key, url)
		}
	}

	// Self link
	if device.Links.Self != "" {
		log.Printf("\n=== Resource Links ===")
		log.Printf("  Self Link: %s", device.Links.Self)
	}
}
