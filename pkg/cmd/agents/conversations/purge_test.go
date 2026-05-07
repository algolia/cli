package conversations

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runPurgeCmd_RefusesDateless(t *testing.T) {
	// The most important test in this package: backend rejects dateless
	// DELETE with 400 "At least one filter is required" — the CLI surfaces
	// this as a flag-level error to skip the round-trip.
	// (Anya's Phase 7 vet caught this; spec marks both date params as
	// optional but the live backend disagrees.)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "purge agent-1 -y", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one of --start-date / --end-date is required")
}

func Test_runPurgeCmd_AcceptsOpenEndedStart(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "1970-01-01", r.URL.Query().Get("startDate"))
		assert.Empty(t, r.URL.Query().Get("endDate"))
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "purge agent-1 --start-date 1970-01-01 -y", out)
	require.NoError(t, err)
}

func Test_runPurgeCmd_DryRunWithDateRange(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	result, err := test.Execute(cmd,
		"purge agent-1 --start-date 2026-01-01 --end-date 2026-01-31 --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/agents/agent-1/conversations?")
	assert.Contains(t, got, "startDate=2026-01-01")
	assert.Contains(t, got, "endDate=2026-01-31")
	assert.Contains(t, got, "scope: between 2026-01-01 and 2026-01-31")
}

func Test_runPurgeCmd_HitsBackendWithDateRange(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "2026-01-01", r.URL.Query().Get("startDate"))
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "purge agent-1 --start-date 2026-01-01 -y", out)
	require.NoError(t, err)
}

func Test_runPurgeCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "purge agent-1 --start-date 2026-01-01", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}
