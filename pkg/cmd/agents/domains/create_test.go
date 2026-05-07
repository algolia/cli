package domains

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runCreateCmd_RequiresDomain(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "create agent-1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--domain is required")
}

func Test_runCreateCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "create agent-1 --domain https://x.test --dry-run", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "Dry run: would POST /1/agents/agent-1/allowed-domains")
	assert.Contains(t, result.String(), `"https://x.test"`)
}

func Test_runCreateCmd_Live(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(
			[]byte(
				`{"id":"d1","appId":"APP","agentId":"agent-1","domain":"https://x.test","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "create agent-1 --domain https://x.test", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id"`)
}
