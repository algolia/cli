package importRecords

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

func Test_runImportCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "objects.json")
	err := os.WriteFile(tmpFile, []byte("{\"objectID\":\"foo\"}"), 0600)
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
			stdin:   `{"objectID": "foo"}`,
			wantOut: "✓ Successfully imported 1 objects to foo in",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully imported 1 objects to foo in",
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "foo -F -",
			stdin:   `{"objectID", "foo"},`,
			wantErr: "failed to parse JSON object on line 0: invalid character ',' after object key",
		},
		{
			name:    "missing file flag",
			cli:     "foo",
			wantErr: "required flag(s) \"file\" not set",
		},
		{
			name:    "non-existant file",
			cli:     "foo -F /tmp/foo",
			wantErr: "open /tmp/foo: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.wantErr == "" {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/batch"),
					httpmock.JSONResponse(search.BatchResponse{}),
				)
			}
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewImportCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.Contains(t, out.String(), tt.wantOut)
		})
	}
}
