// Package sharedtest holds test helpers reused by `algolia agents`
// command tests. Kept in a dedicated package so production code can't
// accidentally import it.
package sharedtest

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/algolia/cli/api/agentstudio"
)

// NewClient returns an AgentStudioClient factory pointing at ts. Used
// by every command-level _test.go that wires an httptest server into
// the cmdutil.Factory.
func NewClient(t *testing.T, ts *httptest.Server) func() (*agentstudio.Client, error) {
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

// WriteTempJSON writes content to a fresh tempfile and returns its
// absolute path. Cleaned up automatically when the test ends.
func WriteTempJSON(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp %s: %v", name, err)
	}
	return path
}
