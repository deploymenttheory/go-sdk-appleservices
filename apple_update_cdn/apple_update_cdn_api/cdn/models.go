package cdn

import (
	"time"
)

// ProgressFunc is called periodically during a download with the current
// number of bytes written and the total expected bytes (0 if unknown).
// Implementations must be safe for concurrent use.
type ProgressFunc func(bytesWritten, totalBytes int64)

// DownloadResult contains the outcome of a completed DownloadFileV1 call.
type DownloadResult struct {
	// URL is the Apple CDN URL that was downloaded.
	URL string
	// DestPath is the local filesystem path where the file was written.
	DestPath string
	// BytesWritten is the total number of bytes streamed to disk.
	BytesWritten int64
	// SHA1 is the hex-encoded SHA-1 checksum of the downloaded bytes.
	SHA1 string
	// SHA256 is the hex-encoded SHA-256 checksum of the downloaded bytes.
	SHA256 string
	// Duration is the wall-clock time elapsed during the download (excluding
	// the initial HEAD request).
	Duration time.Duration
	// Verified is true when at least one checksum was provided by the CDN
	// (via HEAD response headers) and all provided checksums matched the
	// downloaded content. False when no checksums were available or a mismatch
	// was detected.
	Verified bool
}

// CDNURLInfo contains the structured components parsed from an Apple CDN IPSW URL.
//
// Apple CDN IPSW URL format:
//
//	https://updates.cdn-apple.com/{CatalogRelease}/{AssetType}/{AssetID}/{UUID}/{Platform}_{Version}_{Build}_{RestoreType}.ipsw
//
// Example:
//
//	https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw
type CDNURLInfo struct {
	// CatalogRelease is the seasonal release identifier, e.g. "2026WinterFCS" or "2025FallFCS".
	CatalogRelease string
	// AssetType describes the package type, typically "fullrestores".
	AssetType string
	// AssetID is Apple's internal asset identifier, e.g. "122-28781".
	AssetID string
	// UUID is the unique asset delivery identifier, e.g. "DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46".
	UUID string
	// Filename is the full IPSW filename, e.g. "UniversalMac_26.4.1_25E253_Restore.ipsw".
	Filename string
	// Platform is the target platform parsed from the filename, e.g. "UniversalMac".
	Platform string
	// Version is the macOS version parsed from the filename, e.g. "26.4.1".
	Version string
	// Build is the macOS build identifier parsed from the filename, e.g. "25E253".
	Build string
	// RestoreType is the restore type parsed from the filename, e.g. "Restore".
	RestoreType string
}

// CDNFileMetadata contains file metadata retrieved from an Apple CDN IPSW URL
// via a HEAD request. Apple's CDN is S3-backed and returns checksums and file
// size in response headers without requiring the full file to be downloaded.
type CDNFileMetadata struct {
	// URL is the original CDN URL that was queried.
	URL string
	// ContentLength is the file size in bytes from the Content-Length header.
	ContentLength int64
	// SHA1 is the SHA-1 checksum from the x-amz-meta-digest-sh1 header.
	SHA1 string
	// SHA256 is the SHA-256 checksum from the x-amz-meta-digest-sha256 header.
	SHA256 string
	// LastModified is the upload timestamp from the Last-Modified header.
	LastModified time.Time
	// ETag is the entity tag from the ETag header.
	ETag string
	// ContentType is the MIME type from the Content-Type header.
	ContentType string
}
