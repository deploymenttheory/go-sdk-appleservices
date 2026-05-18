package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	submissions map[string]map[string]any
}

func init() {
	mockState.submissions = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`))
}

// loadMockResponse loads JSON response from the mocks folder.
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// SubmissionsMock provides httpmock responders for submission endpoints.
type SubmissionsMock struct{}

// RegisterMocks registers all HTTP mock responders for submissions.
func (m *SubmissionsMock) RegisterMocks() {
	mockState.Lock()
	mockState.submissions = make(map[string]map[string]any)
	mockState.Unlock()

	m.seedTestSubmission()

	// POST /notary/v2/submissions — submit software
	httpmock.RegisterResponder("POST", "https://appstoreconnect.apple.com/notary/v2/submissions", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_submit_software.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to load mock data","labels":[]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to parse mock data","labels":[]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /notary/v2/submissions — list previous submissions
	httpmock.RegisterResponder("GET", "https://appstoreconnect.apple.com/notary/v2/submissions", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_previous_submissions.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to load mock data","labels":[]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to parse mock data","labels":[]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /notary/v2/submissions/{submissionId}/logs — get submission log URL
	httpmock.RegisterResponder("GET", `=~^https://appstoreconnect\.apple\.com/notary/v2/submissions/[^/]+/logs$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		// path: /notary/v2/submissions/{id}/logs — id is parts[4]
		if len(parts) < 5 {
			return httpmock.NewStringResponse(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`), nil
		}
		submissionID := parts[len(parts)-2]

		mockState.Lock()
		_, exists := mockState.submissions[submissionID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_submission_log.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to load mock data","labels":[]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to parse mock data","labels":[]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /notary/v2/submissions/{submissionId} — get submission status
	httpmock.RegisterResponder("GET", `=~^https://appstoreconnect\.apple\.com/notary/v2/submissions/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		submissionID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.submissions[submissionID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_submission_status.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to load mock data","labels":[]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Failed to parse mock data","labels":[]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

// RegisterErrorMocks registers mock responders that return error responses.
func (m *SubmissionsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.submissions = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://appstoreconnect.apple.com/notary/v2/submissions", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"name":"INTERNAL_ERROR","description":"Mock error for testing","labels":[]}`), nil
	})

	httpmock.RegisterResponder("GET", "https://appstoreconnect.apple.com/notary/v2/submissions", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(403, `{"name":"FORBIDDEN","description":"Authentication failure","labels":[]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://appstoreconnect\.apple\.com/notary/v2/submissions/[^/]+/logs$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://appstoreconnect\.apple\.com/notary/v2/submissions/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"name":"NOT_FOUND","description":"The specified identifier can't be found","labels":[]}`), nil
	})
}

// CleanupMockState clears all mock state data.
func (m *SubmissionsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.submissions {
		delete(mockState.submissions, id)
	}
}

func (m *SubmissionsMock) seedTestSubmission() {
	testSubmission := map[string]any{
		"type": "submissions",
		"id":   "2efe2717-52ef-43a5-96dc-0797e4ca1041",
	}
	mockState.Lock()
	mockState.submissions["2efe2717-52ef-43a5-96dc-0797e4ca1041"] = testSubmission
	mockState.Unlock()
}
