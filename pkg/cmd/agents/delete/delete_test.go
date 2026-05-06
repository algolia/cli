package delete

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/test"
)

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

func agentJSON(name, status string) string {
	return `{
		"id":"abc-123","name":"` + name + `","status":"` + status + `",
		"instructions":"x","createdAt":"2025-01-01T00:00:00Z"
	}`
}

func Test_runDeleteCmd_NonTTYRequiresConfirm(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "abc-123", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runDeleteCmd_NonTTYWithConfirmDeletes(t *testing.T) {
	var deleted atomic.Bool

	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(agentJSON("Concierge", "draft")))
		case http.MethodDelete:
			deleted.Store(true)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "abc-123 -y", out)
	require.NoError(t, err)
	assert.True(t, deleted.Load(), "DELETE should have been called")
}

func Test_runDeleteCmd_DryRunDoesNotDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(agentJSON("Concierge", "published")))
		case http.MethodDelete:
			t.Fatal("DELETE called during --dry-run")
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewDeleteCmd(f, nil)
	// --dry-run alone is enough; no --confirm needed because it's non-destructive.
	result, err := test.Execute(cmd, "abc-123 --dry-run", out)
	require.NoError(t, err)

	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/agents/abc-123")
	assert.Contains(t, got, "name:   Concierge")
	assert.Contains(t, got, "status: published")
}

func Test_runDeleteCmd_PropagatesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "missing -y", out)
	require.Error(t, err)
}
