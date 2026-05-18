package submissions

// Submission status values returned by the Notary API
const (
	SubmissionStatusAccepted   = "Accepted"
	SubmissionStatusInProgress = "In Progress"
	SubmissionStatusInvalid    = "Invalid"
	SubmissionStatusRejected   = "Rejected"
)

// NotificationChannelWebhook is the only supported notification channel
const NotificationChannelWebhook = "webhook"
