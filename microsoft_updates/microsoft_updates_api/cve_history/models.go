package cve_history

// CVEHistoryResponse holds the full Office for Mac CVE/security release history.
type CVEHistoryResponse struct {
	Entries []CVEEntry
}

// CVEEntry represents a single security-relevant Office for Mac release with
// associated CVE identifiers.
type CVEEntry struct {
	// ReleaseDate is the release date string as scraped from the page heading.
	ReleaseDate string

	// Version is the Office version string for this release (e.g. "16.108").
	Version string

	// CVEs is the list of CVE identifiers addressed in this release
	// (e.g. ["CVE-2026-12345", "CVE-2026-67890"]).
	CVEs []string
}
