package itunes_search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_buildQueryString(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		params   map[string]string
		expected string
	}{
		{
			name: "single parameter",
			params: map[string]string{
				"term": "Jack Johnson",
			},
			expected: "term=Jack+Johnson",
		},
		{
			name: "multiple parameters",
			params: map[string]string{
				"term":  "Jack Johnson",
				"media": "music",
				"limit": "5",
			},
			expected: "limit=5&media=music&term=Jack+Johnson", // URL parameters are sorted by key
		},
		{
			name: "special characters",
			params: map[string]string{
				"term": "Jack Johnson & Friends",
			},
			expected: "term=Jack+Johnson+%26+Friends",
		},
		{
			name: "empty parameter value",
			params: map[string]string{
				"term": "",
			},
			expected: "term=",
		},
		{
			name:     "empty parameters",
			params:   map[string]string{},
			expected: "",
		},
		{
			name:     "nil parameters",
			params:   nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.buildQueryString(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConstants(t *testing.T) {
	t.Run("BaseSearchURL", func(t *testing.T) {
		assert.Equal(t, "https://itunes.apple.com/search", BaseSearchURL)
	})

	t.Run("BaseLookupURL", func(t *testing.T) {
		assert.Equal(t, "https://itunes.apple.com/lookup", BaseLookupURL)
	})

	t.Run("DefaultLimit", func(t *testing.T) {
		assert.Equal(t, 50, DefaultLimit)
	})

	t.Run("MaxLimit", func(t *testing.T) {
		assert.Equal(t, 200, MaxLimit)
	})
}
