package conversations

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runListCmd_PassesFiltersThrough(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "2026-01-01", q.Get("startDate"))
		assert.Equal(t, "2026-01-31", q.Get("endDate"))
		assert.Equal(t, "true", q.Get("includeFeedback"))
		assert.Equal(t, "0", q.Get("feedbackVote"))
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":20,"totalCount":0,"totalPages":0}}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd,
		"list agent-1 --start-date 2026-01-01 --end-date 2026-01-31 --include-feedback --feedback-vote 0 --output json",
		out)
	require.NoError(t, err)
}

func Test_runListCmd_RejectsFeedbackVoteWithoutInclude(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "list agent-1 --feedback-vote 1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires --include-feedback")
}

func Test_runListCmd_RejectsBadFeedbackVote(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "list agent-1 --include-feedback --feedback-vote 2", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be 0 (down) or 1 (up)")
}

func Test_runListCmd_RequiresAgentID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "list", out)
	require.Error(t, err)
}
