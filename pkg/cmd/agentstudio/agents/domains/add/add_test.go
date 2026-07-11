package add_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/add"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestAddDomain(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents/my-agent/allowed-domains"),
		func(req *http.Request) (*http.Response, error) {
			captured, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(
				`{"id":"domain_1","appId":"app1","agentId":"my-agent","domain":"https://a.example.com","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`,
			)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := add.NewAddCmd(f)
	_, err := test.Execute(cmd, "my-agent https://a.example.com", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"id":"domain_1"`)
	r.Verify(t)

	assert.JSONEq(t, `{"domain":"https://a.example.com"}`, string(captured))
}

func TestAddDomain_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := add.NewAddCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> and a <domain> argument")
}
