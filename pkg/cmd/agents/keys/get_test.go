package keys

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runGetCmd_MasksValueByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "get id1", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"***"`)
	assert.NotContains(t, got, "sk-real")
}
