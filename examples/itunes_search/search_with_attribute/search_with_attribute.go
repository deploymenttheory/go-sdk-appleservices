package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm     = "Taylor Swift"
	mediaType      = "music"
	entityType     = "musicTrack"
	attributeType  = "artistTerm"
	countryCode    = "US"
	resultLimit    = 25
	languageCode   = "en_us"
	explicitFlag   = "No"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Search with Attribute Example ===")
	fmt.Printf("Searching for music tracks by artist '%s' using attribute '%s':\n", searchTerm, attributeType)

	params := itunes_search.NewSearchParams().
		Term(searchTerm).
		Media(mediaType).
		Entity(entityType).
		Attribute(attributeType).
		Country(countryCode).
		Limit(resultLimit).
		Lang(languageCode).
		Explicit(explicitFlag).
		Build()

	response, err := itunesClient.Search(params)
	if err != nil {
		log.Printf("Error searching with attribute: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}