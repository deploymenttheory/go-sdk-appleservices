package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	msapps "github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker"
	"github.com/deploymenttheory/go-api-sdk-apple/microsoft_mac_apps_version_tracker/services/apps"
)

func main() {
	fmt.Println("=== Microsoft Mac Apps - Get App By Bundle ID Example ===")
	fmt.Println()

	// Create a new client with default configuration
	client, err := msapps.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Find Microsoft Excel by bundle ID
	fmt.Println("=== Example 1: Find Microsoft Excel ===")

	excel, err := client.Apps.GetAppByBundleID(ctx, apps.BundleIDExcel)
	if err != nil {
		log.Fatalf("Error finding Excel: %v", err)
	}

	fmt.Printf("Found: %s\n", excel.Name)
	fmt.Printf("Bundle ID: %s\n", excel.BundleID)
	fmt.Printf("Version: %s\n", excel.Version)
	fmt.Printf("Package ID: %s\n", excel.PackageID)
	fmt.Printf("Download URL: %s\n", excel.DownloadURL)
	fmt.Printf("Direct URL: %s\n", excel.DirectURL)
	fmt.Printf("Size: %.2f MB (%.2f GB)\n", excel.SizeMB, excel.SizeMB/1024)
	fmt.Printf("SHA256: %s\n", excel.SHA256)
	fmt.Printf("ETag: %s\n", excel.ETag)

	if excel.InstallKB != nil {
		fmt.Printf("Install Size: %d KB (%.2f MB)\n", *excel.InstallKB, float64(*excel.InstallKB)/1024)
	}
	if excel.NumFiles != nil {
		fmt.Printf("Number of Files: %d\n", *excel.NumFiles)
	}

	// Display components
	if excel.ComponentCount > 0 {
		fmt.Printf("\nBundled Components (%d):\n", excel.ComponentCount)
		for _, comp := range excel.Components {
			fmt.Printf("  - %s (v%s)\n", comp.Name, comp.Version)
			if comp.BundleID != nil {
				fmt.Printf("    Bundle ID: %s\n", *comp.BundleID)
			}
			fmt.Printf("    Package ID: %s\n", comp.PackageID)
		}
	}

	// Example 2: Find Microsoft Teams
	fmt.Println("\n=== Example 2: Find Microsoft Teams ===")

	teams, err := client.Apps.GetAppByBundleID(ctx, apps.BundleIDTeams)
	if err != nil {
		log.Printf("Error finding Teams: %v", err)
	} else {
		fmt.Printf("Found: %s\n", teams.Name)
		fmt.Printf("Version: %s\n", teams.Version)
		fmt.Printf("Bundle ID: %s\n", teams.BundleID)
		fmt.Printf("Size: %.2f MB\n", teams.SizeMB)
		fmt.Printf("Download URL: %s\n", teams.DownloadURL)

		if len(teams.Components) > 0 {
			fmt.Printf("\nComponents:\n")
			for _, comp := range teams.Components {
				fmt.Printf("  - %s (v%s)\n", comp.Name, comp.Version)
			}
		}
	}

	// Example 3: Check for multiple specific applications
	fmt.Println("\n=== Example 3: Check Multiple Applications ===")

	appsToCheck := []struct {
		name     string
		bundleID string
	}{
		{"Excel", apps.BundleIDExcel},
		{"Word", apps.BundleIDWord},
		{"PowerPoint", apps.BundleIDPowerPoint},
		{"Outlook", apps.BundleIDOutlook},
		{"Teams", apps.BundleIDTeams},
		{"OneDrive", apps.BundleIDOneDrive},
		{"Defender", apps.BundleIDDefender},
		{"VS Code", apps.BundleIDVSCode},
		{"Edge", apps.BundleIDEdge},
		{"Company Portal", apps.BundleIDCompanyPortal},
	}

	fmt.Println("Checking for common Microsoft applications:")
	for _, appInfo := range appsToCheck {
		app, err := client.Apps.GetAppByBundleID(ctx, appInfo.bundleID)
		if err != nil {
			fmt.Printf("  ✗ %s: Not found\n", appInfo.name)
		} else {
			fmt.Printf("  ✓ %s: v%s (%.2f MB)\n", app.Name, app.Version, app.SizeMB)
		}
	}

	// Example 4: Find apps with specific components
	fmt.Println("\n=== Example 4: Find Apps with AutoUpdate ===")

	allApps, err := client.Apps.GetLatestApps(ctx)
	if err != nil {
		log.Printf("Error getting apps: %v", err)
	} else {
		fmt.Println("Applications bundled with Microsoft AutoUpdate:")
		for _, app := range allApps.Apps {
			for _, comp := range app.Components {
				if comp.BundleID != nil && *comp.BundleID == apps.BundleIDAutoUpdate {
					fmt.Printf("  - %s (v%s) includes AutoUpdate v%s\n",
						app.Name,
						app.Version,
						comp.Version)
					break
				}
			}
		}
	}

	// Example 5: Try to find a non-existent app (error handling)
	fmt.Println("\n=== Example 5: Error Handling ===")

	_, err = client.Apps.GetAppByBundleID(ctx, "com.nonexistent.app")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 6: Pretty print full JSON response
	fmt.Println("\n=== Example 6: Full JSON Response ===")

	jsonData, err := json.MarshalIndent(excel, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n=== Example Complete ===")
}
