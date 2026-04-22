package edge

// EdgeAllChannelsResponse aggregates the latest macOS release for each Edge channel.
type EdgeAllChannelsResponse struct {
	Stable *EdgeRelease
	Beta   *EdgeRelease
	Dev    *EdgeRelease
	Canary *EdgeRelease
}

// EdgeChannelResponse holds the latest macOS release for a single Edge channel.
type EdgeChannelResponse struct {
	Channel string
	Release *EdgeRelease
}

// EdgeRelease represents the latest Edge release for a given channel on macOS.
type EdgeRelease struct {
	// Channel is the distribution channel (stable, beta, dev, canary).
	Channel string

	// Version is the full Edge version string (e.g. "147.0.3912.72").
	Version string

	// PublishedTime is the ISO 8601 timestamp when this version was published.
	PublishedTime string

	// Artifacts contains the downloadable package entries for macOS.
	Artifacts []EdgeArtifact
}

// EdgeArtifact describes a single downloadable artifact within an Edge release.
type EdgeArtifact struct {
	// ArtifactName identifies the package type (e.g. "pkg", "dmg").
	ArtifactName string

	// Location is the download URL for the artifact.
	Location string

	// SizeInBytes is the artifact file size.
	SizeInBytes int64

	// HashSHA256 is the SHA-256 checksum of the artifact.
	HashSHA256 string
}

// edgeAPIProduct is the raw JSON structure returned by edgeupdates.microsoft.com.
type edgeAPIProduct struct {
	Product  string          `json:"Product"`
	Releases []edgeAPIRelease `json:"Releases"`
}

// edgeAPIRelease is a single release entry in the Edge update API response.
type edgeAPIRelease struct {
	ReleaseVersion string              `json:"ReleaseVersion"`
	PublishedTime  string              `json:"PublishedTime"`
	Platform       string              `json:"Platform"`
	Architecture   string              `json:"Architecture"`
	Artifacts      []edgeAPIArtifact   `json:"Artifacts"`
}

// edgeAPIArtifact is a single artifact entry in an Edge release.
type edgeAPIArtifact struct {
	ArtifactName string `json:"ArtifactName"`
	Location     string `json:"Location"`
	SizeInBytes  int64  `json:"SizeInBytes"`
	Hash         string `json:"Hash"`
	HashAlgorithm string `json:"HashAlgorithm"`
}
