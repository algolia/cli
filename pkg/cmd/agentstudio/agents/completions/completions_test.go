package completions_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/completions"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestCompletions_Message(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents/my-agent/completions"),
		func(req *http.Request) (*http.Response, error) {
			captured, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(`{"id":"resp_1"}`)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := completions.NewCompletionsCmd(f)
	_, err := test.Execute(cmd, `my-agent --message "hello there"`, out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"resp_1"`)
	r.Verify(t)

	var body map[string]any
	require.NoError(t, json.Unmarshal(captured, &body))
	messages, ok := body["messages"].([]any)
	require.True(t, ok, "expected messages array in request body: %s", captured)
	require.Len(t, messages, 1)
	msg := messages[0].(map[string]any)
	assert.Equal(t, "user", msg["role"])
}

func TestCompletions_RequiresMessageOrFile(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := completions.NewCompletionsCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one of `--message` or `--file`")
}

func TestCompletions_MessageAndFileConflict(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := completions.NewCompletionsCmd(f)
	_, err := test.Execute(cmd, `my-agent --message "hi" --file body.json`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one of `--message` or `--file`")
}

func TestCompletions_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := completions.NewCompletionsCmd(f)
	_, err := test.Execute(cmd, "--message hi", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> argument")
}
