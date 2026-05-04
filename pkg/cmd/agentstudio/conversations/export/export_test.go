package export

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_NewExportCmd_writesToStdout(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodGet, "1/agents/a1/conversations/export"),
		httpmock.StringResponse(`{"conversations":[]}`),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewExportCmd(f, nil)
	_, err := test.Execute(cmd, "a1", out)
	require.NoError(t, err)
	assert.Equal(t, `{"conversations":[]}`, out.String())
}

func Test_NewExportCmd_writesToFile(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST(http.MethodGet, "1/agents/a1/conversations/export"),
		httpmock.StringResponse(`payload`),
	)
	defer r.Verify(t)

	dir := t.TempDir()
	dest := filepath.Join(dir, "out.json")

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewExportCmd(f, nil)
	_, err := test.Execute(cmd, "a1 -o "+dest, out)
	require.NoError(t, err)

	got, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "payload", string(got))
}
