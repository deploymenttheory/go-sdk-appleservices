package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/axm"
)

func main() {
	// Create client using GitLab pattern - matches your desired usage exactly
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
-----END EC PRIVATE KEY-----`

	// Parse the private key from PEM format
	privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := axm.NewClient("bd4bd60b-1111-1111-1111-3ed8f6dc4bd1", "BUSINESSAPI.3bb3a62b-1111-1111-1111-a69b86201c5a", privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get device management services
	response, err := client.
		DeviceManagement.
		GetDeviceManagementServices(ctx, nil)

	if err != nil {
		log.Fatalf("Failed to get device management services: %v", err)
	}

	// Pretty print the JSON response
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		fmt.Printf("Device management services: %+v\n", response)
	} else {
		fmt.Printf("Device management services:\n%s\n", responseJSON)
	}

}
