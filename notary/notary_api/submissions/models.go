package submissions

// NewSubmissionRequest is the request body for the Submit Software endpoint.
type NewSubmissionRequest struct {
	SHA256         string                             `json:"sha256"`
	SubmissionName string                             `json:"submissionName"`
	Notifications  []NewSubmissionRequestNotification `json:"notifications,omitempty"`
}

// NewSubmissionRequestNotification specifies a webhook callback for when
// notarization completes.
type NewSubmissionRequestNotification struct {
	Channel string `json:"channel"`
	Target  string `json:"target"`
}

// NewSubmissionResponse is returned by the Submit Software endpoint.
type NewSubmissionResponse struct {
	Data NewSubmissionResponseData `json:"data"`
	Meta map[string]any            `json:"meta"`
}

// NewSubmissionResponseData contains the submission ID and S3 upload credentials.
type NewSubmissionResponseData struct {
	Attributes NewSubmissionResponseAttributes `json:"attributes"`
	ID         string                          `json:"id"`
	Type       string                          `json:"type"`
}

// NewSubmissionResponseAttributes holds the temporary AWS credentials and S3
// bucket/key needed to upload the software for notarization.
type NewSubmissionResponseAttributes struct {
	AWSAccessKeyID     string `json:"awsAccessKeyId"`
	AWSSecretAccessKey string `json:"awsSecretAccessKey"`
	AWSSessionToken    string `json:"awsSessionToken"`
	Bucket             string `json:"bucket"`
	Object             string `json:"object"`
}

// SubmissionResponse is returned by the Get Submission Status endpoint.
type SubmissionResponse struct {
	Data SubmissionResponseData `json:"data"`
	Meta map[string]any         `json:"meta"`
}

// SubmissionResponseData describes a single submission and its current status.
type SubmissionResponseData struct {
	Attributes SubmissionAttributes `json:"attributes"`
	ID         string               `json:"id"`
	Type       string               `json:"type"`
}

// SubmissionAttributes contains the status and metadata of a submission.
type SubmissionAttributes struct {
	CreatedDate string `json:"createdDate"`
	Name        string `json:"name"`
	Status      string `json:"status"`
}

// SubmissionLogURLResponse is returned by the Get Submission Log endpoint.
type SubmissionLogURLResponse struct {
	Data SubmissionLogURLResponseData `json:"data"`
	Meta map[string]any               `json:"meta"`
}

// SubmissionLogURLResponseData contains the URL for downloading the notarization log.
type SubmissionLogURLResponseData struct {
	Attributes SubmissionLogURLAttributes `json:"attributes"`
	ID         string                     `json:"id"`
	Type       string                     `json:"type"`
}

// SubmissionLogURLAttributes holds the temporary URL to download the log JSON.
type SubmissionLogURLAttributes struct {
	DeveloperLogURL string `json:"developerLogUrl"`
}

// SubmissionListResponse is returned by the Get Previous Submissions endpoint.
type SubmissionListResponse struct {
	Data []SubmissionListResponseData `json:"data"`
	Meta map[string]any               `json:"meta"`
}

// SubmissionListResponseData describes one entry in the submission history list.
type SubmissionListResponseData struct {
	Attributes SubmissionAttributes `json:"attributes"`
	ID         string               `json:"id"`
	Type       string               `json:"type"`
}
