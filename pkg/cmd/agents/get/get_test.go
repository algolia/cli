package get

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/test"
)

func newCmdAgainst(
	t *testing.T,
	handler http.Handler,
) (*test.CmdInOut, func(args string) (*test.CmdInOut, error)) {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = func() (*agentstudio.Client, error) {
		return agentstudio.NewClient(agentstudio.Config{
			BaseURL:       ts.URL,
			ApplicationID: "APP123",
			APIKey:        "key-abc",
			HTTPClient:    ts.Client(),
		})
	}

	exec := func(args string) (*test.CmdInOut, error) {
		cmd := NewGetCmd(f, nil)
		return test.Execute(cmd, args, out)
	}
	return out, exec
}

func Test_runGetCmd_DefaultsToJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"id":"abc-123",
			"name":"Concierge",
			"status":"published",
			"instructions":"Be helpful.",
			"createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	_, exec := newCmdAgainst(t, mux)

	result, err := exec("abc-123")
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id":"abc-123"`)
	assert.Contains(t, result.String(), `"name":"Concierge"`)
	assert.Contains(t, result.String(), `"status":"published"`)
}

func Test_runGetCmd_WrapsNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})

	_, exec := newCmdAgainst(t, mux)

	_, err := exec("missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, agentstudio.ErrNotFound), "got %v", err)
}

func Test_runGetCmd_RejectsMissingArg(t *testing.T) {
	_, exec := newCmdAgainst(t, http.NewServeMux())

	_, err := exec("")
	require.Error(t, err) // cobra ExactArgs(1) rejects empty positional list
}
