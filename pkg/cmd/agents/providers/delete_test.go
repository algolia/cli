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

func Test_runDeleteCmd_PrefetchesBeforeDelete(t *testing.T) {
	var getCalls, deleteCalls int
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCalls++
			_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-x"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
		case http.MethodDelete:
			deleteCalls++
			w.WriteHeader(http.StatusNoContent)
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "delete p1 -y", out)
	require.NoError(t, err)
	assert.Equal(t, 1, getCalls)
	assert.Equal(t, 1, deleteCalls)
	assert.NotContains(t, result.String(), "Dry run:")
}

func Test_runDeleteCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "delete p1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runDeleteCmd_PropagatesConflict(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-x"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"detail":"Provider is in use by 3 agent(s)"}`))
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "delete p1 -y", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "in use by 3 agent")
}
