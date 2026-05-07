package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runCreateCmd_RoundTripsBody(t *testing.T) {
	body := `{"name":"openai-prod","providerName":"openai","input":{"apiKey":"sk-XYZ"}}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-XYZ"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := sharedtest.WriteTempJSON(t, "spec.json", body)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath, out)
	require.NoError(t, err)
	assert.NotContains(t, result.String(), "sk-XYZ")
	assert.Contains(t, result.String(), `"apiKey":"***"`)
}

func Test_runCreateCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := sharedtest.WriteTempJSON(
		t,
		"spec.json",
		`{"name":"x","providerName":"openai","input":{"apiKey":"sk-x"}}`,
	)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath+" --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would POST /1/providers")
	// Dry-run shows the unmodified body — masking does NOT apply to
	// dry-run because the user authored the file and is being shown
	// what THEY are about to send.
	assert.Contains(t, got, "sk-x")
}

func Test_runCreateCmd_RejectsInvalidJSON(t *testing.T) {
	specPath := sharedtest.WriteTempJSON(t, "spec.json", `{not json`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "create -F "+specPath+" --dry-run", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}
