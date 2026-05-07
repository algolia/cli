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

func Test_runUpdateCmd_RequiresAtLeastOneField(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewKeysCmd(f)
	_, err := test.Execute(cmd, "update id1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to update")
}

func Test_runUpdateCmd_DryRunNameOnly(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, `update id1 --name renamed --dry-run`, out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "PATCH /1/secret-keys/id1")
	assert.Contains(t, got, `"name": "renamed"`)
	assert.NotContains(t, got, "agentIds")
}

func Test_runUpdateCmd_LivePatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "renamed", got["name"])
		_, _ = w.Write(
			[]byte(
				`{"id":"id1","name":"renamed","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, `update id1 --name renamed`, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"***"`)
}
