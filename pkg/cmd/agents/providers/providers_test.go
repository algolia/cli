package providers

// Shared test harness for every *_test.go in this package.
// Per-verb tests live in <verb>_test.go to mirror the source split.

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/agentstudio"
)

// newClientForServer builds an AgentStudio client factory that talks to
// the given httptest server. Used by every verb test that needs a
// stubbed backend.
func newClientForServer(t *testing.T, ts *httptest.Server) func() (*agentstudio.Client, error) {
	t.Helper()
	return func() (*agentstudio.Client, error) {
		return agentstudio.NewClient(agentstudio.Config{
			BaseURL:       ts.URL,
			ApplicationID: "APP123",
			APIKey:        "k",
			HTTPClient:    ts.Client(),
		})
	}
}

// writeTempJSON dumps content into a t.TempDir-scoped file and returns
// the path. Used by create/update tests to build -F file inputs.
func writeTempJSON(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}
