package onedrive

import "encoding/xml"

// OneDriveAllRingsResponse aggregates metadata for all OneDrive distribution rings.
type OneDriveAllRingsResponse struct {
	Rings []*OneDriveRing
}

// OneDriveRing holds version and download metadata for a single OneDrive distribution ring.
type OneDriveRing struct {
	// Ring is the distribution ring name (e.g. "Production", "Deferred").
	Ring string

	// Version is the OneDrive version string for this ring (e.g. "26.062.0402").
	Version string

	// BuildVersion is the full build version if available.
	BuildVersion string

	// DownloadURL is the installer download URL for this ring.
	DownloadURL string

	// ApplicationID is the Microsoft application identifier (always "ONDR18").
	ApplicationID string

	// BundleID is the macOS bundle identifier (always "com.microsoft.OneDrive").
	BundleID string
}

// oneDriveManifest is the root element of the OneDrive XML manifest from g.live.com.
type oneDriveManifest struct {
	XMLName xml.Name          `xml:"MicrosoftUpdateCatalog"`
	Items   []oneDrivePackage `xml:"UpdateInfo"`
}

// oneDrivePackage represents a single update entry in the OneDrive manifest.
type oneDrivePackage struct {
	Version      string `xml:"Version"`
	BuildVersion string `xml:"BuildVersion"`
	PackageURL   string `xml:"PackageURL"`
	AppID        string `xml:"AppID"`
}
