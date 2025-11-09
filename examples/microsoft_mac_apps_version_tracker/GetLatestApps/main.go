package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	msapps "github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker"
)

func main() {
	fmt.Println("=== Microsoft Mac Apps - Get Latest Apps Example ===")
	fmt.Println()

	// Create a new client with default configuration
	client, err := msapps.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Get all latest Microsoft Mac application versions
	fmt.Println("=== Retrieving Latest Application Versions ===")

	apps, err := client.Apps.GetLatestApps(ctx)
	if err != nil {
		log.Fatalf("Error getting latest apps: %v", err)
	}

	fmt.Printf("API Data Generated: %s\n", apps.Generated)
	fmt.Printf("Total Applications: %d\n\n", len(apps.Apps))

	// Display summary of all apps
	fmt.Println("Available Microsoft Mac Applications:")
	fmt.Println("=" + string(make([]byte, 80)))

	for i, app := range apps.Apps {
		fmt.Printf("\n%d. %s\n", i+1, app.Name)
		fmt.Printf("   Bundle ID: %s\n", app.BundleID)
		fmt.Printf("   Version: %s\n", app.Version)
		fmt.Printf("   Size: %.2f MB (%.2f GB)\n", app.SizeMB, app.SizeMB/1024)
		fmt.Printf("   Package ID: %s\n", app.PackageID)
		fmt.Printf("   Download URL: %s\n", app.DownloadURL)
		fmt.Printf("   Direct URL: %s\n", app.DirectURL)

		if app.ComponentCount > 0 {
			fmt.Printf("   Components (%d):\n", app.ComponentCount)
			for _, comp := range app.Components {
				fmt.Printf("     - %s (v%s)\n", comp.Name, comp.Version)
			}
		}

		// Parse and display timestamps
		detected, err := app.ParseDetectedTime()
		if err == nil {
			fmt.Printf("   Detected: %s\n", detected.Format("2006-01-02 15:04:05"))
		}

		lastModified, err := app.ParseLastModifiedTime()
		if err == nil {
			fmt.Printf("   Last Modified: %s\n", lastModified.Format("2006-01-02 15:04:05"))
		}

		// Display file information
		if app.NumFiles != nil {
			fmt.Printf("   Files: %d\n", *app.NumFiles)
		}
		if app.InstallKB != nil {
			fmt.Printf("   Install Size: %d KB\n", *app.InstallKB)
		}
	}

	// Calculate and display total size
	fmt.Println("\n" + string(make([]byte, 80)))
	fmt.Println("\n=== Storage Requirements ===")

	var totalBytes int64
	for _, app := range apps.Apps {
		totalBytes += app.SizeBytes
	}

	totalMB := float64(totalBytes) / (1024 * 1024)
	totalGB := totalMB / 1024

	fmt.Printf("Total download size for all applications: %.2f MB (%.2f GB)\n", totalMB, totalGB)
	fmt.Printf("Average application size: %.2f MB\n", totalMB/float64(len(apps.Apps)))

	// Find largest and smallest apps
	var largest, smallest *string
	largestSize, smallestSize := 0.0, apps.Apps[0].SizeMB

	for _, app := range apps.Apps {
		if app.SizeMB > largestSize {
			largestSize = app.SizeMB
			name := app.Name
			largest = &name
		}
		if app.SizeMB < smallestSize {
			smallestSize = app.SizeMB
			name := app.Name
			smallest = &name
		}
	}

	if largest != nil {
		fmt.Printf("\nLargest app: %s (%.2f MB)\n", *largest, largestSize)
	}
	if smallest != nil {
		fmt.Printf("Smallest app: %s (%.2f MB)\n", *smallest, smallestSize)
	}

	// Display JSON output for first app
	fmt.Println("\n=== Sample JSON Response (First App) ===")

	if len(apps.Apps) > 0 {
		jsonData, err := json.MarshalIndent(apps.Apps[0], "", "  ")
		if err != nil {
			log.Printf("Error marshaling response to JSON: %v", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
