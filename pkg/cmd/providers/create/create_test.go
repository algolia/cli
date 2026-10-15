package create_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/create"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const createdProviderBody = `{
  "id": "22222222-2222-4222-8222-222222222222",
  "name": "anthropic-prod",
  "providerName": "anthropic",
  "input": {"apiKey": "sk-test"},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z"
}`

func TestCreateProvider_FromJSON(t *testing.T) {
	r := &httpmock.Registry{}
	var gotBody []byte
	r.Register(
		httpmock.REST("POST", "agent-studio/1/providers"),
		func(req *http.Request) (*http.Response, error) {
			gotBody, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(createdProviderBody)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(
		cmd,
		`--json '{"name":"anthropic-prod","providerName":"anthropic","input":{"apiKey":"sk-test"}}'`,
		out,
	)
	require.NoError(t, err)

	assert.JSONEq(t, createdProviderBody, strings.TrimSpace(out.String()))
	assert.JSONEq(
		t,
		`{"name":"anthropic-prod","providerName":"anthropic","input":{"apiKey":"sk-test"}}`,
		string(gotBody),
	)
	r.Verify(t)
}

func TestCreateProvider_MissingJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`--json` is required")
}

func TestCreateProvider_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(cmd, "--json not-json", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing provider JSON")
}

func TestCreateProvider_MissingName(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(
		cmd,
		`--json '{"providerName":"anthropic","input":{"apiKey":"sk-test"}}'`,
		out,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires `name` and a valid `providerName`")
}

func TestCreateProvider_MissingInput(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	_, err := test.Execute(
		cmd,
		`--json '{"name":"anthropic-prod","providerName":"anthropic"}'`,
		out,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires `input`")
}
