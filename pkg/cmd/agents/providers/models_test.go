package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runModelsCmd_NoFlag_HitsCatalogRoute(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/providers/models", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"openai":["gpt-4o","gpt-4o-mini"],"anthropic":["claude-3-5-sonnet"]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "models --output json", out)
	require.NoError(t, err)
	assert.True(t, hit)

	var got map[string][]string
	require.NoError(t, json.Unmarshal([]byte(result.String()), &got))
	assert.Contains(t, got, "openai")
	assert.Contains(t, got["openai"], "gpt-4o")
}

func Test_runModelsCmd_WithFlag_HitsConfiguredProviderRoute(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/providers/p1/models", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`["gpt-4o","my-fine-tune-1"]`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "models --provider-id p1 --output json", out)
	require.NoError(t, err)
	assert.True(t, hit)
	assert.Contains(t, result.String(), "my-fine-tune-1")
}

func Test_runModelsCmd_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/missing/models", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Provider not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "models --provider-id missing", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Provider not found")
}
