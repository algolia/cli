package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/secretkeys/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestListSecretKeys(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/secret-keys"),
		httpmock.StringResponse(
			`{"data":[{"id":"key_1","name":"My Key","value":"secret_***","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"agentIds":[]}],"pagination":{"page":0,"limit":20,"total":1}}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"key_1"`)
	r.Verify(t)
}
