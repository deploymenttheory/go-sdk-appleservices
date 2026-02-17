package client

import (
	"testing"
)

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name  string
		links *Links
		want  bool
	}{
		{
			name:  "Nil links",
			links: nil,
			want:  false,
		},
		{
			name: "Empty next URL",
			links: &Links{
				Next: "",
			},
			want: false,
		},
		{
			name: "Valid next URL",
			links: &Links{
				Next: "https://api-business.apple.com/v1/orgDevices?cursor=abc123",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasNextPage(tt.links)
			if got != tt.want {
				t.Errorf("HasNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasPrevPage(t *testing.T) {
	tests := []struct {
		name  string
		links *Links
		want  bool
	}{
		{
			name:  "Nil links",
			links: nil,
			want:  false,
		},
		{
			name: "Empty prev URL",
			links: &Links{
				Prev: "",
			},
			want: false,
		},
		{
			name: "Valid prev URL",
			links: &Links{
				Prev: "https://api-business.apple.com/v1/orgDevices?cursor=xyz789",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPrevPage(tt.links)
			if got != tt.want {
				t.Errorf("HasPrevPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractParamsFromURL(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		want    map[string]string
		wantErr bool
	}{
		{
			name:   "Single parameter",
			urlStr: "https://api-business.apple.com/v1/orgDevices?cursor=abc123",
			want: map[string]string{
				"cursor": "abc123",
			},
			wantErr: false,
		},
		{
			name:   "Multiple parameters",
			urlStr: "https://api-business.apple.com/v1/orgDevices?cursor=abc123&limit=100",
			want: map[string]string{
				"cursor": "abc123",
				"limit":  "100",
			},
			wantErr: false,
		},
		{
			name:    "No parameters",
			urlStr:  "https://api-business.apple.com/v1/orgDevices",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			urlStr:  "://invalid-url",
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Parameter with special characters",
			urlStr: "https://api-business.apple.com/v1/orgDevices?fields[orgDevices]=serialNumber,deviceModel",
			want: map[string]string{
				"fields[orgDevices]": "serialNumber,deviceModel",
			},
			wantErr: false,
		},
		{
			name:   "Multiple values for same key (takes first)",
			urlStr: "https://api-business.apple.com/v1/orgDevices?key=value1&key=value2",
			want: map[string]string{
				"key": "value1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractParamsFromURL(tt.urlStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractParamsFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("extractParamsFromURL() got %d params, want %d", len(got), len(tt.want))
			}

			for key, wantValue := range tt.want {
				gotValue, ok := got[key]
				if !ok {
					t.Errorf("extractParamsFromURL() missing key %q", key)
					continue
				}
				if gotValue != wantValue {
					t.Errorf("extractParamsFromURL()[%q] = %v, want %v", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestPaginationOptions_AddToQueryBuilder(t *testing.T) {
	tests := []struct {
		name string
		opts *PaginationOptions
		want map[string]string
	}{
		{
			name: "Nil options",
			opts: nil,
			want: map[string]string{},
		},
		{
			name: "Empty options",
			opts: &PaginationOptions{},
			want: map[string]string{},
		},
		{
			name: "Limit only",
			opts: &PaginationOptions{
				Limit: 100,
			},
			want: map[string]string{
				"limit": "100",
			},
		},
		{
			name: "Cursor only",
			opts: &PaginationOptions{
				Cursor: "abc123",
			},
			want: map[string]string{
				"cursor": "abc123",
			},
		},
		{
			name: "Both limit and cursor",
			opts: &PaginationOptions{
				Limit:  50,
				Cursor: "xyz789",
			},
			want: map[string]string{
				"limit":  "50",
				"cursor": "xyz789",
			},
		},
		{
			name: "Zero limit",
			opts: &PaginationOptions{
				Limit:  0,
				Cursor: "abc",
			},
			want: map[string]string{
				"cursor": "abc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.opts.AddToQueryBuilder(qb)

			got := qb.Build()

			if len(got) != len(tt.want) {
				t.Errorf("AddToQueryBuilder() got %d params, want %d", len(got), len(tt.want))
			}

			for key, wantValue := range tt.want {
				gotValue, ok := got[key]
				if !ok {
					t.Errorf("AddToQueryBuilder() missing key %q", key)
					continue
				}
				if gotValue != wantValue {
					t.Errorf("AddToQueryBuilder()[%q] = %v, want %v", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestLinks_AllFields(t *testing.T) {
	links := &Links{
		Self:  "https://api.example.com/v1/resource",
		First: "https://api.example.com/v1/resource?cursor=first",
		Next:  "https://api.example.com/v1/resource?cursor=next",
		Prev:  "https://api.example.com/v1/resource?cursor=prev",
		Last:  "https://api.example.com/v1/resource?cursor=last",
	}

	if links.Self == "" {
		t.Error("Self field is empty")
	}
	if links.First == "" {
		t.Error("First field is empty")
	}
	if links.Next == "" {
		t.Error("Next field is empty")
	}
	if links.Prev == "" {
		t.Error("Prev field is empty")
	}
	if links.Last == "" {
		t.Error("Last field is empty")
	}
}

func TestMeta_Paging(t *testing.T) {
	meta := &Meta{
		Paging: &Paging{
			Total:      1000,
			Limit:      100,
			NextCursor: "abc123",
		},
	}

	if meta.Paging == nil {
		t.Fatal("Paging is nil")
	}

	if meta.Paging.Total != 1000 {
		t.Errorf("Paging.Total = %d, want 1000", meta.Paging.Total)
	}

	if meta.Paging.Limit != 100 {
		t.Errorf("Paging.Limit = %d, want 100", meta.Paging.Limit)
	}

	if meta.Paging.NextCursor != "abc123" {
		t.Errorf("Paging.NextCursor = %v, want 'abc123'", meta.Paging.NextCursor)
	}
}

func TestMeta_NilPaging(t *testing.T) {
	meta := &Meta{
		Paging: nil,
	}

	if meta.Paging != nil {
		t.Error("Expected nil Paging")
	}
}

func TestPaginationOptions_Defaults(t *testing.T) {
	opts := &PaginationOptions{}

	if opts.Limit != 0 {
		t.Errorf("Default Limit = %d, want 0", opts.Limit)
	}

	if opts.Cursor != "" {
		t.Errorf("Default Cursor = %q, want empty string", opts.Cursor)
	}
}

func TestExtractParamsFromURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Empty string",
			urlStr:  "",
			wantErr: false, // URL parsing allows empty string
		},
		{
			name:    "Just query params",
			urlStr:  "?cursor=abc&limit=10",
			wantErr: false,
		},
		{
			name:    "URL with fragment",
			urlStr:  "https://api.example.com/resource?cursor=abc#fragment",
			wantErr: false,
		},
		{
			name:    "URL with empty query value",
			urlStr:  "https://api.example.com/resource?cursor=",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := extractParamsFromURL(tt.urlStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractParamsFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
