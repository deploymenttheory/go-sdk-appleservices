package mocks

import (
	_ "embed"

	"github.com/jarcoal/httpmock"
)

//go:embed update_history_mock.html
var updateHistoryHTML []byte

// RegisterUpdateHistoryMock registers an httpmock responder for the update history page.
func RegisterUpdateHistoryMock() {
	httpmock.RegisterResponder(
		"GET",
		"https://learn.microsoft.com/en-us/officeupdates/update-history-office-for-mac",
		httpmock.NewBytesResponder(200, updateHistoryHTML),
	)
}
