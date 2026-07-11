package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/conversations/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestGetConversation(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents/my-agent/conversations/conv_1"),
		httpmock.StringResponse(
			`{"id":"conv_1","agentId":"my-agent","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","messages":[]}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-agent conv_1", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"conv_1"`)
	r.Verify(t)
}

func TestGetConversation_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> and a <conversation-id> argument")
}
