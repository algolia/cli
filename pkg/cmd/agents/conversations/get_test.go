package conversations

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runGetCmd_PrintsBackendBody(t *testing.T) {
	body := `{"id":"c1","agentId":"agent-1","createdAt":"2026-01-15T00:00:00Z","updatedAt":"2026-01-15T00:01:00Z","messages":[{"role":"user","content":[{"type":"text","text":"hi"}]}]}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations/c1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("includeFeedback"))
		_, _ = w.Write([]byte(body))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewConversationsCmd(f)
	result, err := test.Execute(cmd, "get agent-1 c1 --include-feedback --output json", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id"`)
	assert.Contains(t, result.String(), `"messages"`)
}

func Test_runGetCmd_RequiresBothIDs(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "get agent-1", out)
	require.Error(t, err)
}
