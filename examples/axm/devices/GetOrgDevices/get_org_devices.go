package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
	axm "github.com/deploymenttheory/go-api-sdk-apple/services/axm"
)

func main() {
	// QuickStart automatically handles configuration from:
	// 1. Environment variables (APPLE_ORG_ID, APPLE_KEY_ID, APPLE_PRIVATE_KEY/PATH)
	// 2. config.json file (if it exists)
	// 3. Default configuration values

	fmt.Println("=== Apple School and Business Manager API - Device Examples ===")
	fmt.Println("Using QuickStart configuration...")

	axmClient, err := client.QuickStart("/Users/dafyddwatkins/go-api-sdk-apple/examples/axm/devices/GetOrgDevices/config.example.json")
	if err != nil {
		log.Fatal("Configuration failed. Please either:\n" +
			"  1. Set environment variables: APPLE_ORG_ID, APPLE_KEY_ID, APPLE_PRIVATE_KEY\n" +
			"  2. Create a config.json file based on config.example.json\n" +
			"Error: " + err.Error())
	}
	defer axmClient.Close()

	// Create service wrapper
	axmService := axm.NewClient(axmClient)

	fmt.Println("✓ Successfully authenticated with Apple API")
	fmt.Println()

	// Example 1: Get all devices with automatic pagination
	fmt.Println("=== Example 1: Get All Organization Devices ===")
	getAllDevicesExample(axmService)

	// Example 2: Get devices with field filtering
	fmt.Println("\n=== Example 2: Get Devices with Field Filtering ===")
	getDevicesWithFieldsExample(axmService)

	// Example 3: Get devices with status filtering
	fmt.Println("\n=== Example 3: Get Devices with Status Filtering ===")
	getDevicesWithCustomHeadersExample(axmService)

	// Example 4: Get limited devices with pagination control
	fmt.Println("\n=== Example 4: Get Devices with Pagination Control ===")
	getDevicesWithPaginationExample(axmService)

	fmt.Println("\n✓ All examples completed successfully!")
	fmt.Println("Note: Custom headers are automatically managed by the GetOrgDevices function")
}

func getAllDevicesExample(axmService *axm.Client) {
	// Create query builder (no specific parameters)
	queryBuilder := axmService.NewQueryBuilder()

	// Get all organization devices with automatic pagination
	devices, err := axmService.GetOrgDevices(queryBuilder)
	if err != nil {
		log.Printf("Failed to get organization devices: %v", err)
		return
	}

	fmt.Printf("Retrieved %d devices\n", len(devices))

	// Display first few devices (if any)
	if len(devices) > 0 {
		fmt.Printf("First device: %s (%s)\n", devices[0].Attributes.SerialNumber, devices[0].Attributes.DeviceModel)
		if len(devices) > 1 {
			fmt.Printf("Total devices retrieved: %d\n", len(devices))
		}
	} else {
		fmt.Println("No devices found")
	}
}

func getDevicesWithFieldsExample(axmService *axm.Client) {
	// Create query builder with field filtering
	queryBuilder := axmService.NewQueryBuilder().
		Fields("orgDevices", []string{
			"serialNumber",
			"deviceModel",
			"productFamily",
			"status",
			"addedToOrgDateTime",
		}).
		Limit(10) // Limit to 10 devices for demo

	// Get organization devices with field filtering
	devices, err := axmService.GetOrgDevices(queryBuilder)
	if err != nil {
		log.Printf("Failed to get organization devices with fields: %v", err)
		return
	}

	fmt.Printf("Retrieved %d devices with filtered fields\n", len(devices))

	// Display devices with filtered fields
	for i, device := range devices {
		if i >= 5 { // Show only first 5 for brevity
			break
		}
		fmt.Printf("Device %d:\n", i+1)
		fmt.Printf("  Serial: %s\n", device.Attributes.SerialNumber)
		fmt.Printf("  Model: %s\n", device.Attributes.DeviceModel)
		fmt.Printf("  Family: %s\n", device.Attributes.ProductFamily)
		fmt.Printf("  Status: %s\n", device.Attributes.Status)
		fmt.Printf("  Added: %s\n", device.Attributes.AddedToOrgDateTime)
		fmt.Println()
	}
}

func getDevicesWithCustomHeadersExample(axmService *axm.Client) {
	// This example demonstrates how the headers are handled internally
	// The GetOrgDevices function already sets appropriate headers
	queryBuilder := axmService.NewQueryBuilder().
		Filter("status", "AVAILABLE"). // Filter for available devices only
		Limit(5)

	devices, err := axmService.GetOrgDevices(queryBuilder)
	if err != nil {
		log.Printf("Failed to get devices with custom headers: %v", err)
		return
	}

	fmt.Printf("Retrieved %d available devices\n", len(devices))

	// Display available devices
	for i, device := range devices {
		fmt.Printf("Available Device %d: %s (%s)\n",
			i+1, device.Attributes.SerialNumber, device.Attributes.DeviceModel)
	}
}

func getDevicesWithPaginationExample(axmService *axm.Client) {
	// Create query builder with small limit to demonstrate pagination
	queryBuilder := axmService.NewQueryBuilder().
		Limit(3).                  // Very small limit to show pagination in action
		Sort("addedToOrgDateTime") // Sort by when added to org

	devices, err := axmService.GetOrgDevices(queryBuilder)
	if err != nil {
		log.Printf("Failed to get devices with pagination: %v", err)
		return
	}

	fmt.Printf("Retrieved %d devices with pagination control\n", len(devices))

	// Display devices showing pagination worked
	for i, device := range devices {
		if i >= 10 { // Show first 10 max
			fmt.Printf("... and %d more devices\n", len(devices)-10)
			break
		}
		fmt.Printf("Device %d: %s (Added: %s)\n",
			i+1, device.Attributes.SerialNumber, device.Attributes.AddedToOrgDateTime)
	}

	// Show how to work with the raw device data
	if len(devices) > 0 {
		fmt.Println("\n=== Raw JSON of First Device ===")
		deviceJSON, err := json.MarshalIndent(devices[0], "", "  ")
		if err == nil {
			fmt.Println(string(deviceJSON))
		}
	}
}
