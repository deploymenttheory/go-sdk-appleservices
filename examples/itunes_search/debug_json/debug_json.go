package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/deploymenttheory/go-api-sdk-apple/client/itunes"
	"github.com/deploymenttheory/go-api-sdk-apple/services/itunes_search"
)

func main() {
	baseClient := client.NewClient(client.Config{
		Debug: true,
	})
	defer baseClient.Close()

	itunesClient := itunes_search.NewClient(baseClient)

	fmt.Println("=== Debug JSON Parsing ===")

	params := itunes_search.NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Limit(2).
		Build()

	// Test with our normal search method
	response, err := itunesClient.Search(params)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Parsed Response: %s\n", string(jsonData))

	// Test with raw JSON string
	rawJSON := `{
 "resultCount":2,
 "results": [
{"wrapperType":"track", "kind":"song", "artistId":909253, "collectionId":255144028, "trackId":255145362, "artistName":"Jack Johnson", "collectionName":"Test Album", "trackName":"Test Song", "releaseDate":"2007-06-11T12:00:00Z", "primaryGenreName":"Rock"}
]
}`

	var testResponse itunes_search.SearchResponse
	err = json.Unmarshal([]byte(rawJSON), &testResponse)
	if err != nil {
		log.Printf("JSON Parse Error: %v", err)
		return
	}

	testJSON, _ := json.MarshalIndent(testResponse, "", "  ")
	fmt.Printf("Test Parse Response: %s\n", string(testJSON))
}
