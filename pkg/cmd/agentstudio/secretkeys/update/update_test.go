package update_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/secretkeys/update"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestUpdateSecretKey(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(
		httpmock.REST("PATCH", "agent-studio/1/secret-keys/my-key"),
		func(req *http.Request) (*http.Response, error) {
			captured, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(
				`{"id":"my-key","name":"renamed","value":"secret_***","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"agentIds":[]}`,
			)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "my-key --name renamed", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"name":"renamed"`)
	r.Verify(t)

	assert.JSONEq(t, `{"name":"renamed"}`, string(captured))
}

func TestUpdateSecretKey_RequiresAtLeastOneField(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "my-key", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one of `--name` or `--agent-ids`")
}
