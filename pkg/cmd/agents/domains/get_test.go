package domains

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runGetCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/d1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"d1","appId":"APP","agentId":"agent-1","domain":"x","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "get agent-1 d1", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id"`)
}

func Test_runGetCmd_RequiresBoth(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "get agent-1", out)
	require.Error(t, err)
}
