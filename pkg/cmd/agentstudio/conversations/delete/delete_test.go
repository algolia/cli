package delete

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewDeleteCmd_requiresConfirm(t *testing.T) {
	r := httpmock.Registry{}
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "a1 conv-1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm")
}

func Test_NewDeleteCmd_rejectsAllPlusConvID(t *testing.T) {
	r := httpmock.Registry{}
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "a1 conv-1 --all --confirm", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--all conflicts")
}

func Test_NewDeleteCmd_singleConversation(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodDelete, "1/agents/a1/conversations/conv-1"),
		httpmock.JSONResponse(""),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "a1 conv-1 --confirm", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Deleted conversation conv-1")
}

func Test_NewDeleteCmd_allConversations(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodDelete, "1/agents/a1/conversations"),
		httpmock.JSONResponse(""),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(true, &r, nil, "")
	cmd := NewDeleteCmd(f, nil)
	_, err := test.Execute(cmd, "a1 --all --confirm", out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Deleted all conversations for agent a1")
}
