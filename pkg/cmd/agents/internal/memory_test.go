package internal

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runMemoryCmd_RequiresBody(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewInternalCmd(f)
	_, err := test.Execute(cmd, "memorize agent-1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--body or --file")
}

func Test_runMemoryCmd_BodyAndFileMutex(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewInternalCmd(f)
	_, err := test.Execute(cmd, `memorize agent-1 --body {} -F somefile`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
}

func Test_runMemoryCmd_RejectsInvalidJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewInternalCmd(f)
	_, err := test.Execute(cmd, `memorize agent-1 --body not-json`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func Test_runMemoryCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agents/agent-1/memorize", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewInternalCmd(f)
	result, err := test.Execute(cmd,
		`memorize agent-1 --body '{"providerID":"p","model":"m","messages":[]}' --dry-run`, out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "POST /1/agents/agents/agent-1/memorize")
	assert.Contains(t, got, `"providerID"`)
}

func Test_runMemoryCmd_LiveDoubledPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agents/agent-1/ponder", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"providerID":"p","model":"m","messages":[]}`, string(body))
		_, _ = w.Write([]byte(`{"savedMemories":[],"message":"ok"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewInternalCmd(f)
	result, err := test.Execute(cmd,
		`ponder agent-1 --body '{"providerID":"p","model":"m","messages":[]}'`, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"ok"`)
}
