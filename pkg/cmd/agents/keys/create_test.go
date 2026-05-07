package keys

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runCreateCmd_RequiresName(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewKeysCmd(f)
	_, err := test.Execute(cmd, "create", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--name is required")
}

func Test_runCreateCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "create --name k1 --agent-id a1 --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would POST /1/secret-keys")
	assert.Contains(t, got, `"name": "k1"`)
}

func Test_runCreateCmd_LiveMasksValue(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "k1", got["name"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(
			[]byte(
				`{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "create --name k1", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"***"`)
	assert.NotContains(t, got, "sk-real")
}
