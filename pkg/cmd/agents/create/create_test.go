package create_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/create"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const validAgentBody = `{"name":"My Agent","instructions":"be nice"}`

func TestCreateAgent(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents"),
		httpmock.StringResponse(
			`{"id":"agent_1","name":"My Agent","description":null,"status":"draft","providerId":null,"instructions":"be nice","config":{},"createdAt":"2026-01-01T00:00:00Z","updatedAt":null,"lastUsedAt":null}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	out.InBuf.WriteString(validAgentBody)
	_, err := test.Execute(cmd, "--file -", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"agent_1"`)
	r.Verify(t)
}

func TestCreateAgent_MissingFile(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"file\" not set")
}

func TestCreateAgent_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	out.InBuf.WriteString("not-json")
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing agent JSON")
}
