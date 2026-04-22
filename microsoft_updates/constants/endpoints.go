package constants

// Microsoft Office CDN channel UUIDs and base URLs.
// The production, beta, and preview channels each have a unique UUID that forms
// part of the CDN path. Per-application XML feeds are fetched from:
//
//	{ChannelBaseURL}{ApplicationID}.xml
const (
	// StandaloneCDNBaseURL is the production (stable) Office CDN channel.
	StandaloneCDNBaseURL = "https://officecdnmac.microsoft.com/pr/C1297A47-86C4-4C1F-97FA-950631F94777/MacAutoupdate/"

	// StandaloneBetaCDNBaseURL is the beta (Insider Fast) Office CDN channel.
	StandaloneBetaCDNBaseURL = "https://officecdnmac.microsoft.com/pr/1ac37578-5a24-40fb-892e-b89d85b6dfaa/MacAutoupdate/"

	// StandalonePreviewCDNBaseURL is the preview (Insider Slow) Office CDN channel.
	StandalonePreviewCDNBaseURL = "https://officecdnmac.microsoft.com/pr/4B2D7701-0A4F-49C8-B4CB-0C2D4043F51F/MacAutoupdate/"
)

// Microsoft Edge update API endpoints.
// Each channel returns a JSON array of release objects filtered by platform.
const (
	EdgeUpdateAPIBase    = "https://edgeupdates.microsoft.com/api/products/"
	EdgeStableEndpoint   = EdgeUpdateAPIBase + "stable"
	EdgeBetaEndpoint     = EdgeUpdateAPIBase + "beta"
	EdgeDevEndpoint      = EdgeUpdateAPIBase + "dev"
	EdgeCanaryEndpoint   = EdgeUpdateAPIBase + "canary"
)

// OneDrive distribution ring endpoints.
// fwlink URLs redirect to the latest installer for each ring.
// g.live.com endpoints serve XML manifests with version and download metadata.
const (
	OneDriveLiveFeedInsiders   = "https://g.live.com/0USSDMC_W5T/MacODSUInsiders"
	OneDriveStandaloneManifest = "https://g.live.com/0USSDMC_W5T/StandaloneProductManifest"
	OneDriveFwlinkDeferred     = "https://go.microsoft.com/fwlink/?linkid=861009"
	OneDriveFwlinkUpcoming     = "https://go.microsoft.com/fwlink/?linkid=861010"
	OneDriveFwlinkRollingOut   = "https://go.microsoft.com/fwlink/?linkid=861011"
	OneDriveFwlinkAppNew       = "https://go.microsoft.com/fwlink/?linkid=823060"
)

// iTunes / App Store search API.
// Append entity and term query parameters to target macOS or iOS apps.
const (
	iTunesSearchAPI = "https://itunes.apple.com/search"

	// AppStoreSearchEndpoint is the full iTunes search base URL.
	AppStoreSearchEndpoint = iTunesSearchAPI
)

// Microsoft Learn HTML pages used for update history and CVE scraping.
const (
	UpdateHistoryURL = "https://learn.microsoft.com/en-us/officeupdates/update-history-office-for-mac"
	CVEHistoryURL    = "https://learn.microsoft.com/en-us/officeupdates/release-notes-office-for-mac"
)
