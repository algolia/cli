package bulkdelete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/domains/bulkdelete"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestBulkDeleteDomains(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("DELETE", "agent-studio/1/agents/my-agent/allowed-domains/bulk"),
		httpmock.StringResponse(``),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := bulkdelete.NewBulkDeleteCmd(f)
	_, err := test.Execute(cmd, "my-agent domain_1 domain_2 --confirm", out)
	require.NoError(t, err)
	r.Verify(t)
}

func TestBulkDeleteDomains_RequiresConfirmation(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := bulkdelete.NewBulkDeleteCmd(f)
	_, err := test.Execute(cmd, "my-agent domain_1 domain_2", out)
	require.Error(t, err)
	assert.Empty(t, r.Requests)
}
