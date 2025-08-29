package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

const (
	searchTerm = "Jack Johnson"
	mediaType  = "music"
	resultLimit = 5
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Basic Search Example ===")
	fmt.Printf("Searching for '%s' %s with limit %d:\n", searchTerm, mediaType, resultLimit)

	params := itunes_search.NewSearchParams().
		Term(searchTerm).
		Media(mediaType).
		Limit(resultLimit).
		Build()

	response, err := itunesClient.Search(params)
	if err != nil {
		log.Printf("Error searching for music: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}