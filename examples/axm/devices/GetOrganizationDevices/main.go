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
	fmt.Println("=== Apple Business Manager Test Example ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	privateKey, err := client.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("\nFetching organization devices...")

	options := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
		Limit: 5, // Limit to 5 devices for this example
	}

	response, err := client.Devices.GetOrganizationDevices(ctx, options)
	if err != nil {
		log.Fatalf("Error getting devices: %v", err)
	}

	fmt.Printf("Found %d devices:\n\n", len(response.Data))

	for i, device := range response.Data {
		fmt.Printf("Device %d:\n", i+1)
		fmt.Printf("  ID: %s\n", device.ID)
		fmt.Printf("  Serial: %s\n", device.Attributes.SerialNumber)
		fmt.Printf("  Model: %s\n", device.Attributes.DeviceModel)
		fmt.Printf("  Status: %s\n", device.Attributes.Status)
		fmt.Println()
	}

	if response.Links != nil && response.Links.Next != "" {
		fmt.Println("Note: More devices are available on additional pages.")
		fmt.Printf("Next page URL: %s\n", response.Links.Next)
	}

	fmt.Println("=== Full JSON Response ===")
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n=== Test Complete ===")
}
