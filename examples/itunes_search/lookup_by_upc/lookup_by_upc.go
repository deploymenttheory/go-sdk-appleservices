package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	upcCode = "720642462928"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup by UPC Example ===")
	fmt.Printf("Looking up album by UPC (%s):\n", upcCode)

	params := itunes_search.NewLookupParams().
		UPC(upcCode).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up UPC %s: %v", upcCode, err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}