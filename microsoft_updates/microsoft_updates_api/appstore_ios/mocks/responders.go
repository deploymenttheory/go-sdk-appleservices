package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed appstore_ios_mock.json
var iosAppStoreJSON []byte

// RegisterMicrosoftWordMock registers an httpmock responder for the iTunes iOS
// search API returning a Microsoft Word result. Matches any GET to the search endpoint.
func RegisterMicrosoftWordMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://itunes.apple.com/search",
		httpmock.NewBytesResponder(200, iosAppStoreJSON),
	)
}
