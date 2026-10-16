package create_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/create"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const createdAgentBody = `{
  "id": "agt_1",
  "name": "cli-test-agent",
  "description": null,
  "status": "draft",
  "providerId": null,
  "instructions": "You are a helpful assistant.",
  "config": {},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "lastUsedAt": null
}`

func TestCreateAgent_FromFlags(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("POST", "agent-studio/1/agents"), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(createdAgentBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, `--name cli-test-agent --instructions "You are a helpful assistant."`, out)
	require.NoError(t, err)

	assert.JSONEq(t, createdAgentBody, strings.TrimSpace(out.String()))
	assert.JSONEq(t, `{"name":"cli-test-agent","instructions":"You are a helpful assistant."}`, string(gotBody))
	r.Verify(t)
}

func TestCreateAgent_FromJSON(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("POST", "agent-studio/1/agents"), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(createdAgentBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(
		cmd,
		`--json '{"name":"cli-test-agent","instructions":"You are a helpful assistant.","model":"gpt-4o"}'`,
		out,
	)
	require.NoError(t, err)

	assert.JSONEq(t, createdAgentBody, strings.TrimSpace(out.String()))
	assert.JSONEq(
		t,
		`{"name":"cli-test-agent","instructions":"You are a helpful assistant.","model":"gpt-4o"}`,
		string(gotBody),
	)
	r.Verify(t)
}

func TestCreateAgent_MissingInput(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`--json` or both `--name` and `--instructions` are required")
}

func TestCreateAgent_JSONAndFlagsConflict(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, `--json '{"name":"x","instructions":"y"}' --name x --instructions y`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--json cannot be combined")
}

func TestCreateAgent_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "--json not-json", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing agent JSON")
}

func TestCreateAgent_JSONMissingRequiredFields(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, `--json '{"name":"only-name"}'`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires `name` and `instructions`")
}
