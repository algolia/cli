package create

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewCreateCmd_requiresNameAndInstructions(t *testing.T) {
	r := httpmock.Registry{}
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--name")
}

func Test_NewCreateCmd_postsBodyAndPrintsSuccess(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodPost, "1/agents"),
		httpmock.JSONResponse(map[string]any{
			"id":           "a1",
			"name":         "Helper",
			"status":       "draft",
			"instructions": "be helpful",
			"createdAt":    "2026-01-01T00:00:00Z",
			"updatedAt":    "2026-01-01T00:00:00Z",
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, `--name Helper --instructions "be helpful"`, out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Created agent Helper")
	assert.Contains(t, out.String(), "a1")
}
