package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm   = "tech talk"
	mediaType    = "podcast"
	entityType   = "podcast"
	countryCode  = "US"
	resultLimit  = 10
	languageCode = "en_us"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Podcast Search Example ===")
	fmt.Printf("Searching for podcasts with term '%s':\n", searchTerm)

	params := itunes_search.NewSearchParams().
		Term(searchTerm).
		Media(mediaType).
		Entity(entityType).
		Country(countryCode).
		Limit(resultLimit).
		Lang(languageCode).
		Build()

	response, err := itunesClient.Search(params)
	if err != nil {
		log.Printf("Error searching for podcasts: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
