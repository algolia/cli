package complete

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewCompleteCmd_rejectsBadCompatibilityMode(t *testing.T) {
	r := httpmock.Registry{}
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewCompleteCmd(f, nil)
	_, err := test.Execute(cmd, `a1 --compatibility-mode bogus`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ai-sdk-4, ai-sdk-5")
}

func Test_NewCompleteCmd_postsBodyAndStreamsResponse(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodPost, "1/agents/a1/completions"),
		httpmock.JSONResponse(map[string]any{"answer": "hi"}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewCompleteCmd(f, nil)
	_, err := test.Execute(cmd, `a1 --id conv-1`, out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), `"answer"`)
	assert.Contains(t, out.String(), `"hi"`)
}
