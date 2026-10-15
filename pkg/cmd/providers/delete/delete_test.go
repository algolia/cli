package delete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/delete"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	providerID = "22222222-2222-4222-8222-222222222222"
	deletePath = "agent-studio/1/providers/" + providerID
)

func TestDeleteProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", deletePath), httpmock.StringResponse(""))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, providerID+" --confirm", out)
	require.NoError(t, err)
	assert.Empty(t, out.String())
	r.Verify(t)
}

func TestDeleteProvider_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, providerID, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
	assert.Empty(t, r.Requests)
}

func TestDeleteProvider_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "--confirm", out)
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"agentstudio providers delete requires a <provider-id> argument",
	)
}
