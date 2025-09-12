package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	amgArtistID = "468749"
	entityType  = "album"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup Albums by AMG Artist ID Example ===")
	fmt.Printf("Looking up all albums for Jack Johnson (AMG ID %s):\n", amgArtistID)

	params := itunes_search.NewLookupParams().
		AMGArtistID(amgArtistID).
		Entity(entityType).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up albums for AMG artist ID %s: %v", amgArtistID, err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
