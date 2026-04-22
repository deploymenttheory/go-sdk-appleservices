package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed standalone_word_mock.xml
var wordPlistXML []byte

// RegisterWordMock registers an httpmock responder for the Microsoft Word CDN endpoint.
func RegisterWordMock(baseURL string) {
	httpmock.RegisterResponder(
		"GET",
		baseURL+"MSWD2019.xml",
		httpmock.NewBytesResponder(200, wordPlistXML),
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
