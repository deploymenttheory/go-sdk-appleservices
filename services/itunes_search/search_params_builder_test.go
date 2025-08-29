package itunes_search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSearchParams(t *testing.T) {
	builder := NewSearchParams()
	
	require.NotNil(t, builder)
	assert.NotNil(t, builder.params)
	assert.Equal(t, 0, len(builder.params))
}

func TestSearchParamsBuilder_Term(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		expected map[string]string
	}{
		{
			name: "valid term",
			term: "Jack Johnson",
			expected: map[string]string{
				"term": "Jack Johnson",
			},
		},
		{
			name:     "empty term",
			term:     "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Term(tt.term)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Country(t *testing.T) {
	tests := []struct {
		name     string
		country  string
		expected map[string]string
	}{
		{
			name:    "valid country",
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
			builder := NewSearchParams().Country(tt.country)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Media(t *testing.T) {
	tests := []struct {
		name     string
		media    string
		expected map[string]string
	}{
		{
			name:  "music media",
			media: "music",
			expected: map[string]string{
				"media": "music",
			},
		},
		{
			name:  "movie media",
			media: "movie",
			expected: map[string]string{
				"media": "movie",
			},
		},
		{
			name:     "empty media",
			media:    "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Media(tt.media)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Entity(t *testing.T) {
	tests := []struct {
		name     string
		entity   string
		expected map[string]string
	}{
		{
			name:   "musicTrack entity",
			entity: "musicTrack",
			expected: map[string]string{
				"entity": "musicTrack",
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
			builder := NewSearchParams().Entity(tt.entity)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Attribute(t *testing.T) {
	tests := []struct {
		name      string
		attribute string
		expected  map[string]string
	}{
		{
			name:      "artistTerm attribute",
			attribute: "artistTerm",
			expected: map[string]string{
				"attribute": "artistTerm",
			},
		},
		{
			name:      "empty attribute",
			attribute: "",
			expected:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Attribute(tt.attribute)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Limit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected map[string]string
	}{
		{
			name:  "valid limit",
			limit: 25,
			expected: map[string]string{
				"limit": "25",
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
			limit:    -10,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Limit(tt.limit)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Lang(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected map[string]string
	}{
		{
			name: "english language",
			lang: "en_us",
			expected: map[string]string{
				"lang": "en_us",
			},
		},
		{
			name: "japanese language",
			lang: "ja_jp",
			expected: map[string]string{
				"lang": "ja_jp",
			},
		},
		{
			name:     "empty language",
			lang:     "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Lang(tt.lang)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Version(t *testing.T) {
	tests := []struct {
		name     string
		version  int
		expected map[string]string
	}{
		{
			name:    "version 2",
			version: 2,
			expected: map[string]string{
				"version": "2",
			},
		},
		{
			name:     "zero version",
			version:  0,
			expected: map[string]string{},
		},
		{
			name:     "negative version",
			version:  -1,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Version(tt.version)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Explicit(t *testing.T) {
	tests := []struct {
		name     string
		explicit string
		expected map[string]string
	}{
		{
			name:     "explicit yes",
			explicit: "Yes",
			expected: map[string]string{
				"explicit": "Yes",
			},
		},
		{
			name:     "explicit no",
			explicit: "No",
			expected: map[string]string{
				"explicit": "No",
			},
		},
		{
			name:     "empty explicit",
			explicit: "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Explicit(tt.explicit)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_Callback(t *testing.T) {
	tests := []struct {
		name     string
		callback string
		expected map[string]string
	}{
		{
			name:     "valid callback",
			callback: "wsSearchCB",
			expected: map[string]string{
				"callback": "wsSearchCB",
			},
		},
		{
			name:     "empty callback",
			callback: "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchParams().Callback(tt.callback)
			result := builder.Build()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchParamsBuilder_ChainedCalls(t *testing.T) {
	builder := NewSearchParams().
		Term("Jack Johnson").
		Media("music").
		Entity("musicTrack").
		Country("US").
		Limit(10).
		Lang("en_us").
		Explicit("No")

	result := builder.Build()

	expected := map[string]string{
		"term":     "Jack Johnson",
		"media":    "music",
		"entity":   "musicTrack",
		"country":  "US",
		"limit":    "10",
		"lang":     "en_us",
		"explicit": "No",
	}

	assert.Equal(t, expected, result)
}

func TestSearchParamsBuilder_OverwriteValues(t *testing.T) {
	builder := NewSearchParams().
		Term("Jack Johnson").
		Term("Taylor Swift")

	result := builder.Build()

	expected := map[string]string{
		"term": "Taylor Swift",
	}

	assert.Equal(t, expected, result)
}