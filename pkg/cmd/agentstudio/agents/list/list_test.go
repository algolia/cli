package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestListAgents(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents"),
		httpmock.StringResponse(
			`{"data":[{"id":"agent_1","name":"Agent One","description":null,"status":"published","providerId":null,"instructions":"help","config":{},"createdAt":"2026-01-01T00:00:00Z","updatedAt":null,"lastUsedAt":null}],"pagination":{"page":0,"limit":20,"total":1}}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"agent_1"`)
	r.Verify(t)
}

func TestListAgents_WithFilters(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents"),
		httpmock.StringResponse(`{"data":[],"pagination":{"page":2,"limit":10,"total":0}}`),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "--page 2 --limit 10 --provider-id my-provider", out)
	require.NoError(t, err)
	r.Verify(t)

	require.Len(t, r.Requests, 1)
	q := r.Requests[0].URL.Query()
	assert.Equal(t, "2", q.Get("page"))
	assert.Equal(t, "10", q.Get("limit"))
	assert.Equal(t, "my-provider", q.Get("providerId"))
}
