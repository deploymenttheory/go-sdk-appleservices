package itunes_search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLookupParams(t *testing.T) {
	builder := NewLookupParams()
	
	require.NotNil(t, builder)
	assert.NotNil(t, builder.params)
	assert.Equal(t, 0, len(builder.params))
}

func TestLookupParamsBuilder_ID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		expected map[string]string
	}{
		{
			name: "valid id",
			id:   909253,
			expected: map[string]string{
				"id": "909253",
			},
		},
		{
			name:     "zero id",
			id:       0,
			expected: map[string]string{},
		},
		{
			name:     "negative id",
			id:       -1,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().ID(tt.id)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_IDs(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int
		expected map[string]string
	}{
		{
			name: "multiple valid ids",
			ids:  []int{909253, 123456, 789012},
			expected: map[string]string{
				"id": "909253,123456,789012",
			},
		},
		{
			name: "single id",
			ids:  []int{909253},
			expected: map[string]string{
				"id": "909253",
			},
		},
		{
			name:     "empty ids",
			ids:      []int{},
			expected: map[string]string{},
		},
		{
			name:     "nil ids",
			ids:      nil,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().IDs(tt.ids)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_UPC(t *testing.T) {
	tests := []struct {
		name     string
		upc      string
		expected map[string]string
	}{
		{
			name: "valid upc",
			upc:  "720642462928",
			expected: map[string]string{
				"upc": "720642462928",
			},
		},
		{
			name:     "empty upc",
			upc:      "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().UPC(tt.upc)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_EAN(t *testing.T) {
	tests := []struct {
		name     string
		ean      string
		expected map[string]string
	}{
		{
			name: "valid ean",
			ean:  "1234567890123",
			expected: map[string]string{
				"ean": "1234567890123",
			},
		},
		{
			name:     "empty ean",
			ean:      "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().EAN(tt.ean)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_ISRC(t *testing.T) {
	tests := []struct {
		name     string
		isrc     string
		expected map[string]string
	}{
		{
			name: "valid isrc",
			isrc: "USRC12345678",
			expected: map[string]string{
				"isrc": "USRC12345678",
			},
		},
		{
			name:     "empty isrc",
			isrc:     "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().ISRC(tt.isrc)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_ISBN(t *testing.T) {
	tests := []struct {
		name     string
		isbn     string
		expected map[string]string
	}{
		{
			name: "valid isbn",
			isbn: "9780316069359",
			expected: map[string]string{
				"isbn": "9780316069359",
			},
		},
		{
			name:     "empty isbn",
			isbn:     "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().ISBN(tt.isbn)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_AMGArtistID(t *testing.T) {
	tests := []struct {
		name        string
		amgArtistID string
		expected    map[string]string
	}{
		{
			name:        "valid amg artist id",
			amgArtistID: "468749",
			expected: map[string]string{
				"amgArtistId": "468749",
			},
		},
		{
			name:        "empty amg artist id",
			amgArtistID: "",
			expected:    map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().AMGArtistID(tt.amgArtistID)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_AMGArtistIDs(t *testing.T) {
	tests := []struct {
		name         string
		amgArtistIDs []string
		expected     map[string]string
	}{
		{
			name:         "multiple amg artist ids",
			amgArtistIDs: []string{"468749", "5723"},
			expected: map[string]string{
				"amgArtistId": "468749,5723",
			},
		},
		{
			name:         "single amg artist id",
			amgArtistIDs: []string{"468749"},
			expected: map[string]string{
				"amgArtistId": "468749",
			},
		},
		{
			name:         "empty amg artist ids",
			amgArtistIDs: []string{},
			expected:     map[string]string{},
		},
		{
			name:         "nil amg artist ids",
			amgArtistIDs: nil,
			expected:     map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().AMGArtistIDs(tt.amgArtistIDs)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_AMGAlbumID(t *testing.T) {
	tests := []struct {
		name       string
		amgAlbumID string
		expected   map[string]string
	}{
		{
			name:       "valid amg album id",
			amgAlbumID: "15175",
			expected: map[string]string{
				"amgAlbumId": "15175",
			},
		},
		{
			name:       "empty amg album id",
			amgAlbumID: "",
			expected:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().AMGAlbumID(tt.amgAlbumID)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_AMGAlbumIDs(t *testing.T) {
	tests := []struct {
		name        string
		amgAlbumIDs []string
		expected    map[string]string
	}{
		{
			name:        "multiple amg album ids",
			amgAlbumIDs: []string{"15175", "15176", "15177"},
			expected: map[string]string{
				"amgAlbumId": "15175,15176,15177",
			},
		},
		{
			name:        "single amg album id",
			amgAlbumIDs: []string{"15175"},
			expected: map[string]string{
				"amgAlbumId": "15175",
			},
		},
		{
			name:        "empty amg album ids",
			amgAlbumIDs: []string{},
			expected:    map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().AMGAlbumIDs(tt.amgAlbumIDs)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_AMGVideoID(t *testing.T) {
	tests := []struct {
		name       string
		amgVideoID string
		expected   map[string]string
	}{
		{
			name:       "valid amg video id",
			amgVideoID: "17120",
			expected: map[string]string{
				"amgVideoId": "17120",
			},
		},
		{
			name:       "empty amg video id",
			amgVideoID: "",
			expected:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().AMGVideoID(tt.amgVideoID)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_Entity(t *testing.T) {
	tests := []struct {
		name     string
		entity   string
		expected map[string]string
	}{
		{
			name:   "album entity",
			entity: "album",
			expected: map[string]string{
				"entity": "album",
			},
		},
		{
			name:   "song entity",
			entity: "song",
			expected: map[string]string{
				"entity": "song",
			},
		},
		{
			name:     "empty entity",
			entity:   "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().Entity(tt.entity)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_Limit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected map[string]string
	}{
		{
			name:  "valid limit",
			limit: 10,
			expected: map[string]string{
				"limit": "10",
			},
		},
		{
			name:  "limit exceeds maximum",
			limit: 300,
			expected: map[string]string{
				"limit": "200",
			},
		},
		{
			name:     "zero limit",
			limit:    0,
			expected: map[string]string{},
		},
		{
			name:     "negative limit",
			limit:    -5,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().Limit(tt.limit)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_Sort(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected map[string]string
	}{
		{
			name: "recent sort",
			sort: "recent",
			expected: map[string]string{
				"sort": "recent",
			},
		},
		{
			name:     "empty sort",
			sort:     "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().Sort(tt.sort)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_Country(t *testing.T) {
	tests := []struct {
		name     string
		country  string
		expected map[string]string
	}{
		{
			name:    "US country",
			country: "US",
			expected: map[string]string{
				"country": "US",
			},
		},
		{
			name:     "empty country",
			country:  "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewLookupParams().Country(tt.country)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupParamsBuilder_ChainedCalls(t *testing.T) {
	builder := NewLookupParams().
		AMGArtistID("468749").
		Entity("album").
		Limit(5).
		Sort("recent").
		Country("US")

	result := builder.Build()

	expected := map[string]string{
		"amgArtistId": "468749",
		"entity":     "album",
		"limit":      "5",
		"sort":       "recent",
		"country":    "US",
	}

	assert.Equal(t, expected, result)
}

func TestLookupParamsBuilder_MultipleIDs(t *testing.T) {
	builder := NewLookupParams().
		ID(909253).
		AMGArtistID("468749").
		UPC("720642462928")

	result := builder.Build()

	expected := map[string]string{
		"id":          "909253",
		"amgArtistId": "468749",
		"upc":         "720642462928",
	}

	assert.Equal(t, expected, result)
}