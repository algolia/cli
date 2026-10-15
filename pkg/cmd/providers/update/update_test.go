package update_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/update"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	providerID          = "22222222-2222-4222-8222-222222222222"
	updatePath          = "agent-studio/1/providers/" + providerID
	updatedProviderBody = `{
  "id": "22222222-2222-4222-8222-222222222222",
  "name": "anthropic-prod",
  "providerName": "anthropic",
  "input": {"apiKey": "sk-test"},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z"
}`
)

func TestUpdateProvider_FromFlags(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("PATCH", updatePath), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(updatedProviderBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID+` --name anthropic-prod`, out)
	require.NoError(t, err)

	assert.JSONEq(t, updatedProviderBody, strings.TrimSpace(out.String()))
	assert.JSONEq(t, `{"name":"anthropic-prod"}`, string(gotBody))
	r.Verify(t)
}

func TestUpdateProvider_FromJSON(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(httpmock.REST("PATCH", updatePath), func(req *http.Request) (*http.Response, error) {
		gotBody, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(updatedProviderBody)(req)
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID+` --json '{"input":{"apiKey":"sk-test"}}'`, out)
	require.NoError(t, err)

	assert.JSONEq(t, updatedProviderBody, strings.TrimSpace(out.String()))
	assert.Contains(t, string(gotBody), `"apiKey":"sk-test"`)
	r.Verify(t)
}

func TestUpdateProvider_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "--name x", out)
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"agentstudio providers update requires a <provider-id> argument",
	)
}

func TestUpdateProvider_MissingInput(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`--json` or `--name` is required")
}

func TestUpdateProvider_JSONAndFlagsConflict(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID+` --json '{"name":"x"}' --name x`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--json cannot be combined")
}

func TestUpdateProvider_EmptyJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID+` --json '{}'`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires `name` and/or `input`")
}

func TestUpdateProvider_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, providerID+" --json not-json", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing provider JSON")
}
