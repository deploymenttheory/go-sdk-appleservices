package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	upcCode    = "720642462928"
	entityType = "song"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup UPC with Tracks Example ===")
	fmt.Printf("Looking up album by UPC (%s) including tracks:\n", upcCode)

	params := itunes_search.NewLookupParams().
		UPC(upcCode).
		Entity(entityType).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up UPC with tracks: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
