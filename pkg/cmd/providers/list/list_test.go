package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const listProvidersBody = `{
  "data": [
    {
      "id": "22222222-2222-4222-8222-222222222222",
      "name": "openai",
      "providerName": "openai",
      "input": {"apiKey": "sk-test"},
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

const emptyListProvidersBody = `{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalCount": 0,
    "totalPages": 0
  }
}`

func TestListProviders(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
		wantQ   map[string]string
	}{
		{
			name:    "single provider",
			mock:    listProvidersBody,
			wantOut: listProvidersBody,
		},
		{
			name:    "pagination",
			cli:     "--page 2 --limit 20",
			mock:    emptyListProvidersBody,
			wantOut: emptyListProvidersBody,
			wantQ: map[string]string{
				"page":  "2",
				"limit": "20",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "agent-studio/1/providers"),
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
