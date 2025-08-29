# iTunes Search API Mock Responses

This directory contains JSON mock responses for testing the iTunes Search service without making real API calls.

## Mock Files

- **`search_music_response.json`** - Response for music track searches
- **`search_album_response.json`** - Response for album searches  
- **`search_movie_response.json`** - Response for movie searches
- **`search_app_response.json`** - Response for app/software searches
- **`search_empty_response.json`** - Empty response for no results
- **`lookup_artist_response.json`** - Response for single artist lookup
- **`lookup_multiple_artists_response.json`** - Response for multiple artist lookup
- **`lookup_album_entity_response.json`** - Response for artist lookup with album entity filter

## Usage

The mock responses are used in `mock_test.go` with the `httpmock` library. Each test:

1. Creates a client with httpmock enabled
2. Loads the appropriate JSON response file
3. Registers a mock responder for the iTunes API endpoint
4. Makes the API call and verifies the response

## Example Test Structure

```go
func TestMock_Search_MusicTrack(t *testing.T) {
    client := setupMockClient(t)

    // Load mock response
    mockData, err := loadMockResponse("search_music_response.json")
    require.NoError(t, err)

    // Register mock responder
    httpmock.RegisterResponder("GET", BaseSearchURL,
        func(req *http.Request) (*http.Response, error) {
            // Verify query parameters
            assert.Equal(t, "Jack Johnson", req.URL.Query().Get("term"))
            
            resp := httpmock.NewBytesResponse(200, mockData)
            resp.Header.Set("Content-Type", "application/json")
            return resp, nil
        },
    )

    // Make the API call and test the response
    params := NewSearchParams().Term("Jack Johnson").Build()
    result, err := client.Search(params)
    
    require.NoError(t, err)
    assert.Equal(t, 2, result.ResultCount)
}
```

## Benefits of Mock Testing

- **Fast execution** - No network calls, tests run in milliseconds
- **Reliable** - Tests don't fail due to network issues or API changes
- **Predictable** - Known responses allow precise assertions
- **Offline development** - Tests work without internet connection
- **Parameter validation** - Can verify exact query parameters sent to API

## Running Mock Tests

```bash
# Run only mock tests
go test -v -run TestMock

# Run unit and mock tests (exclude acceptance tests)  
go test -v -short
```