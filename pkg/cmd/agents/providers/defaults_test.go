package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runDefaultsCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/models/defaults", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"openai":"gpt-4.1-mini","anthropic":"claude-haiku-4-5"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "defaults --output json", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "gpt-4.1-mini")
	assert.Contains(t, got, "claude-haiku-4-5")
}

func Test_runDefaultsCmd_RejectsArgs(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "defaults extra-arg", out)
	require.Error(t, err)
}
