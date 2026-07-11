package update_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/update"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestUpdateProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("PATCH", "agent-studio/1/providers/my-provider"),
		httpmock.StringResponse(
			`{"id":"my-provider","name":"Renamed","providerName":"openai","input":{"apiKey":"sk-test"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	out.InBuf.WriteString(`{"name":"Renamed"}`)
	_, err := test.Execute(cmd, "my-provider --file -", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"name":"Renamed"`)
	r.Verify(t)
}

func TestUpdateProvider_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <provider-id> argument")
}
