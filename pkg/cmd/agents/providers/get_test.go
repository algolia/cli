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

func Test_runGetCmd_MasksByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/agent-studio/1/providers/p1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-LEAKED"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioAPIClient = sharedtest.NewAPIClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "get p1", out)
	require.NoError(t, err)
	assert.NotContains(t, result.String(), "sk-LEAKED")
	assert.Contains(t, result.String(), `"apiKey":"***"`)
}

func Test_runGetCmd_RequiresProviderID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "get", out)
	require.Error(t, err)
}

// azure_openai / openai_compatible inputs carry fields the typed SDK oneOf
// drops; they must survive, with only the secret masked.
func Test_runGetCmd_KeepsProviderSpecificInputFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/agent-studio/1/providers/p1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"p1","name":"az","providerName":"azure_openai","input":{"apiKey":"sk-SECRET","azureEndpoint":"https://x.openai.azure.com","azureDeployment":"dep","apiVersion":"2024-01-01"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioAPIClient = sharedtest.NewAPIClient(t, ts)
	result, err := test.Execute(NewProvidersCmd(f), "get p1", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"azureEndpoint":"https://x.openai.azure.com"`)
	assert.Contains(t, got, `"azureDeployment":"dep"`)
	assert.NotContains(t, got, "sk-SECRET")
}
