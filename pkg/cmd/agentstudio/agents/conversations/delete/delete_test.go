package delete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/conversations/delete"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestDeleteConversation(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("DELETE", "agent-studio/1/agents/my-agent/conversations/conv_1"),
		httpmock.StringResponse(``),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-agent conv_1 --confirm", out)
	require.NoError(t, err)
	r.Verify(t)
}

func TestDeleteConversation_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-agent conv_1", out)
	require.Error(t, err)
	assert.Empty(t, r.Requests)
}
