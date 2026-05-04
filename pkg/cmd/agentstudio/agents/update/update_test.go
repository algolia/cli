package update

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewUpdateCmd_appliesFlagAsBody(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodPatch, "1/agents/a1"),
		httpmock.JSONResponse(map[string]any{
			"id":           "a1",
			"name":         "Helper",
			"description":  "renamed via CLI",
			"status":       "draft",
			"instructions": "be helpful",
			"createdAt":    "2026-01-01T00:00:00Z",
			"updatedAt":    "2026-01-01T00:00:00Z",
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewUpdateCmd(f, nil)
	_, err := test.Execute(cmd, `a1 --description "renamed via CLI"`, out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Updated agent a1")
}
