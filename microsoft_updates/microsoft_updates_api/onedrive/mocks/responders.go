package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed onedrive_manifest_mock.xml
var manifestXML []byte

// RegisterProductionManifestMock registers an httpmock responder for the OneDrive
// standalone production manifest endpoint.
func RegisterProductionManifestMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://g.live.com/0USSDMC_W5T/StandaloneProductManifest",
		httpmock.NewBytesResponder(200, manifestXML),
	)
}

// RegisterInsiderManifestMock registers an httpmock responder for the OneDrive
// insider feed endpoint.
func RegisterInsiderManifestMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://g.live.com/0USSDMC_W5T/MacODSUInsiders",
		httpmock.NewBytesResponder(200, manifestXML),
	)
}

// RegisterErrorMock registers a 500 error responder for the given URL.
func RegisterErrorMock(url string) {
	httpmock.RegisterResponder(
		"GET",
		url,
		httpmock.NewStringResponder(500, `Internal Server Error`),
	)
}
