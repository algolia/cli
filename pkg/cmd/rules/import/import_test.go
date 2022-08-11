package importRules

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
	tmpFile := filepath.Join(t.TempDir(), "rules.json")
	err := ioutil.WriteFile(tmpFile, []byte("{\"objectID\":\"test\"}"), 0600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cli     string
		stdin   string
		wantOut string
		wantErr string
	}{
		{
			name:    "from stdin",
			cli:     "foo -F -",
			stdin:   `{"objectID":"test"}`,
			wantOut: "✓ Successfully imported 1 rules to foo\n",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully imported 1 rules to foo\n",
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "foo -F -",
			stdin:   `{"objectID", "test"},`,
			wantErr: "failed to parse JSON rule on line 0: invalid character ',' after object key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.wantErr == "" {
				r.Register(httpmock.REST("POST", "1/indexes/foo/rules/batch"), httpmock.JSONResponse(search.UpdateTaskRes{}))
			}
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewImportCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
