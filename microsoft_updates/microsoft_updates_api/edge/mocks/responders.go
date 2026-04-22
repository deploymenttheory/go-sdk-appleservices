package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed edge_stable_mock.json
var stableJSON []byte

// RegisterStableMock registers an httpmock responder for the Edge stable endpoint.
func RegisterStableMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://edgeupdates.microsoft.com/api/products/stable",
		httpmock.NewBytesResponder(200, stableJSON),
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
