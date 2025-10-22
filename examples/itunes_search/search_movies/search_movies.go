package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm   = "star wars"
	mediaType    = "movie"
	entityType   = "movie"
	countryCode  = "US"
	resultLimit  = 10
	languageCode = "en_us"
	explicitFlag = "Yes"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Movie Search Example ===")
	fmt.Printf("Searching for movies with term '%s':\n", searchTerm)

	params := itunes_search.NewSearchParams().
		Term(searchTerm).
		Media(mediaType).
		Entity(entityType).
		Country(countryCode).
		Limit(resultLimit).
		Lang(languageCode).
		Explicit(explicitFlag).
		Build()

	response, err := itunesClient.Search(params)
	if err != nil {
		log.Printf("Error searching for movies: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
