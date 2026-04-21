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

// RegisterListAllMacFirmware registers the mock responder for the v3 condensed endpoint.
func RegisterListAllMacFirmware() {
	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, loadFixture("validate_list_all_mac_firmware.json")))
}

// RegisterGetDeviceFirmware registers the mock responder for the v4 device endpoint.
// identifier should match the path segment, e.g. "Mac14,3".
func RegisterGetDeviceFirmware(identifier string) {
	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v4/device/"+identifier,
		jsonResponder(200, loadFixture("validate_get_device_firmware.json")))
}
