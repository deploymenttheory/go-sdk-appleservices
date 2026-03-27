package search

// SearchOptions configures a request to the iTunes Search endpoint.
// At minimum Term must be set.
type SearchOptions struct {
	// Term is the URL-encoded text string you want to search for.
	Term string
	// Country is the two-letter ISO 3166-1 country code (e.g. "us", "gb").
	Country string
	// Media is the media type to search for (e.g. "music", "movie", "software").
	// Defaults to "all" when empty.
	Media string
	// Entity is the type of results to return (e.g. "song", "album", "musicArtist").
	Entity string
	// Attribute is the attribute to search within (e.g. "artistTerm", "titleTerm").
	Attribute string
	// Callback is the name of the JavaScript callback function for JSONP.
	Callback string
	// Limit is the number of results to return. Maximum is 200. Defaults to 50.
	Limit int
	// Lang is the language to use (e.g. "en_us", "ja_jp"). Defaults to "en_us".
	Lang string
	// Version is the iTunes Search API version (1 or 2). Defaults to 2.
	Version int
	// Explicit controls whether explicit content is returned ("Yes" / "No").
	Explicit string
}

// LookupOptions configures a request to the iTunes Lookup endpoint.
// At least one identifier field must be set.
type LookupOptions struct {
	// ID looks up a single item by its iTunes ID.
	ID int
	// IDs looks up multiple items by their iTunes IDs (comma-joined).
	IDs []int
	// UPC looks up an album or video by its UPC barcode.
	UPC string
	// EAN looks up an album by its European Article Number.
	EAN string
	// ISRC looks up a song by its International Standard Recording Code.
	ISRC string
	// ISBN looks up a book by its International Standard Book Number.
	ISBN string
	// AMGArtistID looks up an artist by their All Music Guide artist ID.
	AMGArtistID string
	// AMGArtistIDs looks up multiple artists by their AMG artist IDs.
	AMGArtistIDs []string
	// AMGAlbumID looks up an album by its AMG album ID.
	AMGAlbumID string
	// AMGAlbumIDs looks up multiple albums by their AMG album IDs.
	AMGAlbumIDs []string
	// AMGVideoID looks up a video by its AMG video ID.
	AMGVideoID string
	// Entity is the type of related content to return alongside the lookup result.
	Entity string
	// Limit is the maximum number of results to return. Maximum is 200.
	Limit int
	// Sort sorts results by a field (e.g. "recent").
	Sort string
	// Country is the two-letter ISO 3166-1 country code to scope the lookup.
	Country string
}

// SearchResponse is the top-level response from both the Search and Lookup endpoints.
type SearchResponse struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

// Result represents a single item returned by the iTunes Search or Lookup API.
type Result struct {
	WrapperType                        string   `json:"wrapperType,omitempty"`
	Kind                               string   `json:"kind,omitempty"`
	ArtistID                           int64    `json:"artistId,omitempty"`
	CollectionID                       int64    `json:"collectionId,omitempty"`
	TrackID                            int64    `json:"trackId,omitempty"`
	ArtistName                         string   `json:"artistName,omitempty"`
	CollectionName                     string   `json:"collectionName,omitempty"`
	TrackName                          string   `json:"trackName,omitempty"`
	CollectionCensoredName             string   `json:"collectionCensoredName,omitempty"`
	TrackCensoredName                  string   `json:"trackCensoredName,omitempty"`
	ArtistViewURL                      string   `json:"artistViewUrl,omitempty"`
	CollectionViewURL                  string   `json:"collectionViewUrl,omitempty"`
	TrackViewURL                       string   `json:"trackViewUrl,omitempty"`
	PreviewURL                         string   `json:"previewUrl,omitempty"`
	ArtworkURL30                       string   `json:"artworkUrl30,omitempty"`
	ArtworkURL60                       string   `json:"artworkUrl60,omitempty"`
	ArtworkURL100                      string   `json:"artworkUrl100,omitempty"`
	CollectionPrice                    float64  `json:"collectionPrice,omitempty"`
	TrackPrice                         float64  `json:"trackPrice,omitempty"`
	ReleaseDate                        string   `json:"releaseDate,omitempty"`
	CollectionExplicitness             string   `json:"collectionExplicitness,omitempty"`
	TrackExplicitness                  string   `json:"trackExplicitness,omitempty"`
	DiscCount                          int      `json:"discCount,omitempty"`
	DiscNumber                         int      `json:"discNumber,omitempty"`
	TrackCount                         int      `json:"trackCount,omitempty"`
	TrackNumber                        int      `json:"trackNumber,omitempty"`
	TrackTimeMillis                    int      `json:"trackTimeMillis,omitempty"`
	Country                            string   `json:"country,omitempty"`
	Currency                           string   `json:"currency,omitempty"`
	PrimaryGenreName                   string   `json:"primaryGenreName,omitempty"`
	RadioStationURL                    string   `json:"radioStationUrl,omitempty"`
	IsStreamable                       bool     `json:"isStreamable,omitempty"`
	ContentAdvisoryRating              string   `json:"contentAdvisoryRating,omitempty"`
	CollectionHDPrice                  float64  `json:"collectionHdPrice,omitempty"`
	TrackHDPrice                       float64  `json:"trackHdPrice,omitempty"`
	TrackRentalPrice                   float64  `json:"trackRentalPrice,omitempty"`
	CollectionHDRentalPrice            float64  `json:"collectionHdRentalPrice,omitempty"`
	TrackHDRentalPrice                 float64  `json:"trackHdRentalPrice,omitempty"`
	LongDescription                    string   `json:"longDescription,omitempty"`
	ShortDescription                   string   `json:"shortDescription,omitempty"`
	HasITunesExtras                    bool     `json:"hasITunesExtras,omitempty"`
	SellerName                         string   `json:"sellerName,omitempty"`
	Features                           []string `json:"features,omitempty"`
	SupportedDevices                   []string `json:"supportedDevices,omitempty"`
	Advisories                         []string `json:"advisories,omitempty"`
	ScreenshotUrls                     []string `json:"screenshotUrls,omitempty"`
	IPadScreenshotUrls                 []string `json:"ipadScreenshotUrls,omitempty"`
	AppletvScreenshotUrls              []string `json:"appletvScreenshotUrls,omitempty"`
	ArtistLinkURL                      string   `json:"artistLinkUrl,omitempty"`
	CollectionArtistID                 int64    `json:"collectionArtistId,omitempty"`
	CollectionArtistName               string   `json:"collectionArtistName,omitempty"`
	CollectionArtistViewURL            string   `json:"collectionArtistViewUrl,omitempty"`
	Description                        string   `json:"description,omitempty"`
	Version                            string   `json:"version,omitempty"`
	FileSizeBytes                      string   `json:"fileSizeBytes,omitempty"`
	MinimumOSVersion                   string   `json:"minimumOsVersion,omitempty"`
	AverageUserRating                  float64  `json:"averageUserRating,omitempty"`
	UserRatingCount                    int      `json:"userRatingCount,omitempty"`
	AverageUserRatingForCurrentVersion float64  `json:"averageUserRatingForCurrentVersion,omitempty"`
	UserRatingCountForCurrentVersion   int      `json:"userRatingCountForCurrentVersion,omitempty"`
	FormattedPrice                     string   `json:"formattedPrice,omitempty"`
	Price                              float64  `json:"price,omitempty"`
	BundleID                           string   `json:"bundleId,omitempty"`
	GenreIDs                           []string `json:"genreIds,omitempty"`
	Genres                             []string `json:"genres,omitempty"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A,omitempty"`
}
