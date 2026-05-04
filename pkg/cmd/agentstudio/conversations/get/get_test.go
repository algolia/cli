package get

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewGetCmd_callsEndpoint(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodGet, "1/agents/a1/conversations/conv-1"),
		httpmock.JSONResponse(map[string]any{
			"id":        "conv-1",
			"agentId":   "a1",
			"title":     "Hello",
			"createdAt": "2026-01-01T00:00:00Z",
			"updatedAt": "2026-01-01T00:00:00Z",
			"messages":  []any{},
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewGetCmd(f, nil)
	_, err := test.Execute(cmd, "a1 conv-1", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), `"conv-1"`)
	assert.Contains(t, out.String(), `"Hello"`)
}
