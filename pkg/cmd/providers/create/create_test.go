package create_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/create"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const validProviderBody = `{"name":"My Provider","providerName":"openai","input":{"apiKey":"sk-test"}}`

func TestCreateProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "agent-studio/1/providers"),
		httpmock.StringResponse(
			`{"id":"provider_1","name":"My Provider","providerName":"openai","input":{"apiKey":"sk-test"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	out.InBuf.WriteString(validProviderBody)
	_, err := test.Execute(cmd, "--file -", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"provider_1"`)
	r.Verify(t)
}

func TestCreateProvider_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := create.NewCreateCmd(f)
	out.InBuf.WriteString("not-json")
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing provider JSON")
}
