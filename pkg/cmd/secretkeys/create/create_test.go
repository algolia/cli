package create_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/secretkeys/create"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestCreateSecretKey(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(
		httpmock.REST("POST", "agent-studio/1/secret-keys"),
		func(req *http.Request) (*http.Response, error) {
			captured, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(
				`{"id":"key_1","name":"my-key","value":"secret_***","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"agentIds":["agent_1"]}`,
			)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "my-key --agent-ids agent_1", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"key_1"`)
	r.Verify(t)

	assert.JSONEq(t, `{"name":"my-key","agentIds":["agent_1"]}`, string(captured))
}

func TestCreateSecretKey_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <name> argument")
}
