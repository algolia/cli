package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const listAgentsBody = `{
  "data": [
    {
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
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalCount": 1,
    "totalPages": 1
  }
}`

const emptyListAgentsBody = `{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalCount": 0,
    "totalPages": 0
  }
}`

func TestListAgents(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
		wantQ   map[string]string
	}{
		{
			name:    "single agent",
			mock:    listAgentsBody,
			wantOut: listAgentsBody,
		},
		{
			name:    "empty list",
			mock:    emptyListAgentsBody,
			wantOut: emptyListAgentsBody,
		},
		{
			name:    "pagination and provider filter",
			cli:     "--page 2 --limit 20 --provider-id p1",
			mock:    emptyListAgentsBody,
			wantOut: emptyListAgentsBody,
			wantQ: map[string]string{
				"page":       "2",
				"limit":      "20",
				"providerId": "p1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "agent-studio/1/agents"),
				httpmock.StringResponse(tt.mock),
			)

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
