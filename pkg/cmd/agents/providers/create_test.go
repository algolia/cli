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

func Test_runCreateCmd_RejectsInvalidJSON(t *testing.T) {
	specPath := sharedtest.WriteTempJSON(t, "spec.json", `{not json`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "create -F "+specPath, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func Test_runCreateCmd_RequiresFileFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "create", out)
	require.Error(t, err)
}

func Test_runCreateCmd_File_PostsOpenAIProviderBody(t *testing.T) {
	specJSON := `{"name":"prod","providerName":"openai","input":{"apiKey":"sk-env"}}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"input":{"apiKey":"sk-env"},"name":"prod","providerName":"openai"}`,
			string(bytes.TrimSpace(body)))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"prod","providerName":"openai",
			"input":{"apiKey":"sk-env"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := sharedtest.WriteTempJSON(t, "spec.json", specJSON)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"name":"prod"`)
}

func Test_runCreateCmd_FileAcceptsAzureOpenAI(t *testing.T) {
	specJSON := `{"name":"azure1","providerName":"azure_openai","input":{"apiKey":"k","azureEndpoint":"https://x.openai.azure.com","azureDeployment":"d"}}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"azure1","providerName":"azure_openai",
			"input":{"apiKey":"k"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := sharedtest.WriteTempJSON(t, "spec.json", specJSON)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"name":"azure1"`)
}
