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
		"term":      "Jack Johnson",
		"media":     "music",
		"entity":    "musicTrack",
		"attribute": "artistTerm",
		"limit":     "25",
		"country":   "US",
		"lang":      "en_us",
		"explicit":  "No",
		"version":   "2",
	}
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Advanced Search with Raw Parameters Example ===")
	fmt.Println("Using custom parameter map for advanced search:")

	response, err := itunesClient.Search(customParams)
	if err != nil {
		log.Printf("Error performing advanced search: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
