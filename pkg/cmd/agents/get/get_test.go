package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const getAgentBody = `{
  "id": "agt_1",
  "name": "test-agent",
  "description": "d",
  "status": "draft",
  "providerId": "p1",
  "instructions": "You are helpful.",
  "config": {},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "lastUsedAt": null
}`

func TestGetAgent(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents/agt_1"),
		httpmock.StringResponse(getAgentBody),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "agt_1", out)
	require.NoError(t, err)

	assert.JSONEq(t, getAgentBody, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestGetAgent_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agentstudio agents get requires a <agent-id> argument")
}
