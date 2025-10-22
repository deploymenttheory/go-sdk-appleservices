package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devicemanagement"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Device Management Services Example ===")

	// Use credentials directly for testing
	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	// Parse the private key
	privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Create client using GitLab pattern - matches the v3 pattern exactly
	client, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Example 1: Get all MDM servers with default options
	fmt.Println("\n=== Example 1: Get All MDM Servers (Default Options) ===")

	response, err := client.DeviceManagement.GetDeviceManagementServices(ctx, nil)
	if err != nil {
		log.Printf("Error getting MDM servers: %v", err)
	} else {
		fmt.Printf("Found %d MDM servers\n", len(response.Data))

		for i, server := range response.Data {
			fmt.Printf("MDM Server %d:\n", i+1)
			fmt.Printf("  ID: %s\n", server.ID)
			fmt.Printf("  Type: %s\n", server.Type)
			if server.Attributes != nil {
				fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
				fmt.Printf("  Server Type: %s\n", server.Attributes.ServerType)
				if server.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", server.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
				if server.Attributes.UpdatedDateTime != nil {
					fmt.Printf("  Updated: %s\n", server.Attributes.UpdatedDateTime.Format(time.RFC3339))
				}
			}
			fmt.Println()
		}

		// Check pagination
		if response.Links != nil && response.Links.Next != "" {
			fmt.Println("More MDM servers available on next page!")
		}
	}

	// Example 2: Get MDM servers with specific fields
	fmt.Println("\n=== Example 2: Get MDM Servers with Specific Fields ===")

	options := &devicemanagement.GetMDMServersOptions{
		Fields: []string{
			"serverName",
			"serverType",
			"createdDateTime",
		},
		Limit: 10,
	}

	response, err = client.DeviceManagement.GetDeviceManagementServices(ctx, options)
	if err != nil {
		log.Printf("Error getting MDM servers with options: %v", err)
	} else {
		fmt.Printf("Found %d MDM servers (limited to %d)\n", len(response.Data), options.Limit)

		for i, server := range response.Data {
			fmt.Printf("MDM Server %d:\n", i+1)
			fmt.Printf("  ID: %s\n", server.ID)
			if server.Attributes != nil {
				fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
				fmt.Printf("  Type: %s\n", server.Attributes.ServerType)
				if server.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", server.Attributes.CreatedDateTime.Format(time.RFC3339))
				}
			}
			fmt.Println()
		}
	}

	// Example 3: Get MDM servers with pagination limit
	fmt.Println("\n=== Example 3: Get MDM Servers with Pagination Limit ===")

	paginationOptions := &devicemanagement.GetMDMServersOptions{
		Limit: 5,
	}

	paginatedResponse, err := client.DeviceManagement.GetDeviceManagementServices(ctx, paginationOptions)
	if err != nil {
		log.Printf("Error getting paginated MDM servers: %v", err)
	} else {
		fmt.Printf("Retrieved %d MDM servers (limit: %d)\n", len(paginatedResponse.Data), paginationOptions.Limit)

		if paginatedResponse.Links != nil && paginatedResponse.Links.Next != "" {
			fmt.Printf("Next page available: %s\n", paginatedResponse.Links.Next)
		}

		if paginatedResponse.Meta != nil && paginatedResponse.Meta.Paging != nil {
			fmt.Printf("Pagination info - Limit: %d\n", paginatedResponse.Meta.Paging.Limit)
			if paginatedResponse.Meta.Paging.Total > 0 {
				fmt.Printf("Total servers: %d\n", paginatedResponse.Meta.Paging.Total)
			}
		}
	}

	// Example 4: Get all available fields
	fmt.Println("\n=== Example 4: Get MDM Servers with All Available Fields ===")

	allFieldsOptions := &devicemanagement.GetMDMServersOptions{
		Fields: []string{
			"serverName",
			"serverType",
			"createdDateTime",
			"updatedDateTime",
			"devices",
		},
	}

	allFieldsResponse, err := client.DeviceManagement.GetDeviceManagementServices(ctx, allFieldsOptions)
	if err != nil {
		log.Printf("Error getting MDM servers with all fields: %v", err)
	} else {
		fmt.Printf("Retrieved %d MDM servers with all fields\n", len(allFieldsResponse.Data))

		for i, server := range allFieldsResponse.Data {
			fmt.Printf("Complete MDM Server %d Info:\n", i+1)
			fmt.Printf("  ID: %s\n", server.ID)
			fmt.Printf("  Type: %s\n", server.Type)

			if server.Attributes != nil {
				fmt.Printf("  Server Name: %s\n", server.Attributes.ServerName)
				fmt.Printf("  Server Type: %s\n", server.Attributes.ServerType)

				if server.Attributes.CreatedDateTime != nil {
					fmt.Printf("  Created: %s\n", server.Attributes.CreatedDateTime.Format(time.RFC3339))
				}

				if server.Attributes.UpdatedDateTime != nil {
					fmt.Printf("  Updated: %s\n", server.Attributes.UpdatedDateTime.Format(time.RFC3339))
				}
			}

			if server.Relationships != nil && server.Relationships.Devices != nil {
				fmt.Printf("  Devices Relationship: %+v\n", server.Relationships.Devices)
			}

			fmt.Println()
		}
	}

	// Example 5: Error handling - test with invalid options
	fmt.Println("\n=== Example 5: Error Handling ===")

	// Test with very large limit (should be capped at 1000)
	largeOptions := &devicemanagement.GetMDMServersOptions{
		Limit: 5000, // This will be capped at 1000
	}

	largeResponse, err := client.DeviceManagement.GetDeviceManagementServices(ctx, largeOptions)
	if err != nil {
		log.Printf("Error with large limit: %v", err)
	} else {
		fmt.Printf("Large limit test - Retrieved %d servers (requested limit was capped)\n", len(largeResponse.Data))
	}

	// Example 6: Pretty print JSON response
	fmt.Println("\n=== Example 6: Full JSON Response ===")
	if response != nil {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Printf("Error marshaling response to JSON: %v", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
