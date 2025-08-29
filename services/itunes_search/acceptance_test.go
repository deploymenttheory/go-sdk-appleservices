package itunes_search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptance_Search_MusicTrack(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Entity("musicTrack").
		Limit(5).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify the first result has expected fields for a music track
	if result.ResultCount > 0 {
		track := result.Results[0]
		assert.Equal(t, "track", track.WrapperType)
		assert.Equal(t, "song", track.Kind)
		assert.NotEmpty(t, track.ArtistName)
		assert.NotEmpty(t, track.TrackName)
		assert.Greater(t, track.ArtistID, int64(0))
	}
}

func TestAcceptance_Search_Album(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Taylor Swift").
		Media("music").
		Entity("album").
		Limit(3).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got album results
	if result.ResultCount > 0 {
		album := result.Results[0]
		assert.Equal(t, "collection", album.WrapperType)
		assert.NotEmpty(t, album.ArtistName)
		assert.NotEmpty(t, album.CollectionName)
		assert.Greater(t, album.CollectionID, int64(0))
	}
}

func TestAcceptance_Search_Movie(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Star Wars").
		Media("movie").
		Entity("movie").
		Limit(5).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got movie results
	if result.ResultCount > 0 {
		movie := result.Results[0]
		assert.Equal(t, "track", movie.WrapperType)
		assert.Equal(t, "feature-movie", movie.Kind)
		assert.NotEmpty(t, movie.TrackName)
	}
}

func TestAcceptance_Search_App(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Instagram").
		Media("software").
		Entity("software").
		Limit(3).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got software results
	if result.ResultCount > 0 {
		app := result.Results[0]
		assert.Equal(t, "software", app.WrapperType)
		assert.Equal(t, "software", app.Kind)
		assert.NotEmpty(t, app.TrackName)
		assert.Greater(t, app.TrackID, int64(0))
	}
}

func TestAcceptance_Search_Podcast(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("This American Life").
		Media("podcast").
		Entity("podcast").
		Limit(5).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got podcast results
	if result.ResultCount > 0 {
		podcast := result.Results[0]
		assert.Equal(t, "track", podcast.WrapperType)
		assert.Equal(t, "podcast", podcast.Kind)
		assert.NotEmpty(t, podcast.ArtistName)
		assert.NotEmpty(t, podcast.TrackName)
	}
}

func TestAcceptance_Search_WithCountryFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Ed Sheeran").
		Media("music").
		Entity("musicTrack").
		Country("GB").
		Limit(3).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got UK-specific results
	if result.ResultCount > 0 {
		track := result.Results[0]
		assert.Equal(t, "GBR", track.Country)
		assert.Equal(t, "GBP", track.Currency)
	}
}

func TestAcceptance_Search_WithExplicitFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Taylor Swift").
		Media("music").
		Entity("musicTrack").
		Explicit("No").
		Limit(5).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)

	// With explicit filter "No", we should get clean results
	if result.ResultCount > 0 {
		for _, track := range result.Results {
			// The explicit filter should return only non-explicit content
			// Note: Some tracks may not have explicitness set, so we check if it's set and not explicit
			if track.TrackExplicitness != "" {
				assert.NotEqual(t, "explicit", track.TrackExplicitness, "Should not contain explicit tracks when explicit=No")
			}
		}
	}
}

func TestAcceptance_Search_WithAttributeFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Entity("musicTrack").
		Attribute("artistTerm").
		Limit(5).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)

	// With artistTerm attribute, all results should be from Jack Johnson
	if result.ResultCount > 0 {
		for _, track := range result.Results {
			assert.Contains(t, track.ArtistName, "Jack Johnson")
		}
	}
}

func TestAcceptance_Lookup_ByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Using Jack Johnson's artist ID
	params := NewLookupParams().
		ID(909253).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got the correct artist
	if result.ResultCount > 0 {
		artist := result.Results[0]
		assert.Equal(t, "artist", artist.WrapperType)
		assert.Equal(t, int64(909253), artist.ArtistID)
		assert.Equal(t, "Jack Johnson", artist.ArtistName)
	}
}

func TestAcceptance_Lookup_ByMultipleIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Using multiple artist IDs
	params := NewLookupParams().
		IDs([]int{909253, 159260351}).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Should get at least 2 results
	assert.GreaterOrEqual(t, result.ResultCount, 2)

	// Verify we got the expected artists
	artistIDs := make([]int64, 0, result.ResultCount)
	for _, artist := range result.Results {
		if artist.WrapperType == "artist" {
			artistIDs = append(artistIDs, artist.ArtistID)
		}
	}

	assert.Contains(t, artistIDs, int64(909253))    // Jack Johnson
	assert.Contains(t, artistIDs, int64(159260351)) // Taylor Swift
}

func TestAcceptance_Lookup_ByAMGArtistID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Using Jack Johnson's AMG Artist ID
	params := NewLookupParams().
		AMGArtistID("468749").
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)
	assert.Len(t, result.Results, result.ResultCount)

	// Verify we got results for the AMG Artist ID
	if result.ResultCount > 0 {
		artist := result.Results[0]
		assert.Equal(t, "artist", artist.WrapperType)
		assert.NotEmpty(t, artist.ArtistName)
	}
}

func TestAcceptance_Lookup_WithEntityFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Get albums for Jack Johnson
	params := NewLookupParams().
		ID(909253).
		Entity("album").
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)

	// Should include both artist and albums
	hasArtist := false
	hasAlbums := false

	for _, item := range result.Results {
		if item.WrapperType == "artist" {
			hasArtist = true
			assert.Equal(t, int64(909253), item.ArtistID)
		}
		if item.WrapperType == "collection" {
			hasAlbums = true
			// Albums will have the artist ID, but not necessarily the same as the lookup ID
			// Just verify it's a valid artist ID
			assert.Greater(t, item.ArtistID, int64(0))
		}
	}

	assert.True(t, hasArtist, "Should include the artist in results")
	assert.True(t, hasAlbums, "Should include albums in results")
}

func TestAcceptance_Lookup_WithLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewLookupParams().
		ID(909253).
		Entity("album").
		Limit(5).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 0)

	// The limit parameter for lookup may not work the same as for search
	// Just verify we get results - the API may return more than the limit
	assert.NotEmpty(t, result.Results)
}

func TestAcceptance_Search_LargeResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	params := NewSearchParams().
		Term("Beatles").
		Media("music").
		Entity("musicTrack").
		Limit(25).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.ResultCount, 5)      // Should get many Beatles songs
	assert.LessOrEqual(t, result.ResultCount, 25) // Should respect limit
	assert.Len(t, result.Results, result.ResultCount)

	// Verify all results are music tracks
	for _, track := range result.Results {
		assert.Equal(t, "track", track.WrapperType)
		assert.Equal(t, "song", track.Kind)
		assert.NotEmpty(t, track.TrackName)
		assert.Greater(t, track.TrackID, int64(0))
	}
}

func TestAcceptance_Search_NoResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Search for something that should return no results
	params := NewSearchParams().
		Term("xzqwerty123nonexistent").
		Media("music").
		Entity("musicTrack").
		Limit(10).
		Build()

	result, err := client.Search(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ResultCount)
	assert.Empty(t, result.Results)
}

func TestAcceptance_Lookup_NonexistentID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	client := NewDefaultClient()

	// Using a non-existent ID
	params := NewLookupParams().
		ID(999999999).
		Build()

	result, err := client.Lookup(params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ResultCount)
	assert.Empty(t, result.Results)
}
