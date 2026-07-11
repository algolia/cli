package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/conversations/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestListConversations(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents/my-agent/conversations"),
		httpmock.StringResponse(
			`{"data":[{"id":"conv_1","agentId":"my-agent","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}],"pagination":{"page":0,"limit":20,"total":1}}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"conv_1"`)
	r.Verify(t)
}

func TestListConversations_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> argument")
}
