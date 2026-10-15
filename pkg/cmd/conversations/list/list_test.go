package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/conversations/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	agentID  = "11111111-1111-4111-8111-111111111111"
	listPath = "agent-studio/1/agents/" + agentID + "/conversations"
	listBody = `{
  "data": [
    {
      "id": "33333333-3333-4333-8333-333333333333",
      "agentId": "11111111-1111-4111-8111-111111111111",
      "createdAt": "2026-01-01T00:00:00Z",
      "updatedAt": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalCount": 1,
    "totalPages": 1
  }
}`
	emptyListBody = `{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalCount": 0,
    "totalPages": 0
  }
}`
)

func TestListConversations(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
		wantQ   map[string]string
	}{
		{
			name:    "single conversation",
			cli:     agentID,
			mock:    listBody,
			wantOut: listBody,
		},
		{
			name:    "pagination and dates",
			cli:     agentID + " --page 2 --limit 20 --start-date 2026-01-01 --end-date 2026-01-31 --include-feedback",
			mock:    emptyListBody,
			wantOut: emptyListBody,
			wantQ: map[string]string{
				"page":            "2",
				"limit":           "20",
				"startDate":       "2026-01-01",
				"endDate":         "2026-01-31",
				"includeFeedback": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(httpmock.REST("GET", listPath), httpmock.StringResponse(tt.mock))

			f, out := test.NewFactory(false, r, nil, "")
			cmd := list.NewListCmd(f)
			_, err := test.Execute(cmd, tt.cli, out)
			require.NoError(t, err)

			assert.JSONEq(t, tt.wantOut, strings.TrimSpace(out.String()))
			r.Verify(t)

			if tt.wantQ != nil {
				require.Len(t, r.Requests, 1)
				q := r.Requests[0].URL.Query()
				for k, v := range tt.wantQ {
					assert.Equal(t, v, q.Get(k), "query %s", k)
				}
			}
		})
	}
}

func TestListConversations_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agentstudio conversations list requires a <agent-id> argument")
}
