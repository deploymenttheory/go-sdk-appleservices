package main

import (
	"context"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/axm"
)

func main() {
	// This matches your desired GitLab pattern exactly
	client, err := axm.NewClient("keyID", "issuerID", "privateKey")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Use the client just like GitLab's pattern
	response, err := client.DeviceManagement.GetDeviceManagementServices(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get device management services: %v", err)
	}

	fmt.Printf("Device management services: %v", response)

	// Also works for devices
	devices, err := client.Devices.GetOrganizationDevices(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	fmt.Printf("Devices: %v", devices)
}