package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-apple/axm"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/client"
	"github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get AppleCare Coverage Information Example ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	// Parse the private key (supports both RSA and ECDSA)
	privateKey, err := client.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// First, get a list of devices to find a device ID to query
	fmt.Println("\nStep 1: Getting organization devices to find a device ID...")

	listOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
		Limit: 5,
	}

	devicesResponse, err := client.Devices.GetOrganizationDevicesV1(ctx, listOptions)
	if err != nil {
		log.Fatalf("Error getting organization devices: %v", err)
	}

	if len(devicesResponse.Data) == 0 {
		log.Fatalf("No devices found in organization")
	}

	// Use the first device for AppleCare coverage information
	deviceID := devicesResponse.Data[0].ID
	fmt.Printf("Found device ID: %s (Serial: %s)\n", deviceID, devicesResponse.Data[0].Attributes.SerialNumber)

	// Step 2: Get AppleCare coverage information for the specific device
	fmt.Printf("\nStep 2: Getting AppleCare coverage information for device %s...\n", deviceID)

	// Example 1: Get all available AppleCare coverage fields
	fmt.Println("\n=== Example 1: Get All Available AppleCare Coverage Information ===")

	allFieldsOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldAppleCareStatus,
			devices.FieldAppleCarePaymentType,
			devices.FieldAppleCareDescription,
			devices.FieldAppleCareAgreementNumber,
			devices.FieldAppleCareStartDateTime,
			devices.FieldAppleCareEndDateTime,
			devices.FieldAppleCareIsRenewable,
			devices.FieldAppleCareIsCanceled,
			devices.FieldAppleCareContractCancelDateTime,
		},
		Limit: 100,
	}

	coverageInfo, err := client.Devices.GetAppleCareInformationByDeviceIDV1(ctx, deviceID, allFieldsOptions)
	if err != nil {
		log.Fatalf("Error getting AppleCare coverage information: %v", err)
	}

	// Display detailed AppleCare coverage information
	fmt.Printf("Found %d AppleCare coverage plan(s) for device %s:\n\n", len(coverageInfo.Data), deviceID)

	for i, coverage := range coverageInfo.Data {
		fmt.Printf("Coverage Plan %d:\n", i+1)
		fmt.Printf("  ID: %s\n", coverage.ID)
		fmt.Printf("  Type: %s\n", coverage.Type)
		fmt.Printf("  Description: %s\n", coverage.Attributes.Description)
		fmt.Printf("  Status: %s\n", coverage.Attributes.Status)
		fmt.Printf("  Payment Type: %s\n", coverage.Attributes.PaymentType)
		fmt.Printf("  Is Renewable: %v\n", coverage.Attributes.IsRenewable)
		fmt.Printf("  Is Canceled: %v\n", coverage.Attributes.IsCanceled)

		if coverage.Attributes.AgreementNumber != nil {
			fmt.Printf("  Agreement Number: %s\n", *coverage.Attributes.AgreementNumber)
		} else {
			fmt.Printf("  Agreement Number: None\n")
		}

		if coverage.Attributes.StartDateTime != nil {
			fmt.Printf("  Start Date: %s\n", coverage.Attributes.StartDateTime.Format(time.RFC3339))
		}

		if coverage.Attributes.EndDateTime != nil {
			fmt.Printf("  End Date: %s\n", coverage.Attributes.EndDateTime.Format(time.RFC3339))
		} else {
			fmt.Printf("  End Date: Ongoing/Not specified\n")
		}

		if coverage.Attributes.ContractCancelDateTime != nil {
			fmt.Printf("  Contract Cancel Date: %s\n", coverage.Attributes.ContractCancelDateTime.Format(time.RFC3339))
		}

		fmt.Println()
	}

	// Example 2: Get only specific fields
	fmt.Println("\n=== Example 2: Get Only Specific AppleCare Fields ===")

	specificFieldsOptions := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldAppleCareDescription,
			devices.FieldAppleCareStatus,
			devices.FieldAppleCarePaymentType,
		},
	}

	specificCoverageInfo, err := client.Devices.GetAppleCareInformationByDeviceIDV1(ctx, deviceID, specificFieldsOptions)
	if err != nil {
		log.Printf("Error getting specific AppleCare coverage information: %v", err)
	} else {
		fmt.Printf("Specific Fields Only:\n")
		for i, coverage := range specificCoverageInfo.Data {
			fmt.Printf("  Plan %d: %s (Status: %s, Payment: %s)\n",
				i+1,
				coverage.Attributes.Description,
				coverage.Attributes.Status,
				coverage.Attributes.PaymentType)
		}
	}

	// Example 3: Get AppleCare coverage with no field filtering (all fields)
	fmt.Println("\n=== Example 3: Get AppleCare Coverage (No Field Filtering) ===")

	noFilterCoverageInfo, err := client.Devices.GetAppleCareInformationByDeviceIDV1(ctx, deviceID, nil)
	if err != nil {
		log.Printf("Error getting unfiltered AppleCare coverage information: %v", err)
	} else {
		fmt.Printf("Unfiltered AppleCare coverage information retrieved successfully\n")
		fmt.Printf("Number of coverage plans: %d\n", len(noFilterCoverageInfo.Data))
	}

	// Example 4: Check for active AppleCare+ subscription
	fmt.Println("\n=== Example 4: Check for Active AppleCare+ Coverage ===")

	hasAppleCare := false
	for _, coverage := range coverageInfo.Data {
		if coverage.Attributes.Status == devices.AppleCareStatusActive &&
			coverage.Attributes.PaymentType != devices.PaymentTypeNone {
			hasAppleCare = true
			fmt.Printf("✓ Active AppleCare+ found:\n")
			fmt.Printf("  Type: %s\n", coverage.Attributes.Description)
			fmt.Printf("  Payment: %s\n", coverage.Attributes.PaymentType)

			if coverage.Attributes.EndDateTime != nil {
				daysRemaining := int(time.Until(*coverage.Attributes.EndDateTime).Hours() / 24)
				fmt.Printf("  Days Remaining: %d\n", daysRemaining)
			} else if coverage.Attributes.IsRenewable {
				fmt.Printf("  Status: Auto-renewing subscription\n")
			}
		}
	}

	if !hasAppleCare {
		fmt.Println("No active AppleCare+ coverage found (only standard warranty)")
	}

	// Example 5: Filter by payment type
	fmt.Println("\n=== Example 5: Categorize Coverage by Payment Type ===")

	warrantyPlans := []devices.AppleCareCoverage{}
	subscriptionPlans := []devices.AppleCareCoverage{}
	abeSubscriptionPlans := []devices.AppleCareCoverage{}

	for _, coverage := range coverageInfo.Data {
		switch coverage.Attributes.PaymentType {
		case devices.PaymentTypeNone:
			warrantyPlans = append(warrantyPlans, coverage)
		case devices.PaymentTypeSubscription:
			subscriptionPlans = append(subscriptionPlans, coverage)
		case devices.PaymentTypeABESubscription:
			abeSubscriptionPlans = append(abeSubscriptionPlans, coverage)
		}
	}

	fmt.Printf("Warranty Plans (No Payment): %d\n", len(warrantyPlans))
	for _, plan := range warrantyPlans {
		fmt.Printf("  - %s\n", plan.Attributes.Description)
	}

	fmt.Printf("AppleCare+ Subscriptions: %d\n", len(subscriptionPlans))
	for _, plan := range subscriptionPlans {
		fmt.Printf("  - %s\n", plan.Attributes.Description)
	}

	fmt.Printf("AppleCare+ for Business Essentials: %d\n", len(abeSubscriptionPlans))
	for _, plan := range abeSubscriptionPlans {
		fmt.Printf("  - %s\n", plan.Attributes.Description)
	}

	// Example 6: Try to get coverage for a non-existent device (error handling)
	fmt.Println("\n=== Example 6: Error Handling (Non-existent Device) ===")

	fakeDeviceID := "non-existent-device-id"
	_, err = client.Devices.GetAppleCareInformationByDeviceIDV1(ctx, fakeDeviceID, nil)
	if err != nil {
		fmt.Printf("Expected error for non-existent device: %v\n", err)
	}

	// Example 7: Pretty print full JSON response
	fmt.Println("\n=== Example 7: Full JSON Response ===")
	jsonData, err := json.MarshalIndent(coverageInfo, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	// Example 8: Check coverage expiration dates
	fmt.Println("\n=== Example 8: Coverage Expiration Summary ===")

	now := time.Now()
	for _, coverage := range coverageInfo.Data {
		if coverage.Attributes.IsCanceled {
			fmt.Printf("✗ %s: CANCELED\n", coverage.Attributes.Description)
			if coverage.Attributes.ContractCancelDateTime != nil {
				fmt.Printf("  Canceled on: %s\n", coverage.Attributes.ContractCancelDateTime.Format("2006-01-02"))
			}
			continue
		}

		if coverage.Attributes.EndDateTime == nil {
			fmt.Printf("○ %s: ONGOING (No end date)\n", coverage.Attributes.Description)
			continue
		}

		daysUntilExpiry := int(coverage.Attributes.EndDateTime.Sub(now).Hours() / 24)

		if daysUntilExpiry < 0 {
			fmt.Printf("✗ %s: EXPIRED (%s)\n",
				coverage.Attributes.Description,
				coverage.Attributes.EndDateTime.Format("2006-01-02"))
		} else if daysUntilExpiry < 30 {
			fmt.Printf("⚠ %s: EXPIRING SOON (%d days, %s)\n",
				coverage.Attributes.Description,
				daysUntilExpiry,
				coverage.Attributes.EndDateTime.Format("2006-01-02"))
		} else {
			fmt.Printf("✓ %s: ACTIVE (%d days remaining, expires %s)\n",
				coverage.Attributes.Description,
				daysUntilExpiry,
				coverage.Attributes.EndDateTime.Format("2006-01-02"))
		}

		if coverage.Attributes.IsRenewable {
			fmt.Printf("  Auto-renewable: Yes\n")
		}
	}

	// Example 9: Get coverage for multiple devices
	fmt.Println("\n=== Example 9: Get AppleCare Coverage for Multiple Devices ===")

	if len(devicesResponse.Data) > 1 {
		fmt.Printf("Checking AppleCare coverage for first %d devices:\n\n", min(3, len(devicesResponse.Data)))

		for i, dev := range devicesResponse.Data[:min(3, len(devicesResponse.Data))] {
			fmt.Printf("Device %d (Serial: %s):\n", i+1, dev.Attributes.SerialNumber)

			coverage, err := client.Devices.GetAppleCareInformationByDeviceIDV1(ctx, dev.ID, &devices.RequestQueryOptions{
				Fields: []string{
					devices.FieldAppleCareDescription,
					devices.FieldAppleCareStatus,
					devices.FieldAppleCarePaymentType,
				},
			})

			if err != nil {
				fmt.Printf("  Error: %v\n", err)
				continue
			}

			if len(coverage.Data) == 0 {
				fmt.Printf("  No AppleCare coverage found\n")
			} else {
				for _, plan := range coverage.Data {
					fmt.Printf("  - %s (%s)\n", plan.Attributes.Description, plan.Attributes.Status)
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("\n=== Example Complete ===")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
