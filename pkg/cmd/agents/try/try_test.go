package try

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func writeTempJSON(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func Test_runTryCmd_StreamingHappyPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "true", r.URL.Query().Get("stream"))
		assert.Equal(t, "ai-sdk-5", r.URL.Query().Get("compatibilityMode"))

		body, _ := io.ReadAll(r.Body)
		var req map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &req))
		assert.JSONEq(t, `[{"role":"user","content":"hello"}]`, string(req["messages"]))
		assert.JSONEq(t, `{"model":"gpt-4o-mini"}`, string(req["configuration"]))

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"type\":\"text-delta\",\"delta\":\"hi\"}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"gpt-4o-mini"}`)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewTryCmd(f, nil)
	result, err := test.Execute(cmd, "-c "+cfgPath+" -m hello", out)
	require.NoError(t, err)

	got := strings.TrimSpace(result.String())
	require.NotEmpty(t, got)
	var event map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &event))
	assert.Equal(t, "text-delta", event["type"])
}

func Test_runTryCmd_NoStreamReturnsBufferedJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "false", r.URL.Query().Get("stream"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"role":"assistant","content":"hi"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewTryCmd(f, nil)
	result, err := test.Execute(cmd, "-c "+cfgPath+" -m hi --no-stream", out)
	require.NoError(t, err)
	assert.Equal(t, `{"role":"assistant","content":"hi"}`, result.String())
}

// No DryRun test on purpose. The whole command IS the dry-run in the
// conversational-ai sense (unsaved-config completion); a CLI-level
// --dry-run would be "dry-run a dry-run" — exactly the conceptual
// collision that motivated the rename from `agents test`. See
// AGENTS.md → "On `--dry-run`".

func Test_runTryCmd_RejectsNeitherInputNorMessage(t *testing.T) {
	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewTryCmd(f, nil)
	_, err := test.Execute(cmd, "-c "+cfgPath, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "one of --input or --message")
}

func Test_runTryCmd_RejectsBothInputAndMessage(t *testing.T) {
	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)
	msgPath := writeTempJSON(t, "msgs.json", `[{"role":"user","content":"x"}]`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewTryCmd(f, nil)
	_, err := test.Execute(cmd, "-c "+cfgPath+" -i "+msgPath+" -m hi", out)
	require.Error(t, err)
	// cobra rejects mutually exclusive flags before our handler runs;
	// it phrases this as "none of the others can be" — match on a
	// substring that's robust across cobra versions.
	assert.Contains(t, err.Error(), "[input message]")
}

func Test_runTryCmd_RequiresConfigFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewTryCmd(f, nil)
	_, err := test.Execute(cmd, "-m hi", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `required flag(s) "config"`)
}

func Test_runTryCmd_CompatibilityV4(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ai-sdk-4", r.URL.Query().Get("compatibilityMode"))
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(`0:"hello"` + "\n"))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewTryCmd(f, nil)
	result, err := test.Execute(cmd, "-c "+cfgPath+" -m hi --compatibility v4", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"text"`)
}

func Test_runTryCmd_ForwardsCompletionFlagsToWire(t *testing.T) {
	// One end-to-end check that all four Phase 5 flags map onto the
	// expected query params + header. Exhaustive matrix lives in
	// api/agentstudio/completions_test.go; this test exists so a
	// regression in the cobra→opts→client wiring (forgetting one
	// field, transposing No* polarity, etc.) is caught at the cmd
	// layer.
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "false", r.URL.Query().Get("cache"))
		assert.Equal(t, "false", r.URL.Query().Get("memory"))
		assert.Equal(t, "false", r.URL.Query().Get("analytics"))
		assert.Equal(t, "ey.signed.jwt", r.Header.Get("X-Algolia-Secure-User-Token"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"role":"assistant","content":"ok"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewTryCmd(f, nil)
	_, err := test.Execute(cmd,
		"-c "+cfgPath+" -m hi --no-stream --no-cache --no-memory --no-analytics --secure-user-token ey.signed.jwt",
		out)
	require.NoError(t, err)
}

func Test_runTryCmd_RejectsInvalidCompatibility(t *testing.T) {
	cfgPath := writeTempJSON(t, "cfg.json", `{"model":"x"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewTryCmd(f, nil)
	_, err := test.Execute(cmd, "-c "+cfgPath+" -m hi --compatibility v9", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `invalid --compatibility "v9"`)
}
