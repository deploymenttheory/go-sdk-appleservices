package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	isbnCode = "9780316069359"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup by ISBN Example ===")
	fmt.Printf("Looking up book by 13-digit ISBN (%s):\n", isbnCode)

	params := itunes_search.NewLookupParams().
		ISBN(isbnCode).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up ISBN %s: %v", isbnCode, err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
