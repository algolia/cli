package list

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewListCmd_rendersTable(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodGet, "1/agents/a1/conversations"),
		httpmock.JSONResponse(map[string]any{
			"data": []any{map[string]any{
				"id":           "conv-1",
				"agentId":      "a1",
				"title":        "Hello",
				"messageCount": 3,
				"createdAt":    "2026-01-01T00:00:00Z",
				"updatedAt":    "2026-01-01T00:00:00Z",
			}},
			"pagination": map[string]any{"page": 1, "limit": 10, "totalCount": 1, "totalPages": 1},
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewListCmd(f, nil)
	_, err := test.Execute(cmd, "a1", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "conv-1")
	assert.Contains(t, out.String(), "Hello")
	assert.Contains(t, out.String(), "3")
}
