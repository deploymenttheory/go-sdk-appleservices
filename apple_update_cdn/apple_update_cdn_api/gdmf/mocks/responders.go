package mocks

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// loadFixture reads a JSON fixture file from the mocks directory.
func loadFixture(filename string) string {
	data, err := os.ReadFile(filepath.Join("mocks", filename))
	if err != nil {
		return `{}`
	}
	return string(data)
}

// jsonResponder returns an httpmock.Responder that serves JSON with the given status code.
func jsonResponder(status int, body string) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, body)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	}
}

// RegisterGetPublicVersions registers the mock responder for the GDMF v2 pmv endpoint.
func RegisterGetPublicVersions() {
	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		jsonResponder(200, loadFixture("validate_get_public_versions.json")))
}
