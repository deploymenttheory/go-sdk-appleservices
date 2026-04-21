package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

const (
	// TestIPSWURL is the canonical test CDN URL used in all CDN mock tests.
	TestIPSWURL = "https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw"

	// TestSHA1 is the expected SHA-1 checksum returned by the mock CDN HEAD response.
	TestSHA1 = "03078f4af82bff5473398ca49f99288c76253fe8"

	// TestSHA256 is the expected SHA-256 checksum returned by the mock CDN HEAD response.
	TestSHA256 = "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90"

	// TestContentLength is the expected file size in bytes.
	TestContentLength = int64(19734779897)
)

// RegisterGetFileMetadata registers a HEAD responder for TestIPSWURL that returns
// realistic Apple CDN response headers (checksums, size, last-modified).
func RegisterGetFileMetadata() {
	httpmock.RegisterResponder("HEAD", TestIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Type", "application/x-apple-aspen-config")
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", TestContentLength))
			resp.Header.Set("x-amz-meta-digest-sh1", TestSHA1)
			resp.Header.Set("x-amz-meta-digest-sha256", TestSHA256)
			resp.Header.Set("Last-Modified", "Mon, 27 Jan 2025 12:00:00 GMT")
			resp.Header.Set("Etag", `"abc123etag"`)
			return resp, nil
		})
}
