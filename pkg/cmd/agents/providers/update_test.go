package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runUpdateCmd_HitsCorrectPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"renamed","providerName":"openai",
			"input":{"apiKey":"sk-XYZ"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	patchPath := writeTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "renamed")
	assert.NotContains(t, result.String(), "sk-XYZ")
}

func Test_runUpdateCmd_DryRunIncludesProviderID(t *testing.T) {
	patchPath := writeTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath+" --dry-run", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "Dry run: would PATCH /1/providers/p1")
}
