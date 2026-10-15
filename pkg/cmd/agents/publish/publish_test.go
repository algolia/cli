package publish_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/publish"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const publishedAgentBody = `{
  "id": "agt_1",
  "name": "test-agent",
  "description": "d",
  "status": "published",
  "providerId": "p1",
  "instructions": "You are helpful.",
  "config": {},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "lastUsedAt": null
}`

func TestPublishAgent(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents/agt_1/publish"),
		httpmock.StringResponse(publishedAgentBody),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := publish.NewPublishCmd(f)
	_, err := test.Execute(cmd, "agt_1", out)
	require.NoError(t, err)

	assert.JSONEq(t, publishedAgentBody, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestPublishAgent_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := publish.NewPublishCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <agent-id> argument")
}
