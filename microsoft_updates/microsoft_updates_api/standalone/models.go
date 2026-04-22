package standalone

import "encoding/xml"

// StandaloneResponse holds all packages fetched across one CDN channel.
type StandaloneResponse struct {
	// Packages contains one entry per application ID fetched from the CDN.
	Packages []*Package
}

// Package represents a single Microsoft application update entry.
// Fields map to the Apple plist dict keys returned by the Microsoft CDN XML.
type Package struct {
	// ApplicationID is the Microsoft CDN application identifier (e.g. "MSWD2019").
	ApplicationID string

	// Title is the human-readable application name from the plist (e.g. "Microsoft Word").
	Title string

	// ShortVersion is the user-facing version string (e.g. "16.108.1").
	ShortVersion string

	// FullVersion is the build version string (e.g. "16.108.26041915").
	FullVersion string

	// MinimumOS is the minimum macOS version required (e.g. "14.0").
	MinimumOS string

	// UpdateVersion is the full update version as returned by the CDN.
	UpdateVersion string

	// Location is the download URL for the full installer package.
	Location string

	// AppOnlyLocation is the download URL for the app-only delta update (may be empty).
	AppOnlyLocation string

	// Hash is the base64-encoded SHA-1 hash of the full installer.
	Hash string

	// HashSHA256 is the base64-encoded SHA-256 hash of the full installer.
	HashSHA256 string

	// AppOnlyHash is the SHA-1 hash of the app-only update (may be empty).
	AppOnlyHash string

	// AppOnlyHashSHA256 is the SHA-256 hash of the app-only update (may be empty).
	AppOnlyHashSHA256 string

	// Date is the release date string as provided by the CDN.
	Date string
}

// plistArray is the top-level plist XML structure returned by the Microsoft CDN.
// The CDN returns an Apple plist with an array of dict entries.
type plistArray struct {
	XMLName xml.Name   `xml:"plist"`
	Items   []plistDict `xml:"array>dict"`
}

// plistDict holds the alternating key/value children of a plist <dict> element.
// Because standard encoding/xml cannot natively handle alternating key/value pairs,
// we capture all child elements and post-process them in pairs.
type plistDict struct {
	Children []plistNode `xml:",any"`
}

// plistNode is a generic plist element — either a <key>, <string>, or <data>.
type plistNode struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// toPackage converts a raw plist dict into a typed Package. It iterates the
// alternating key/value children and maps known keys to Package fields.
func (d *plistDict) toPackage(appID string) *Package {
	p := &Package{ApplicationID: appID}
	children := d.Children
	for i := 0; i+1 < len(children); i += 2 {
		key := children[i].Value
		val := children[i+1].Value
		switch key {
		case "Title":
			p.Title = val
		case "Update Version":
			p.UpdateVersion = val
		case "Short Version":
			p.ShortVersion = val
		case "Minimum OS":
			p.MinimumOS = val
		case "Location":
			p.Location = val
		case "App Only Location":
			p.AppOnlyLocation = val
		case "Hash":
			p.Hash = val
		case "Hash SHA-256":
			p.HashSHA256 = val
		case "App Only Hash":
			p.AppOnlyHash = val
		case "App Only Hash SHA-256":
			p.AppOnlyHashSHA256 = val
		case "Date":
			p.Date = val
		case "Full Version":
			p.FullVersion = val
		}
	}
	// Derive FullVersion from UpdateVersion when not set explicitly.
	if p.FullVersion == "" {
		p.FullVersion = p.UpdateVersion
	}
	return p
}
