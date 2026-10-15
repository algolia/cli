package update_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/update"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const updatedAgentBody = `{
  "id": "agt_1",
  "name": "renamed-bot",
  "description": null,
  "status": "draft",
  "providerId": null,
  "instructions": "You are a helpful assistant.",
  "config": {},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "lastUsedAt": null
}`

func TestUpdateAgent_FromFlags(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("PATCH", "agent-studio/1/agents/agt_1"), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(updatedAgentBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "agt_1 --name renamed-bot", out)
	require.NoError(t, err)

	assert.JSONEq(t, updatedAgentBody, strings.TrimSpace(out.String()))
	assert.JSONEq(t, `{"name":"renamed-bot"}`, string(gotBody))
	r.Verify(t)
}

func TestUpdateAgent_FromFile(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("PATCH", "agent-studio/1/agents/agt_1"), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(updatedAgentBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	out.InBuf.WriteString(`{"name":"renamed-bot","instructions":"Be brief."}`)
	_, err := test.Execute(cmd, "agt_1 --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, updatedAgentBody, strings.TrimSpace(out.String()))
	assert.JSONEq(t, `{"name":"renamed-bot","instructions":"Be brief."}`, string(gotBody))
	r.Verify(t)
}

func TestUpdateAgent_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "--name x", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <agent-id> argument")
}

func TestUpdateAgent_MissingInput(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "agt_1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`--file` or at least one of `--name` or `--instructions` is required")
}

func TestUpdateAgent_FileAndFlagsConflict(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "agt_1 --file patch.json --name x", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--file cannot be combined")
}

func TestUpdateAgent_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	out.InBuf.WriteString("not-json")
	_, err := test.Execute(cmd, "agt_1 --file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing agent JSON")
}
