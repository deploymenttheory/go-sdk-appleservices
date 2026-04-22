package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed cve_history_mock.html
var cveHistoryHTML []byte

// RegisterCVEHistoryMock registers an httpmock responder for the CVE history page.
func RegisterCVEHistoryMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://learn.microsoft.com/en-us/officeupdates/release-notes-office-for-mac",
		httpmock.NewBytesResponder(200, cveHistoryHTML),
	)
}
