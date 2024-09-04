package set

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runExportCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	err := os.WriteFile(tmpFile, []byte("{\"enableReRanking\":false}"), 0600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cli     string
		stdin   string
		wantOut string
	}{
		{
			name:    "from stdin",
			cli:     "foo -F -",
			stdin:   `{"enableReRanking": true}`,
			wantOut: "✓ Imported settings on foo\n",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Imported settings on foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("PUT", "1/indexes/foo/settings"),
				httpmock.JSONResponse(search.UpdatedAtResponse{}),
			)
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewImportCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
