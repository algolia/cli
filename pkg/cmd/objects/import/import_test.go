package importRecords

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runExportCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "objects.json")
	err := ioutil.WriteFile(tmpFile, []byte("{\"objectID\":\"foo\"}"), 0600)
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
			stdin:   `{"objectID": "foo"}`,
			wantOut: "✓ Successfully imported 1 objects to foo in",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully imported 1 objects to foo in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("POST", "1/indexes/foo/batch"), httpmock.JSONResponse(search.BatchRes{}))
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewImportCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Contains(t, out.String(), tt.wantOut)
		})
	}
}
