package delete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/secretkeys/delete"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestDeleteSecretKey(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", "agent-studio/1/secret-keys/my-key"), httpmock.StringResponse(``))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-key --confirm", out)
	require.NoError(t, err)
	r.Verify(t)
}

func TestDeleteSecretKey_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-key", out)
	require.Error(t, err)
	assert.Empty(t, r.Requests)
}
