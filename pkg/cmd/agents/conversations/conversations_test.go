package conversations

// Shared test harness for every *_test.go in this package.
// Per-verb tests live in <verb>_test.go to mirror the source split.

import (
	"net/http/httptest"
	"testing"

	"github.com/algolia/cli/api/agentstudio"
)

// newClientForServer builds an AgentStudio client factory that talks to
// the given httptest server.
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
