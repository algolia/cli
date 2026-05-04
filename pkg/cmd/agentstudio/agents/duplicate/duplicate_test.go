package duplicate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewDuplicateCmd_postsAndPrintsResult(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodPost, "1/agents/a1/duplicate"),
		httpmock.JSONResponse(map[string]any{
			"id":           "a1-copy",
			"name":         "Helper (copy)",
			"status":       "draft",
			"instructions": "be helpful",
			"createdAt":    "2026-01-01T00:00:00Z",
			"updatedAt":    "2026-01-01T00:00:00Z",
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewDuplicateCmd(f, nil)
	_, err := test.Execute(cmd, "a1", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Duplicated a1 -> a1-copy")
}
