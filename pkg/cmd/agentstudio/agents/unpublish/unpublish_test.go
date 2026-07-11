package unpublish_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/unpublish"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestUnpublishAgent(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents/my-agent/unpublish"),
		httpmock.StringResponse(
			`{"id":"my-agent","name":"My Agent","description":null,"status":"draft","providerId":null,"instructions":"be nice","config":{},"createdAt":"2026-01-01T00:00:00Z","updatedAt":null,"lastUsedAt":null}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := unpublish.NewUnpublishCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"status":"draft"`)
	r.Verify(t)
}

func TestUnpublishAgent_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := unpublish.NewUnpublishCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> argument")
}
