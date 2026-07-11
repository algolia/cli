package update_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/update"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestUpdateAgent(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("PATCH", "agent-studio/1/agents/my-agent"),
		httpmock.StringResponse(
			`{"id":"my-agent","name":"Renamed","description":null,"status":"draft","providerId":null,"instructions":"be nice","config":{},"createdAt":"2026-01-01T00:00:00Z","updatedAt":null,"lastUsedAt":null}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	out.InBuf.WriteString(`{"name":"Renamed"}`)
	_, err := test.Execute(cmd, "my-agent --file -", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"name":"Renamed"`)
	r.Verify(t)
}

func TestUpdateAgent_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := update.NewUpdateCmd(f)
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> argument")
}
