package export_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/conversations/export"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	agentID    = "11111111-1111-4111-8111-111111111111"
	exportPath = "agent-studio/1/agents/" + agentID + "/conversations/export"
	exportBody = `[{
  "id": "33333333-3333-4333-8333-333333333333",
  "agentId": "11111111-1111-4111-8111-111111111111",
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "messages": []
}]`
)

func TestExportConversations(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", exportPath), httpmock.StringResponse(exportBody))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := export.NewExportCmd(f)
	_, err := test.Execute(cmd, agentID+" --start-date 2026-01-01 --end-date 2026-01-31", out)
	require.NoError(t, err)

	assert.JSONEq(t, exportBody, strings.TrimSpace(out.String()))
	r.Verify(t)

	require.Len(t, r.Requests, 1)
	q := r.Requests[0].URL.Query()
	assert.Equal(t, "2026-01-01", q.Get("startDate"))
	assert.Equal(t, "2026-01-31", q.Get("endDate"))
}

func TestExportConversations_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := export.NewExportCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"agentstudio conversations export requires a <agent-id> argument",
	)
}
