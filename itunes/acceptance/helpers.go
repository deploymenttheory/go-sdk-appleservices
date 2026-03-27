package acceptance

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// RequireClient ensures the shared iTunes client is initialised. Unlike the
// AXM client, the iTunes client requires no credentials — if initialisation
// fails it is a hard error rather than a skip.
func RequireClient(t *testing.T) {
	t.Helper()
	if Client == nil {
		err := InitClient()
		require.NoError(t, err, "failed to initialise iTunes client")
	}
}

// NewContext returns a context bound to the configured request timeout.
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), Config.RequestTimeout)
}

// LogTestStage logs a named stage with optional GitHub Actions ::notice annotation.
func LogTestStage(t *testing.T, stage, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::notice title=%s::%s\n", stage, formatted)
	}
	if Config.Verbose {
		t.Logf("[%s] %s", stage, formatted)
	}
}

// LogTestSuccess logs a successful step.
func LogTestSuccess(t *testing.T, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::notice title=Success::%s\n", formatted)
	}
	if Config.Verbose {
		t.Logf("OK: %s", formatted)
	}
}

// LogTestWarning logs a non-fatal warning.
func LogTestWarning(t *testing.T, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::warning title=Warning::%s\n", formatted)
	}
	if Config.Verbose {
		t.Logf("WARNING: %s", formatted)
	}
}

func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}
