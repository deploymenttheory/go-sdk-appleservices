package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

var (
	customParams = map[string]string{
		"amgArtistId": "468749,5723",
		"entity":      "album",
		"limit":       "10",
		"sort":        "recent",
		"country":     "US",
	}
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Advanced Lookup with Raw Parameters Example ===")
	fmt.Println("Using custom parameter map for advanced lookup:")

	response, err := itunesClient.Lookup(customParams)
	if err != nil {
		log.Printf("Error performing advanced lookup: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
