package publish

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewPublishCmd_callsEndpoint(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodPost, "1/agents/a1/publish"),
		httpmock.JSONResponse(map[string]any{
			"id":           "a1",
			"name":         "Helper",
			"status":       "published",
			"instructions": "be helpful",
			"createdAt":    "2026-01-01T00:00:00Z",
			"updatedAt":    "2026-01-01T00:00:00Z",
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewPublishCmd(f, nil)
	_, err := test.Execute(cmd, "a1", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Published agent a1")
}
