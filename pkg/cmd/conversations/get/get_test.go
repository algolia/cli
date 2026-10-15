package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/conversations/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	agentID = "11111111-1111-4111-8111-111111111111"
	convID  = "33333333-3333-4333-8333-333333333333"
	getPath = "agent-studio/1/agents/" + agentID + "/conversations/" + convID
	getBody = `{
  "id": "33333333-3333-4333-8333-333333333333",
  "agentId": "11111111-1111-4111-8111-111111111111",
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "messages": []
}`
)

func TestGetConversation(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", getPath), httpmock.StringResponse(getBody))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(
		cmd,
		agentID+" "+convID+" --include-feedback --include-message-events",
		out,
	)
	require.NoError(t, err)

	assert.JSONEq(t, getBody, strings.TrimSpace(out.String()))
	r.Verify(t)

	require.Len(t, r.Requests, 1)
	q := r.Requests[0].URL.Query()
	assert.Equal(t, "true", q.Get("includeFeedback"))
	assert.Equal(t, "true", q.Get("includeMessageEvents"))
}

func TestGetConversation_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, agentID, out)
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"agentstudio conversations get requires <agent-id> and <conversation-id>",
	)
}
