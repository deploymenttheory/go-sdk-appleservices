package appstore_macos

// AppStoreResponse holds the search results from the iTunes Search API
// for macOS Microsoft applications.
type AppStoreResponse struct {
	// ResultCount is the number of results returned by the iTunes API.
	ResultCount int `json:"resultCount"`

	// Results contains the individual app entries.
	Results []AppEntry `json:"results"`
}

// AppEntry represents a single macOS application in the App Store.
type AppEntry struct {
	// TrackID is the numeric App Store identifier.
	TrackID int `json:"trackId"`

	// TrackName is the application's display name in the App Store.
	TrackName string `json:"trackName"`

	// BundleID is the macOS bundle identifier (e.g. "com.microsoft.Word").
	BundleID string `json:"bundleId"`

	// Version is the current App Store version string (e.g. "16.108.0").
	Version string `json:"version"`

	// CurrentVersionReleaseDate is the ISO 8601 release date of the current version.
	CurrentVersionReleaseDate string `json:"currentVersionReleaseDate"`

	// MinimumOsVersion is the minimum macOS version required.
	MinimumOsVersion string `json:"minimumOsVersion"`

	// ReleaseNotes contains the release notes for the current version.
	ReleaseNotes string `json:"releaseNotes"`

	// ArtworkUrl512 is the URL of the 512×512 app icon.
	ArtworkUrl512 string `json:"artworkUrl512"`

	// TrackViewURL is the App Store URL for this application.
	TrackViewURL string `json:"trackViewUrl"`

	// Price is the App Store price (0.0 for free apps).
	Price float64 `json:"price"`

	// SellerName is the developer name as shown in the App Store.
	SellerName string `json:"sellerName"`
}
