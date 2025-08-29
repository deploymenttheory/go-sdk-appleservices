package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm   = "photo editor"
	mediaType    = "software"
	entityType   = "software"
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

	fmt.Println("=== Software Search Example ===")
	fmt.Printf("Searching for software with term '%s':\n", searchTerm)

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
		log.Printf("Error searching for software: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}