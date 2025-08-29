package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	artistID = 909253
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup by iTunes ID Example ===")
	fmt.Printf("Looking up Jack Johnson by iTunes artist ID (%d):\n", artistID)

	params := itunes_search.NewLookupParams().
		ID(artistID).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up ID %d: %v", artistID, err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
