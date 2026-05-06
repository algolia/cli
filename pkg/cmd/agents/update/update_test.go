package update

import (
	"encoding/json"
	"io"
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

func Test_runUpdateCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"Renamed"}`, string(body))

		_, _ = w.Write([]byte(`{
			"id":"abc-123","name":"Renamed","status":"draft",
			"instructions":"x","createdAt":"2025-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, `{"name":"Renamed"}`)
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewUpdateCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123 -F -", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"name":"Renamed"`)
}

func Test_runUpdateCmd_DryRunStructuredIncludesAgentID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, `{"name":"X"}`)
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewUpdateCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123 -F - --dry-run --output json", out)
	require.NoError(t, err)

	var summary map[string]any
	require.NoError(t, json.Unmarshal([]byte(result.String()), &summary))
	assert.Equal(t, "update_agent", summary["action"])
	assert.Equal(t, "PATCH /1/agents/abc-123", summary["request"])
	assert.Equal(t, "abc-123", summary["agentId"])
	assert.Equal(t, true, summary["dryRun"])
}

func Test_runUpdateCmd_RequiresAgentID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, `{"name":"X"}`)
	cmd := NewUpdateCmd(f, nil)
	_, err := test.Execute(cmd, "-F -", out)
	require.Error(t, err)
	// validators.ExactArgs(1) message — not pinning exact wording, just that it errors.
}
