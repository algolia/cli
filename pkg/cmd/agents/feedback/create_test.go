package feedback

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runCreateCmd_RequiresAllRequired(t *testing.T) {
	cases := []struct {
		argv string
		msg  string
	}{
		{"create", "--agent-id is required"},
		{"create --agent-id a", "--message-id is required"},
		{"create --agent-id a --message-id m", "--vote is required"},
		{"create --agent-id a --message-id m --vote 5", "must be 0 or 1"},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			f, out := test.NewFactory(false, nil, nil, "")
			cmd := NewFeedbackCmd(f)
			_, err := test.Execute(cmd, tc.argv, out)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func Test_runCreateCmd_TooManyTags(t *testing.T) {
	tags := strings.Repeat("x,", 11)
	tags = strings.TrimSuffix(tags, ",")
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewFeedbackCmd(f)
	_, err := test.Execute(cmd, "create --agent-id a --message-id m --vote 1 --tags "+tags, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at most 10")
}

func Test_runCreateCmd_NotesTooLong(t *testing.T) {
	notes := strings.Repeat("x", 1001)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewFeedbackCmd(f)
	_, err := test.Execute(cmd, "create --agent-id a --message-id m --vote 0 --notes "+notes, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1000-character")
}

// Limits are in characters, not bytes: 20 three-byte runes (60 bytes) is a
// valid tag, 51 runes is not.
func Test_runCreateCmd_TagLimitCountsRunes(t *testing.T) {
	const argv = "create --agent-id a --message-id m --vote 1 --tags "

	mux := http.NewServeMux()
	mux.HandleFunc("/agent-studio/1/feedback", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(
			[]byte(
				`{"id":"fb1","agentId":"a","messageId":"m","vote":1,"tags":[],"notes":null,"model":null,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioAPIClient = sharedtest.NewAPIClient(t, ts)
	_, err := test.Execute(NewFeedbackCmd(f), argv+strings.Repeat("日", 20), out)
	require.NoError(t, err)

	f, out = test.NewFactory(false, nil, nil, "")
	_, err = test.Execute(NewFeedbackCmd(f), argv+strings.Repeat("日", 51), out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "50-character")
}

func Test_runCreateCmd_Live(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/agent-studio/1/feedback", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "m1", got["messageId"])
		assert.EqualValues(t, 1, got["vote"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(
			[]byte(
				`{"id":"fb1","agentId":"a1","messageId":"m1","vote":1,"tags":[],"notes":null,"model":null,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioAPIClient = sharedtest.NewAPIClient(t, ts)
	cmd := NewFeedbackCmd(f)
	result, err := test.Execute(cmd, "create --agent-id a1 --message-id m1 --vote 1", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id"`)
}
