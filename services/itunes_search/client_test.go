package itunes_search

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-apple/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	baseClient := &client.Client{}
	itunesClient := NewClient(baseClient)

	require.NotNil(t, itunesClient)
	assert.Equal(t, baseClient, itunesClient.baseClient)
	assert.Equal(t, baseClient.Logger, itunesClient.logger)
}

func TestNewDefaultClient(t *testing.T) {
	itunesClient := NewDefaultClient()

	require.NotNil(t, itunesClient)
	require.NotNil(t, itunesClient.baseClient)
	require.NotNil(t, itunesClient.logger)
}

func TestClient_Search_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty parameters",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "at least one search parameter is required",
		},
		{
			name:    "nil parameters",
			params:  nil,
			wantErr: true,
			errMsg:  "at least one search parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewDefaultClient()
			result, err := client.Search(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestClient_Lookup_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty parameters",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "at least one lookup parameter is required",
		},
		{
			name:    "nil parameters",
			params:  nil,
			wantErr: true,
			errMsg:  "at least one lookup parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewDefaultClient()
			result, err := client.Lookup(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestClient_Search_MockServer(t *testing.T) {
	// Mock iTunes API response
	mockResponse := SearchResponse{
		ResultCount: 2,
		Results: []Result{
			{
				WrapperType:       "track",
				Kind:              "song",
				ArtistID:          909253,
				ArtistName:        "Jack Johnson",
				TrackName:         "Imagine",
				CollectionName:    "Test Album",
				PrimaryGenreName:  "Rock",
				Country:           "USA",
				Currency:          "USD",
				TrackPrice:        1.29,
				CollectionPrice:   9.99,
				ReleaseDate:       "2007-06-11T12:00:00Z",
			},
			{
				WrapperType:       "track",
				Kind:              "song",
				ArtistID:          909253,
				ArtistName:        "Jack Johnson",
				TrackName:         "Flake",
				CollectionName:    "Test Album 2",
				PrimaryGenreName:  "Rock",
				Country:           "USA",
				Currency:          "USD",
				TrackPrice:        0.99,
				CollectionPrice:   8.99,
				ReleaseDate:       "2001-02-06T08:00:00Z",
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		assert.Equal(t, "Jack+Johnson", r.URL.Query().Get("term"))
		assert.Equal(t, "music", r.URL.Query().Get("media"))
		assert.Equal(t, "2", r.URL.Query().Get("limit"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// This test is limited due to the tight coupling with the baseClient
	// In a real scenario, you would need to mock the HTTP client properly
	t.Skip("Skipping mock server test due to tight coupling with baseClient HTTP implementation")
}

func TestSearchResponse_JSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"resultCount": 2,
		"results": [
			{
				"wrapperType": "track",
				"kind": "song",
				"artistId": 909253,
				"artistName": "Jack Johnson",
				"trackName": "Imagine",
				"primaryGenreName": "Rock",
				"trackPrice": 1.29
			},
			{
				"wrapperType": "track", 
				"kind": "song",
				"artistId": 909253,
				"artistName": "Jack Johnson",
				"trackName": "Flake",
				"primaryGenreName": "Rock",
				"trackPrice": 0.99
			}
		]
	}`

	var response SearchResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	require.NoError(t, err)
	assert.Equal(t, 2, response.ResultCount)
	assert.Len(t, response.Results, 2)

	// Test first result
	result1 := response.Results[0]
	assert.Equal(t, "track", result1.WrapperType)
	assert.Equal(t, "song", result1.Kind)
	assert.Equal(t, int64(909253), result1.ArtistID)
	assert.Equal(t, "Jack Johnson", result1.ArtistName)
	assert.Equal(t, "Imagine", result1.TrackName)
	assert.Equal(t, "Rock", result1.PrimaryGenreName)
	assert.Equal(t, 1.29, result1.TrackPrice)

	// Test second result
	result2 := response.Results[1]
	assert.Equal(t, "track", result2.WrapperType)
	assert.Equal(t, "song", result2.Kind)
	assert.Equal(t, int64(909253), result2.ArtistID)
	assert.Equal(t, "Jack Johnson", result2.ArtistName)
	assert.Equal(t, "Flake", result2.TrackName)
	assert.Equal(t, "Rock", result2.PrimaryGenreName)
	assert.Equal(t, 0.99, result2.TrackPrice)
}

func TestResult_AllFields(t *testing.T) {
	jsonData := `{
		"wrapperType": "track",
		"kind": "song",
		"artistId": 909253,
		"collectionId": 255144028,
		"trackId": 255145362,
		"artistName": "Jack Johnson",
		"collectionName": "Test Album",
		"trackName": "Imagine",
		"collectionCensoredName": "Test Album",
		"trackCensoredName": "Imagine",
		"artistViewUrl": "https://music.apple.com/us/artist/jack-johnson/909253",
		"collectionViewUrl": "https://music.apple.com/us/album/test-album/255144028",
		"trackViewUrl": "https://music.apple.com/us/album/imagine/255144028",
		"previewUrl": "https://audio-ssl.itunes.apple.com/preview.m4a",
		"artworkUrl30": "https://is1-ssl.mzstatic.com/image/30x30bb.jpg",
		"artworkUrl60": "https://is1-ssl.mzstatic.com/image/60x60bb.jpg",
		"artworkUrl100": "https://is1-ssl.mzstatic.com/image/100x100bb.jpg",
		"collectionPrice": 19.99,
		"trackPrice": 1.29,
		"releaseDate": "2007-06-11T12:00:00Z",
		"collectionExplicitness": "notExplicit",
		"trackExplicitness": "notExplicit",
		"discCount": 1,
		"discNumber": 1,
		"trackCount": 34,
		"trackNumber": 15,
		"trackTimeMillis": 219080,
		"country": "USA",
		"currency": "USD",
		"primaryGenreName": "Rock",
		"isStreamable": true
	}`

	var result Result
	err := json.Unmarshal([]byte(jsonData), &result)

	require.NoError(t, err)
	
	// Test all major fields
	assert.Equal(t, "track", result.WrapperType)
	assert.Equal(t, "song", result.Kind)
	assert.Equal(t, int64(909253), result.ArtistID)
	assert.Equal(t, int64(255144028), result.CollectionID)
	assert.Equal(t, int64(255145362), result.TrackID)
	assert.Equal(t, "Jack Johnson", result.ArtistName)
	assert.Equal(t, "Test Album", result.CollectionName)
	assert.Equal(t, "Imagine", result.TrackName)
	assert.Equal(t, "Test Album", result.CollectionCensoredName)
	assert.Equal(t, "Imagine", result.TrackCensoredName)
	assert.Equal(t, "https://music.apple.com/us/artist/jack-johnson/909253", result.ArtistViewURL)
	assert.Equal(t, "https://music.apple.com/us/album/test-album/255144028", result.CollectionViewURL)
	assert.Equal(t, "https://music.apple.com/us/album/imagine/255144028", result.TrackViewURL)
	assert.Equal(t, "https://audio-ssl.itunes.apple.com/preview.m4a", result.PreviewURL)
	assert.Equal(t, "https://is1-ssl.mzstatic.com/image/30x30bb.jpg", result.ArtworkURL30)
	assert.Equal(t, "https://is1-ssl.mzstatic.com/image/60x60bb.jpg", result.ArtworkURL60)
	assert.Equal(t, "https://is1-ssl.mzstatic.com/image/100x100bb.jpg", result.ArtworkURL100)
	assert.Equal(t, 19.99, result.CollectionPrice)
	assert.Equal(t, 1.29, result.TrackPrice)
	assert.Equal(t, "2007-06-11T12:00:00Z", result.ReleaseDate)
	assert.Equal(t, "notExplicit", result.CollectionExplicitness)
	assert.Equal(t, "notExplicit", result.TrackExplicitness)
	assert.Equal(t, 1, result.DiscCount)
	assert.Equal(t, 1, result.DiscNumber)
	assert.Equal(t, 34, result.TrackCount)
	assert.Equal(t, 15, result.TrackNumber)
	assert.Equal(t, 219080, result.TrackTimeMillis)
	assert.Equal(t, "USA", result.Country)
	assert.Equal(t, "USD", result.Currency)
	assert.Equal(t, "Rock", result.PrimaryGenreName)
	assert.Equal(t, true, result.IsStreamable)
}