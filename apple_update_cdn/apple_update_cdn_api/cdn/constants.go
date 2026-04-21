package cdn

const (
	// CDNHost is the hostname of Apple's firmware CDN.
	CDNHost = "updates.cdn-apple.com"

	// HeaderSHA1 is the response header name containing the SHA-1 checksum.
	// Note: Apple uses "sh1" (not "sha1") in this header name.
	HeaderSHA1 = "x-amz-meta-digest-sh1"

	// HeaderSHA256 is the response header name containing the SHA-256 checksum.
	HeaderSHA256 = "x-amz-meta-digest-sha256"

	// IPSWExtension is the file extension for Apple firmware files.
	IPSWExtension = ".ipsw"

	// ExpectedPathSegments is the number of path segments expected in a valid CDN URL
	// after stripping the leading slash: catalogRelease/assetType/assetID/uuid/filename
	ExpectedPathSegments = 5

	// ExpectedFilenameSegments is the number of underscore-delimited segments expected
	// in the IPSW filename after stripping the .ipsw extension: platform_version_build_restoreType
	ExpectedFilenameSegments = 4
)
