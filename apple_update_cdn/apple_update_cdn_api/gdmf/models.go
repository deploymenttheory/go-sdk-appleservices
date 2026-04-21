package gdmf

// GDMFResponse is the top-level response from the Apple GDMF API
// (https://gdmf.apple.com/v2/pmv).
//
// PublicAssetSets contains firmware versions Apple has publicly released.
// AssetSets may contain additional entries visible to devices enrolled in
// seed programmes.
// PublicBackgroundSecurityImprovements lists rapid security responses.
type GDMFResponse struct {
	PublicAssetSets                      *PlatformAssetSets `json:"PublicAssetSets"`
	AssetSets                            *PlatformAssetSets `json:"AssetSets"`
	PublicBackgroundSecurityImprovements *PlatformAssetSets `json:"PublicBackgroundSecurityImprovements"`
}

// PlatformAssetSets groups asset entries by operating system platform.
type PlatformAssetSets struct {
	// IOS contains firmware entries for iPhone, iPad, and Apple Watch.
	IOS []*AssetEntry `json:"iOS"`
	// MacOS contains firmware entries for Mac computers.
	MacOS []*AssetEntry `json:"macOS"`
	// VisionOS contains firmware entries for Apple Vision Pro.
	VisionOS []*AssetEntry `json:"visionOS"`
}

// AssetEntry represents a single firmware version entry in the GDMF feed.
// This is Apple's authoritative source for which versions are currently signed
// and which devices they support.
type AssetEntry struct {
	// ProductVersion is the human-readable OS version, e.g. "26.4.1" or "15.7.5".
	ProductVersion string `json:"ProductVersion"`
	// Build is the OS build identifier, e.g. "25E253".
	Build string `json:"Build"`
	// PostingDate is the date this version was made publicly available (YYYY-MM-DD).
	PostingDate string `json:"PostingDate"`
	// ExpirationDate is the date after which this version may no longer be signed (YYYY-MM-DD).
	ExpirationDate string `json:"ExpirationDate"`
	// SupportedDevices lists board configuration identifiers (e.g. "J473AP") and
	// Mac model identifiers (e.g. "Mac-1E7E29AD0135F9BC") that support this version.
	SupportedDevices []string `json:"SupportedDevices"`
}
