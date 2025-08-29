package itunes_search

type SearchResponse struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

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
	ShortDescription                   string   `json:"shortDescription,omitempty"`
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
	GenreIDS                           []string `json:"genreIds,omitempty"`
	Genres                             []string `json:"genres,omitempty"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A,omitempty"`
}
