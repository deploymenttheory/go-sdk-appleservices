package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm   = "宇多田ヒカル"
	mediaType    = "music"
	entityType   = "musicTrack"
	countryCode  = "JP"
	resultLimit  = 10
	languageCode = "ja_jp"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Japanese Language Search Example ===")
	fmt.Printf("Searching for music with Japanese language setting (term: '%s'):\n", searchTerm)

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
		log.Printf("Error searching with Japanese language: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
