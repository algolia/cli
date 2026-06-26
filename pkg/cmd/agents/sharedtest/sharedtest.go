// Package sharedtest holds test helpers reused by `algolia agents`
// command tests. Kept in a dedicated package so production code can't
// accidentally import it.
package sharedtest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/call"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"

	"github.com/algolia/cli/api/agentstudio"
)

// httptestRequester adapts an httptest server's *http.Client to the SDK's
// transport.Requester interface so the SDK client talks to the test server.
type httptestRequester struct{ client *http.Client }

func (r httptestRequester) Request(req *http.Request, _, _ time.Duration) (*http.Response, error) {
	return r.client.Do(req)
}

// NewClient returns a local AgentStudioClient factory pointing at ts. Used
// by command tests that exercise the hand-rolled client (run/try streaming,
// internal endpoints, duplicate).
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

// NewAPIClient returns an official SDK AgentStudioAPIClient factory pointing at
// ts. The SDK prepends the /agent-studio/1/... path itself, so test handlers
// must register routes under that prefix. The host is taken from ts.URL with an
// http scheme so httptest.NewServer (plain HTTP) is reached without TLS.
func NewAPIClient(t *testing.T, ts *httptest.Server) func() (*agentStudio.APIClient, error) {
	t.Helper()
	return func() (*agentStudio.APIClient, error) {
		u, err := url.Parse(ts.URL)
		if err != nil {
			return nil, err
		}
		return agentStudio.NewClientWithConfig(agentStudio.AgentStudioConfiguration{
			Configuration: transport.Configuration{
				AppID:         "APP123",
				ApiKey:        "k",
				DefaultHeader: make(map[string]string),
				Requester:     httptestRequester{client: ts.Client()},
				Hosts: []transport.StatefulHost{
					transport.NewStatefulHost("http", u.Host, call.IsReadWrite),
				},
			},
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
