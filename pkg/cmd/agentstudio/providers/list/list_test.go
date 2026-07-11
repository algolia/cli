package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestListProviders(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/providers"),
		httpmock.StringResponse(
			`{"data":[{"id":"provider_1","name":"My Provider","providerName":"openai","input":{"apiKey":"sk-***"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}],"pagination":{"page":0,"limit":20,"total":1}}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"provider_1"`)
	r.Verify(t)
}
