package delete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/conversations/delete"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	agentID    = "11111111-1111-4111-8111-111111111111"
	convID     = "33333333-3333-4333-8333-333333333333"
	deletePath = "agent-studio/1/agents/" + agentID + "/conversations/" + convID
)

func TestDeleteConversation(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", deletePath), httpmock.StringResponse(""))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, agentID+" "+convID+" --confirm", out)
	require.NoError(t, err)
	assert.Empty(t, out.String())
	r.Verify(t)
}

func TestDeleteConversation_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, agentID+" "+convID, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
	assert.Empty(t, r.Requests)
}

func TestDeleteConversation_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, agentID+" --confirm", out)
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"agentstudio conversations delete requires <agent-id> and <conversation-id>",
	)
}
