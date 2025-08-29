package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

var (
	amgArtistIDs = []string{"468749", "5723"}
)

const (
	entityType  = "song"
	resultLimit = 5
	sortOrder   = "recent"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Lookup Recent Songs for Multiple Artists Example ===")
	fmt.Printf("Looking up %d most recent songs for multiple artists by AMG artist IDs (%v):\n", resultLimit, amgArtistIDs)

	params := itunes_search.NewLookupParams().
		AMGArtistIDs(amgArtistIDs).
		Entity(entityType).
		Limit(resultLimit).
		Sort(sortOrder).
		Build()

	response, err := itunesClient.Lookup(params)
	if err != nil {
		log.Printf("Error looking up recent songs for multiple AMG artist IDs: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}