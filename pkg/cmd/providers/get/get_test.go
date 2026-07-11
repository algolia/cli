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

func TestGetProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/providers/my-provider"),
		httpmock.StringResponse(
			`{"id":"my-provider","name":"My Provider","providerName":"openai","input":{"apiKey":"sk-***"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-provider", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"my-provider"`)
	r.Verify(t)
}

func TestGetProvider_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <provider-id> argument")
}
