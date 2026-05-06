package cache

import (
	"net/http"
	"net/http/httptest"
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

func Test_runInvalidateCmd_NoBefore_HitsBackendWithoutQuery(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/agents/abc-123/cache", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "", r.URL.RawQuery)
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewCacheCmd(f)
	_, err := test.Execute(cmd, "invalidate abc-123 -y", out)
	require.NoError(t, err)
	assert.True(t, hit)
}

func Test_runInvalidateCmd_WithBefore_PassesQueryParam(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/cache", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2026-01-15", r.URL.Query().Get("before"))
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewCacheCmd(f)
	_, err := test.Execute(cmd, "invalidate abc-123 --before 2026-01-15 -y", out)
	require.NoError(t, err)
}

func Test_runInvalidateCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/cache", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewCacheCmd(f)
	result, err := test.Execute(cmd, "invalidate abc-123 --before 2026-01-15 --dry-run", out)
	require.NoError(t, err)

	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/agents/abc-123/cache?before=2026-01-15")
	assert.Contains(t, got, "scope: cached completions created before 2026-01-15")
}

func Test_runInvalidateCmd_DryRunNoBefore_DescribesAllScope(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCacheCmd(f)
	result, err := test.Execute(cmd, "invalidate abc-123 --dry-run", out)
	require.NoError(t, err)

	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/agents/abc-123/cache")
	assert.NotContains(t, got, "?before=")
	assert.Contains(t, got, "all cached completions for this agent")
}

func Test_runInvalidateCmd_RequiresAgentID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCacheCmd(f)
	_, err := test.Execute(cmd, "invalidate", out)
	require.Error(t, err)
}

func Test_runInvalidateCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	// test.NewFactory(false, ...) configures IO with non-TTY stdin/stdout/stderr.
	// CanPrompt() must return false in that case, and the command must
	// refuse without --confirm — same contract as `agents delete`.
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCacheCmd(f)
	_, err := test.Execute(cmd, "invalidate abc-123", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runInvalidateCmd_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing/cache", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewCacheCmd(f)
	_, err := test.Execute(cmd, "invalidate missing -y", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Agent not found")
}
