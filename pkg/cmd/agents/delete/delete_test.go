package delete

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func agentJSON(name, status string) string {
	return `{
		"id":"abc-123","name":"` + name + `","status":"` + status + `",
		"instructions":"x","createdAt":"2025-01-01T00:00:00Z"
	}`
}

func writeTestJSONResponse(w http.ResponseWriter, body []byte) {
	var out io.Writer = w
	_, _ = out.Write(body)
}

func Test_runDeleteCmd_NonTTYRequiresConfirm(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "abc-123", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runDeleteCmd_NonTTYWithConfirmDeletes(t *testing.T) {
	var deleted atomic.Bool

	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeTestJSONResponse(w, []byte(agentJSON("Concierge", "draft")))
		case http.MethodDelete:
			deleted.Store(true)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "abc-123 -y", out)
	require.NoError(t, err)
	assert.True(t, deleted.Load(), "DELETE should have been called")
}

func Test_runDeleteCmd_PropagatesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "missing -y", out)
	require.Error(t, err)
}
