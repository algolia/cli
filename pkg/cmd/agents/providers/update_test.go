package providers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
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

	patchPath := sharedtest.WriteTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "renamed")
	assert.NotContains(t, result.String(), "sk-XYZ")
}

func Test_runUpdateCmd_DryRunIncludesProviderID(t *testing.T) {
	patchPath := sharedtest.WriteTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath+" --dry-run", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "Dry run: would PATCH /1/providers/p1")
}

func Test_runUpdateCmd_Flags_PatchesRename(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p9", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"new-label"}`, string(bytes.TrimSpace(body)))
		_, _ = w.Write([]byte(`{
			"id":"p9","name":"new-label","providerName":"openai",
			"input":{"apiKey":"sk-XYZ"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(
		cmd,
		`update p9 --name new-label`,
		out,
	)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "new-label")
}

func Test_runUpdateCmd_FlagsEnvRotatesKey(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p9", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"input":{"apiKey":"sk-env-rot"}}`, string(bytes.TrimSpace(body)))
		_, _ = w.Write([]byte(`{
			"id":"p9","name":"x","providerName":"openai",
			"input":{"apiKey":"sk-env-rot"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	t.Setenv("ROT_KEY", "sk-env-rot")

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, `update p9 --api-key-env ROT_KEY`, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"apiKey":"***"`)
}

func Test_runUpdateCmd_FlagsRejectWithFile(t *testing.T) {
	p := sharedtest.WriteTempJSON(t, "patch.json", `{"name":"n"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	cli := `update p9 --name clash -F ` + p
	_, err := test.Execute(cmd, cli, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "combine")
}
