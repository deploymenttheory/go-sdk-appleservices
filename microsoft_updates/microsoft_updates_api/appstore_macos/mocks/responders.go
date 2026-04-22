package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed appstore_macos_mock.json
var macOSAppStoreJSON []byte

// RegisterMicrosoftWordMock registers an httpmock responder for the iTunes search
// API returning a Microsoft Word macOS result. Matches any GET to the search endpoint.
func RegisterMicrosoftWordMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://itunes.apple.com/search",
		httpmock.NewBytesResponder(200, macOSAppStoreJSON),
	)
}

// RegisterErrorMock registers a 500 error responder for the given URL.
func RegisterErrorMock(u string) {
	httpmock.RegisterResponder(
		"GET",
		u,
		httpmock.NewStringResponder(500, `Internal Server Error`),
	)
}
