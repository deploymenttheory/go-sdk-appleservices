package itunes_search

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadMockResponse loads JSON response from the mocks folder
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) *Client {
	baseClient := client.NewDefaultClient()
	httpmock.ActivateNonDefault(baseClient.HTTP.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)
	return NewClient(baseClient)
}

func TestMock_Search_MusicTrack(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("search_music_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		func(req *http.Request) (*http.Response, error) {
			// Verify expected query parameters
			assert.Equal(t, "Jack Johnson", req.URL.Query().Get("term"))
			assert.Equal(t, "music", req.URL.Query().Get("media"))
			assert.Equal(t, "musicTrack", req.URL.Query().Get("entity"))
			assert.Equal(t, "2", req.URL.Query().Get("limit"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	params := NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Entity("musicTrack").
		Limit(2).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ResultCount)
	assert.Len(t, result.Results, 2)

	// Verify first track
	track1 := result.Results[0]
	assert.Equal(t, "track", track1.WrapperType)
	assert.Equal(t, "song", track1.Kind)
	assert.Equal(t, int64(909253), track1.ArtistID)
	assert.Equal(t, "Jack Johnson", track1.ArtistName)
	assert.Equal(t, "Better Together", track1.TrackName)
	assert.Equal(t, "In Between Dreams", track1.CollectionName)

	// Verify second track
	track2 := result.Results[1]
	assert.Equal(t, "track", track2.WrapperType)
	assert.Equal(t, "song", track2.Kind)
	assert.Equal(t, int64(909253), track2.ArtistID)
	assert.Equal(t, "Jack Johnson", track2.ArtistName)
	assert.Equal(t, "Banana Pancakes", track2.TrackName)

	// Verify that exactly one HTTP call was made
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_Album(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("search_album_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		func(req *http.Request) (*http.Response, error) {
			// Verify expected query parameters
			assert.Equal(t, "Jack Johnson", req.URL.Query().Get("term"))
			assert.Equal(t, "music", req.URL.Query().Get("media"))
			assert.Equal(t, "album", req.URL.Query().Get("entity"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Entity("album").
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ResultCount)
	assert.Len(t, result.Results, 1)

	// Verify album
	album := result.Results[0]
	assert.Equal(t, "collection", album.WrapperType)
	assert.Equal(t, int64(909253), album.ArtistID)
	assert.Equal(t, "Jack Johnson", album.ArtistName)
	assert.Equal(t, "In Between Dreams", album.CollectionName)
	assert.Equal(t, 9.99, album.CollectionPrice)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_Movie(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("search_movie_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "Star Wars", req.URL.Query().Get("term"))
			assert.Equal(t, "movie", req.URL.Query().Get("media"))
			assert.Equal(t, "movie", req.URL.Query().Get("entity"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewSearchParams().
		Term("Star Wars").
		Media("movie").
		Entity("movie").
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ResultCount)
	assert.Len(t, result.Results, 2)

	// Verify first movie
	movie1 := result.Results[0]
	assert.Equal(t, "track", movie1.WrapperType)
	assert.Equal(t, "feature-movie", movie1.Kind)
	assert.Equal(t, "Star Wars: Episode IV - A New Hope", movie1.TrackName)
	assert.Equal(t, "George Lucas", movie1.ArtistName)
	assert.Equal(t, "PG", movie1.ContentAdvisoryRating)
	assert.True(t, movie1.HasITunesExtras)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_App(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("search_app_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "Instagram", req.URL.Query().Get("term"))
			assert.Equal(t, "software", req.URL.Query().Get("media"))
			assert.Equal(t, "software", req.URL.Query().Get("entity"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewSearchParams().
		Term("Instagram").
		Media("software").
		Entity("software").
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ResultCount)
	assert.Len(t, result.Results, 1)

	// Verify app
	app := result.Results[0]
	assert.Equal(t, "software", app.WrapperType)
	assert.Equal(t, "software", app.Kind)
	assert.Equal(t, "Instagram", app.TrackName)
	assert.Equal(t, "Instagram, Inc.", app.ArtistName)
	assert.Equal(t, "Photo & Video", app.PrimaryGenreName)
	assert.Equal(t, "4+", app.ContentAdvisoryRating)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_NoResults(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("search_empty_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		httpmock.NewBytesResponder(200, mockData),
	)


	params := NewSearchParams().
		Term("nonexistentxyz123").
		Media("music").
		Entity("musicTrack").
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ResultCount)
	assert.Empty(t, result.Results)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Register mock responder that returns an error
	httpmock.RegisterResponder("GET", BaseSearchURL,
		httpmock.NewStringResponder(404, "Not Found"),
	)


	params := NewSearchParams().
		Term("test").
		Build()

	result, err := client.Search(params)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "search request failed with status 404")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Search_InvalidJSON(t *testing.T) {
	client := setupMockClient(t)

	// Register mock responder that returns invalid JSON
	httpmock.RegisterResponder("GET", BaseSearchURL,
		httpmock.NewStringResponder(200, "invalid json"),
	)


	params := NewSearchParams().
		Term("test").
		Build()

	result, err := client.Search(params)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal search response")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Lookup_ByID(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("lookup_artist_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseLookupURL,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "909253", req.URL.Query().Get("id"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewLookupParams().
		ID(909253).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ResultCount)
	assert.Len(t, result.Results, 1)

	// Verify artist
	artist := result.Results[0]
	assert.Equal(t, "artist", artist.WrapperType)
	assert.Equal(t, int64(909253), artist.ArtistID)
	assert.Equal(t, "Jack Johnson", artist.ArtistName)
	assert.Equal(t, "Rock", artist.PrimaryGenreName)
	assert.NotEmpty(t, artist.RadioStationURL)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Lookup_MultipleIDs(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("lookup_multiple_artists_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseLookupURL,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "909253,159260351", req.URL.Query().Get("id"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewLookupParams().
		IDs([]int{909253, 159260351}).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ResultCount)
	assert.Len(t, result.Results, 2)

	// Verify first artist
	artist1 := result.Results[0]
	assert.Equal(t, "artist", artist1.WrapperType)
	assert.Equal(t, int64(909253), artist1.ArtistID)
	assert.Equal(t, "Jack Johnson", artist1.ArtistName)

	// Verify second artist
	artist2 := result.Results[1]
	assert.Equal(t, "artist", artist2.WrapperType)
	assert.Equal(t, int64(159260351), artist2.ArtistID)
	assert.Equal(t, "Taylor Swift", artist2.ArtistName)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Lookup_WithEntityFilter(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response
	mockData, err := loadMockResponse("lookup_album_entity_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseLookupURL,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "909253", req.URL.Query().Get("id"))
			assert.Equal(t, "album", req.URL.Query().Get("entity"))

			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	params := NewLookupParams().
		ID(909253).
		Entity("album").
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 6, result.ResultCount)
	assert.Len(t, result.Results, 6)

	// Should include both artist and albums
	hasArtist := false
	hasAlbums := false
	albumCount := 0

	for _, item := range result.Results {
		if item.WrapperType == "artist" {
			hasArtist = true
			assert.Equal(t, int64(909253), item.ArtistID)
			assert.Equal(t, "Jack Johnson", item.ArtistName)
		}
		if item.WrapperType == "collection" {
			hasAlbums = true
			albumCount++
			assert.Equal(t, int64(909253), item.ArtistID)
			assert.Equal(t, "Jack Johnson", item.ArtistName)
			assert.NotEmpty(t, item.CollectionName)
		}
	}

	assert.True(t, hasArtist, "Should include the artist in results")
	assert.True(t, hasAlbums, "Should include albums in results")
	assert.Equal(t, 5, albumCount, "Should have 5 albums")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Lookup_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Register mock responder that returns an error
	httpmock.RegisterResponder("GET", BaseLookupURL,
		httpmock.NewStringResponder(500, "Internal Server Error"),
	)


	params := NewLookupParams().
		ID(909253).
		Build()

	result, err := client.Lookup(params)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "lookup request failed with status 500")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_Lookup_NoResults(t *testing.T) {
	client := setupMockClient(t)

	// Load mock response for no results
	mockData, err := loadMockResponse("search_empty_response.json")
	require.NoError(t, err)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseLookupURL,
		httpmock.NewBytesResponder(200, mockData),
	)


	params := NewLookupParams().
		ID(999999999).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ResultCount)
	assert.Empty(t, result.Results)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_ParameterValidation(t *testing.T) {
	client := setupMockClient(t)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseSearchURL,
		func(req *http.Request) (*http.Response, error) {
			// Verify URL encoding
			query := req.URL.RawQuery
			// Check that the term parameter exists and contains the expected content
			term := req.URL.Query().Get("term")
			assert.Equal(t, "Jack Johnson & Friends", term)
			assert.Contains(t, query, "media=music")
			assert.Contains(t, query, "entity=musicTrack")
			assert.Contains(t, query, "limit=10")
			assert.Contains(t, query, "explicit=No")

			// Return minimal valid response
			mockResponse := SearchResponse{ResultCount: 0, Results: []Result{}}
			responseBytes, _ := json.Marshal(mockResponse)

			resp := httpmock.NewBytesResponse(200, responseBytes)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	// Test parameter encoding and validation
	params := NewSearchParams().
		Term("Jack Johnson & Friends").
		Media("music").
		Entity("musicTrack").
		Limit(10).
		Explicit("No").
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestMock_BuilderChaining(t *testing.T) {
	client := setupMockClient(t)

	// Register mock responder
	httpmock.RegisterResponder("GET", BaseLookupURL,
		func(req *http.Request) (*http.Response, error) {
			// Verify all chained parameters
			assert.Equal(t, "468749", req.URL.Query().Get("amgArtistId"))
			assert.Equal(t, "album", req.URL.Query().Get("entity"))
			assert.Equal(t, "5", req.URL.Query().Get("limit"))
			assert.Equal(t, "recent", req.URL.Query().Get("sort"))
			assert.Equal(t, "US", req.URL.Query().Get("country"))

			// Return minimal valid response
			mockResponse := SearchResponse{ResultCount: 0, Results: []Result{}}
			responseBytes, _ := json.Marshal(mockResponse)

			resp := httpmock.NewBytesResponse(200, responseBytes)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)


	// Test method chaining
	params := NewLookupParams().
		AMGArtistID("468749").
		Entity("album").
		Limit(5).
		Sort("recent").
		Country("US").
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}