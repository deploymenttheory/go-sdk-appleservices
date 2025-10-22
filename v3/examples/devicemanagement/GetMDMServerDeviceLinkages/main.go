package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/v3/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/v3/devicemanagement"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get MDM Server Device Linkages Example ===")

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

	// Step 1: Get MDM servers to find an MDM server ID
	fmt.Println("\nStep 1: Getting MDM servers to find an MDM server ID...")

	serversResponse, err := client.DeviceManagement.GetDeviceManagementServices(ctx, &devicemanagement.GetMDMServersOptions{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Error getting MDM servers: %v", err)
	}

	if len(serversResponse.Data) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	// Use the first MDM server for examples
	mdmServerID := serversResponse.Data[0].ID
	serverName := ""
	if serversResponse.Data[0].Attributes != nil {
		serverName = serversResponse.Data[0].Attributes.ServerName
	}

	fmt.Printf("Using MDM Server ID: %s", mdmServerID)
	if serverName != "" {
		fmt.Printf(" (Name: %s)", serverName)
	}
	fmt.Println()

	// Example 1: Get all device linkages for MDM server (default options)
	fmt.Println("\n=== Example 1: Get All Device Linkages (Default Options) ===")

	linkagesResponse, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, mdmServerID, nil)
	if err != nil {
		log.Printf("Error getting device linkages: %v", err)
	} else {
		fmt.Printf("Found %d device linkages for MDM server %s\n", len(linkagesResponse.Data), mdmServerID)

		for i, linkage := range linkagesResponse.Data {
			fmt.Printf("Device Linkage %d:\n", i+1)
			fmt.Printf("  Type: %s\n", linkage.Type)
			fmt.Printf("  Device ID: %s\n", linkage.ID)
		}

		// Check pagination
		if linkagesResponse.Links != nil && linkagesResponse.Links.Next != "" {
			fmt.Println("More device linkages available on next page!")
		}

		// Show pagination metadata
		if linkagesResponse.Meta != nil && linkagesResponse.Meta.Paging != nil {
			fmt.Printf("Pagination - Limit: %d", linkagesResponse.Meta.Paging.Limit)
			if linkagesResponse.Meta.Paging.NextCursor != "" {
				fmt.Printf(", Next Cursor: %s", linkagesResponse.Meta.Paging.NextCursor)
			}
			fmt.Println()
		}
	}

	// Example 2: Get device linkages with specific limit
	fmt.Println("\n=== Example 2: Get Device Linkages with Limit ===")

	limitOptions := &devicemanagement.GetMDMServerDeviceLinkagesOptions{
		Limit: 10,
	}

	limitedResponse, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, mdmServerID, limitOptions)
	if err != nil {
		log.Printf("Error getting limited device linkages: %v", err)
	} else {
		fmt.Printf("Retrieved %d device linkages (limit: %d)\n", len(limitedResponse.Data), limitOptions.Limit)

		for i, linkage := range limitedResponse.Data {
			fmt.Printf("  %d. Device ID: %s (Type: %s)\n", i+1, linkage.ID, linkage.Type)
		}

		if limitedResponse.Links != nil && limitedResponse.Links.Next != "" {
			fmt.Printf("Next page URL: %s\n", limitedResponse.Links.Next)
		}
	}

	// Example 3: Get device linkages with small limit for pagination demo
	fmt.Println("\n=== Example 3: Pagination Demo (Small Limit) ===")

	smallLimitOptions := &devicemanagement.GetMDMServerDeviceLinkagesOptions{
		Limit: 3,
	}

	paginationResponse, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, mdmServerID, smallLimitOptions)
	if err != nil {
		log.Printf("Error getting paginated device linkages: %v", err)
	} else {
		fmt.Printf("Page 1: Retrieved %d device linkages\n", len(paginationResponse.Data))

		for i, linkage := range paginationResponse.Data {
			fmt.Printf("  %d. Device ID: %s\n", i+1, linkage.ID)
		}

		// Demonstrate pagination
		if paginationResponse.Links != nil && paginationResponse.Links.Next != "" {
			fmt.Println("\nFetching next page...")

			// Note: In a real application, you would use the client's GetNextPage method
			// or extract cursor from the next URL for subsequent requests
			fmt.Printf("Next page would be available at: %s\n", paginationResponse.Links.Next)
		}
	}

	// Example 4: Error handling - invalid MDM server ID
	fmt.Println("\n=== Example 4: Error Handling (Invalid MDM Server ID) ===")

	invalidServerID := "invalid-mdm-server-id"
	_, err = client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, invalidServerID, nil)
	if err != nil {
		fmt.Printf("Expected error for invalid MDM server ID: %v\n", err)
	}

	// Example 5: Error handling - empty MDM server ID
	fmt.Println("\n=== Example 5: Error Handling (Empty MDM Server ID) ===")

	_, err = client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, "", nil)
	if err != nil {
		fmt.Printf("Expected error for empty MDM server ID: %v\n", err)
	}

	// Example 6: Test with maximum limit
	fmt.Println("\n=== Example 6: Test with Maximum Limit ===")

	maxLimitOptions := &devicemanagement.GetMDMServerDeviceLinkagesOptions{
		Limit: 1000, // API maximum
	}

	maxResponse, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, mdmServerID, maxLimitOptions)
	if err != nil {
		log.Printf("Error with max limit: %v", err)
	} else {
		fmt.Printf("Max limit test - Retrieved %d device linkages (limit: %d)\n",
			len(maxResponse.Data), maxLimitOptions.Limit)
	}

	// Example 7: Test with over-maximum limit (should be capped)
	fmt.Println("\n=== Example 7: Test with Over-Maximum Limit (Should be Capped) ===")

	overMaxOptions := &devicemanagement.GetMDMServerDeviceLinkagesOptions{
		Limit: 5000, // This should be capped at 1000
	}

	cappedResponse, err := client.DeviceManagement.GetMDMServerDeviceLinkages(ctx, mdmServerID, overMaxOptions)
	if err != nil {
		log.Printf("Error with over-max limit: %v", err)
	} else {
		fmt.Printf("Over-max limit test - Retrieved %d device linkages (requested: %d, should be capped at 1000)\n",
			len(cappedResponse.Data), overMaxOptions.Limit)
	}

	// Example 8: Pretty print JSON response
	fmt.Println("\n=== Example 8: Full JSON Response ===")
	if linkagesResponse != nil {
		jsonData, err := json.MarshalIndent(linkagesResponse, "", "  ")
		if err != nil {
			log.Printf("Error marshaling response to JSON: %v", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
