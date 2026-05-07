package run

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runRunCmd_StreamingHappyPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("stream"))

		body, _ := io.ReadAll(r.Body)
		var req map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &req))
		assert.JSONEq(t, `[{"role":"user","content":"hello"}]`, string(req["messages"]))
		// `agents run` must NOT send `configuration`; the persisted agent
		// has its own. Sending one would 422 the request on the backend.
		_, hasCfg := req["configuration"]
		assert.False(t, hasCfg, "configuration must not be present for agents run")

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(`data: {"type":"text-delta","delta":"hi"}` + "\n\n"))
		_, _ = w.Write([]byte(`data: [DONE]` + "\n\n"))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewRunCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123 -m hello", out)
	require.NoError(t, err)

	got := strings.TrimSpace(result.String())
	var event map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &event))
	assert.Equal(t, "text-delta", event["type"])
}

func Test_runRunCmd_DryRunIncludesAgentID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/completions", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewRunCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123 -m hi --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would POST /1/agents/abc-123/completions")
	assert.Contains(t, got, `"content": "hi"`)
}

func Test_runRunCmd_RequiresAgentID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewRunCmd(f, nil)
	_, err := test.Execute(cmd, "-m hi", out)
	require.Error(t, err)
}

func Test_runRunCmd_ForwardsCompletionFlagsToWire(t *testing.T) {
	// One end-to-end check that all four Phase 5 flags map onto the
	// expected query params + header. Exhaustive matrix lives in
	// api/agentstudio/completions_test.go; this exists so a regression
	// in the cobra→opts→client wiring (forgetting one field,
	// transposing No* polarity, etc.) is caught at the cmd layer.
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "false", r.URL.Query().Get("cache"))
		assert.Equal(t, "false", r.URL.Query().Get("memory"))
		assert.Equal(t, "false", r.URL.Query().Get("analytics"))
		assert.Equal(t, "ey.signed.jwt", r.Header.Get("X-Algolia-Secure-User-Token"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"role":"assistant","content":"ok"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewRunCmd(f, nil)
	_, err := test.Execute(cmd,
		"abc-123 -m hi --no-stream --no-cache --no-memory --no-analytics --secure-user-token ey.signed.jwt",
		out)
	require.NoError(t, err)
}

func Test_runRunCmd_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing/completions", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewRunCmd(f, nil)
	_, err := test.Execute(cmd, "missing -m hi", out)
	require.Error(t, err)
}
