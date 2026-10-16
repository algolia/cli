package purge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/conversations/purge"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	agentID   = "11111111-1111-4111-8111-111111111111"
	purgePath = "agent-studio/1/agents/" + agentID + "/conversations"
)

func TestPurgeConversations(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", purgePath), httpmock.StringResponse(""))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := purge.NewPurgeCmd(f)
	_, err := test.Execute(cmd, agentID+" --end-date 2026-01-01 --confirm", out)
	require.NoError(t, err)
	assert.Empty(t, out.String())
	r.Verify(t)

	require.Len(t, r.Requests, 1)
	assert.Equal(t, "2026-01-01", r.Requests[0].URL.Query().Get("endDate"))
}

func TestPurgeConversations_RequiresDate(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := purge.NewPurgeCmd(f)
	_, err := test.Execute(cmd, agentID+" --confirm", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`--start-date` or `--end-date` is required")
	assert.Empty(t, r.Requests)
}

func TestPurgeConversations_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := purge.NewPurgeCmd(f)
	_, err := test.Execute(cmd, agentID+" --end-date 2026-01-01", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
	assert.Empty(t, r.Requests)
}
