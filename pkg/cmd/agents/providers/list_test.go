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

func Test_runListCmd_MasksSecretsByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"data":[{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-LEAKED-DATA"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}],
			"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "list --output json", out)
	require.NoError(t, err)
	got := result.String()
	assert.NotContains(t, got, "sk-LEAKED-DATA", "raw apiKey must NOT appear in default output")
	assert.Contains(t, got, `"apiKey":"***"`, "apiKey must be replaced with masking sentinel")
}

func Test_runListCmd_ShowSecretRevealsRawKey(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"data":[{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-INTENTIONAL-EXPORT"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}],
			"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "list --output json --show-secret", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "sk-INTENTIONAL-EXPORT")
}
