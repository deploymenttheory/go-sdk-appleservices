package firmware

import "time"

// AllFirmwaresV3Response is the top-level response from the ipsw.me v3 condensed endpoint.
// The Devices map is keyed by Apple model identifier (e.g. "Mac14,3").
// Non-Mac entries (iTunes, etc.) are present in the raw API but filtered by
// ListAllMacFirmwareV3.
type AllFirmwaresV3Response struct {
	Devices map[string]*MacDevice `json:"devices"`
}

// MacDevice represents a single Mac model and its full firmware history.
type MacDevice struct {
	// Name is the human-readable product name, e.g. "Mac mini (M2, 2023)".
	Name string `json:"name"`
	// BoardConfig is the internal board configuration identifier, e.g. "J473AP".
	BoardConfig string `json:"BoardConfig"`
	// Platform is the internal platform string, e.g. "REALBRIDGE".
	Platform string `json:"platform"`
	// CPID is the chip identifier.
	CPID int `json:"cpid"`
	// BDID is the board identifier.
	BDID int `json:"bdid"`
	// Firmwares is the list of all known firmware versions for this device.
	Firmwares []*FirmwareV3 `json:"firmwares"`
}

// FirmwareV3 represents a single macOS IPSW firmware entry as returned by the
// ipsw.me v3 API. All macOS IPSW files are Universal — the same URL is shared
// across all Apple Silicon Mac models for a given version.
type FirmwareV3 struct {
	// Version is the human-readable macOS version, e.g. "26.4.1".
	Version string `json:"version"`
	// BuildID is the macOS build identifier, e.g. "25E253".
	BuildID string `json:"buildid"`
	// SHA1Sum is the SHA-1 checksum of the IPSW file.
	SHA1Sum string `json:"sha1sum"`
	// MD5Sum is the MD5 checksum of the IPSW file.
	MD5Sum string `json:"md5sum"`
	// Size is the file size in bytes.
	Size int64 `json:"size"`
	// ReleaseDate is the date the firmware was publicly released.
	ReleaseDate time.Time `json:"releasedate"`
	// UploadDate is the date the firmware was uploaded to the CDN.
	UploadDate time.Time `json:"uploaddate"`
	// URL is the full Apple CDN download URL for the IPSW file.
	URL string `json:"url"`
	// Signed indicates whether Apple's signing servers currently accept this firmware.
	Signed bool `json:"signed"`
	// Filename is the base filename of the IPSW, e.g. "UniversalMac_26.4.1_25E253_Restore.ipsw".
	Filename string `json:"filename"`
}

// DeviceFirmwaresV4Response is the response from the ipsw.me v4 device endpoint.
// It is richer than the v3 response, including SHA-256 checksums and filesize.
type DeviceFirmwaresV4Response struct {
	// Name is the human-readable product name.
	Name string `json:"name"`
	// Identifier is the Apple model identifier, e.g. "Mac14,3".
	Identifier string `json:"identifier"`
	// Firmwares is the list of all known firmware versions for this device.
	Firmwares []*FirmwareV4 `json:"firmwares"`
	// Boards lists the board configuration identifiers for this device.
	Boards []string `json:"boards"`
	// BoardConfig is the primary board configuration identifier.
	BoardConfig string `json:"boardconfig"`
	// Platform is the internal platform string.
	Platform string `json:"platform"`
	// CPID is the chip identifier.
	CPID int `json:"cpid"`
	// BDID is the board identifier.
	BDID int `json:"bdid"`
}

// FirmwareV4 represents a single macOS IPSW firmware entry as returned by the
// ipsw.me v4 device endpoint. Includes SHA-256 in addition to SHA-1/MD5.
type FirmwareV4 struct {
	// Identifier is the Apple model identifier this firmware entry is associated with.
	Identifier string `json:"identifier"`
	// Version is the human-readable macOS version, e.g. "26.4.1".
	Version string `json:"version"`
	// BuildID is the macOS build identifier, e.g. "25E253".
	BuildID string `json:"buildid"`
	// SHA1Sum is the SHA-1 checksum of the IPSW file.
	SHA1Sum string `json:"sha1sum"`
	// MD5Sum is the MD5 checksum of the IPSW file.
	MD5Sum string `json:"md5sum"`
	// SHA256Sum is the SHA-256 checksum of the IPSW file.
	SHA256Sum string `json:"sha256sum"`
	// FileSize is the file size in bytes.
	FileSize int64 `json:"filesize"`
	// URL is the full Apple CDN download URL for the IPSW file.
	URL string `json:"url"`
	// ReleaseDate is the date the firmware was publicly released.
	ReleaseDate time.Time `json:"releasedate"`
	// UploadDate is the date the firmware was uploaded to the CDN.
	UploadDate time.Time `json:"uploaddate"`
	// Signed indicates whether Apple's signing servers currently accept this firmware.
	Signed bool `json:"signed"`
}
