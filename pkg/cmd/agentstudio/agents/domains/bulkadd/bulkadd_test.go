package bulkadd_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/bulkadd"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestBulkAddDomains(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(
		httpmock.REST("POST", "agent-studio/1/agents/my-agent/allowed-domains/bulk"),
		func(req *http.Request) (*http.Response, error) {
			captured, _ = io.ReadAll(req.Body)
			return httpmock.StringResponse(`{"domains":[]}`)(req)
		},
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := bulkadd.NewBulkAddCmd(f)
	_, err := test.Execute(cmd, "my-agent https://a.example.com https://b.example.com", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"domains"`)
	r.Verify(t)

	assert.JSONEq(t, `{"domains":["https://a.example.com","https://b.example.com"]}`, string(captured))
}

func TestBulkAddDomains_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := bulkadd.NewBulkAddCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.Error(t, err)
	assert.NotEmpty(t, err.Error())
}
