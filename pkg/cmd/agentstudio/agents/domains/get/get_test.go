package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestGetDomain(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/agents/my-agent/allowed-domains/domain_1"),
		httpmock.StringResponse(
			`{"id":"domain_1","appId":"app1","agentId":"my-agent","domain":"https://a.example.com","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
		),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-agent domain_1", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"domain":"https://a.example.com"`)
	r.Verify(t)
}

func TestGetDomain_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> and a <domain-id> argument")
}
