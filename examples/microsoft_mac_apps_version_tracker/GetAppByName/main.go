package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	msapps "github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker"
	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/services/apps"
)

func main() {
	fmt.Println("=== Microsoft Mac Apps - Get App By Name Example ===")
	fmt.Println()

	// Create a new client with default configuration
	client, err := msapps.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Find Microsoft Teams by name
	fmt.Println("=== Example 1: Find Microsoft Teams ===")

	teams, err := client.Apps.GetAppByName(ctx, apps.AppNameTeams)
	if err != nil {
		log.Fatalf("Error finding Teams: %v", err)
	}

	fmt.Printf("Found: %s\n", teams.Name)
	fmt.Printf("Bundle ID: %s\n", teams.BundleID)
	fmt.Printf("Version: %s\n", teams.Version)
	fmt.Printf("Size: %.2f MB (%.2f GB)\n", teams.SizeMB, teams.SizeMB/1024)
	fmt.Printf("Package ID: %s\n", teams.PackageID)
	fmt.Printf("Download URL: %s\n", teams.DownloadURL)
	fmt.Printf("Direct URL: %s\n", teams.DirectURL)

	// Display components
	if len(teams.Components) > 0 {
		fmt.Printf("\nBundled Components (%d):\n", teams.ComponentCount)
		for _, comp := range teams.Components {
			fmt.Printf("  - %s (v%s)\n", comp.Name, comp.Version)
		}
	}

	// Display timestamps
	detected, err := teams.ParseDetectedTime()
	if err == nil {
		fmt.Printf("\nDetected: %s\n", detected.Format("2006-01-02 15:04:05"))
	}

	lastModified, err := teams.ParseLastModifiedTime()
	if err == nil {
		fmt.Printf("Last Modified: %s\n", lastModified.Format("2006-01-02 15:04:05"))
	}

	// Example 2: Find Microsoft Defender
	fmt.Println("\n=== Example 2: Find Defender for Mac ===")

	defender, err := client.Apps.GetAppByName(ctx, apps.AppNameDefender)
	if err != nil {
		log.Printf("Error finding Defender: %v", err)
	} else {
		fmt.Printf("Found: %s\n", defender.Name)
		fmt.Printf("Version: %s\n", defender.Version)
		fmt.Printf("Bundle ID: %s\n", defender.BundleID)
		fmt.Printf("Size: %.2f MB\n", defender.SizeMB)

		// Show DLP components
		if len(defender.Components) > 0 {
			fmt.Printf("\nDefender includes %d components:\n", len(defender.Components))
			for _, comp := range defender.Components {
				fmt.Printf("  - %s (v%s)\n", comp.Name, comp.Version)
			}
		}
	}

	// Example 3: Find Office 365 Suite
	fmt.Println("\n=== Example 3: Find Office 365 for Mac ===")

	office, err := client.Apps.GetAppByName(ctx, apps.AppNameOffice365)
	if err != nil {
		log.Printf("Error finding Office 365: %v", err)
	} else {
		fmt.Printf("Found: %s\n", office.Name)
		fmt.Printf("Version: %s\n", office.Version)
		fmt.Printf("Size: %.2f GB\n", office.SizeMB/1024)
		fmt.Printf("\nOffice 365 Suite includes:\n")
		for _, comp := range office.Components {
			if comp.BundleID != nil {
				fmt.Printf("  - %s (v%s)\n", comp.Name, comp.Version)
			}
		}
	}

	// Example 4: Check for multiple apps by name
	fmt.Println("\n=== Example 4: Check Multiple Applications ===")

	appNames := []string{
		apps.AppNameExcel,
		apps.AppNameWord,
		apps.AppNamePowerPoint,
		apps.AppNameOutlook,
		apps.AppNameOneNote,
		apps.AppNameTeams,
		apps.AppNameOneDrive,
		apps.AppNameEdge,
		apps.AppNameVSCode,
	}

	fmt.Println("Checking for Microsoft applications:")
	for _, appName := range appNames {
		app, err := client.Apps.GetAppByName(ctx, appName)
		if err != nil {
			fmt.Printf("  ✗ %s: Not found\n", appName)
		} else {
			fmt.Printf("  ✓ %s: v%s\n", app.Name, app.Version)
		}
	}

	// Example 5: Compare versions and timestamps
	fmt.Println("\n=== Example 5: Version and Update Analysis ===")

	allApps, err := client.Apps.GetLatestApps(ctx)
	if err != nil {
		log.Printf("Error getting all apps: %v", err)
	} else {
		fmt.Println("Recent updates (within last 30 days):")
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

		for _, app := range allApps.Apps {
			detected, err := app.ParseDetectedTime()
			if err == nil && detected.After(thirtyDaysAgo) {
				fmt.Printf("  - %s (v%s) - detected %s\n",
					app.Name,
					app.Version,
					detected.Format("2006-01-02"))
			}
		}
	}

	// Example 6: Generate deployment script
	fmt.Println("\n=== Example 6: Generate Download Script for Specific Apps ===")

	targetApps := []string{
		apps.AppNameExcel,
		apps.AppNameWord,
		apps.AppNamePowerPoint,
		apps.AppNameTeams,
	}

	fmt.Println("#!/bin/bash")
	fmt.Println("# Microsoft Mac Apps Download Script")
	fmt.Println("# Generated:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	for _, appName := range targetApps {
		app, err := client.Apps.GetAppByName(ctx, appName)
		if err != nil {
			fmt.Printf("# Error: %s not found\n", appName)
			continue
		}

		fmt.Printf("# %s v%s (%.2f MB)\n", app.Name, app.Version, app.SizeMB)
		fmt.Printf("curl -L '%s' -o '%s.pkg'\n", app.DownloadURL, app.PackageID)
		fmt.Printf("# SHA256: %s\n", app.SHA256)
		fmt.Println()
	}

	// Example 7: Try to find a non-existent app (error handling)
	fmt.Println("=== Example 7: Error Handling ===")

	_, err = client.Apps.GetAppByName(ctx, "Non-existent Application")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 8: Pretty print full JSON response
	fmt.Println("\n=== Example 8: Full JSON Response ===")

	jsonData, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n=== Example Complete ===")
}
