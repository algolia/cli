package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	providerID = "22222222-2222-4222-8222-222222222222"
	getPath    = "agent-studio/1/providers/" + providerID
	getBody    = `{
  "id": "22222222-2222-4222-8222-222222222222",
  "name": "openai",
  "providerName": "openai",
  "input": {"apiKey": "sk-test"},
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z"
}`
)

func TestGetProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", getPath), httpmock.StringResponse(getBody))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, providerID, out)
	require.NoError(t, err)

	assert.JSONEq(t, getBody, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestGetProvider_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agentstudio providers get requires a <provider-id> argument")
}
