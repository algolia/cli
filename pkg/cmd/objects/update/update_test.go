package update

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runUpdateCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "objects.json")
	err := os.WriteFile(tmpFile, []byte(`{"objectID":"foo"}`), 0600)
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
			wantOut: "✓ Successfully updated 1 objects on foo in",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully updated 1 objects on foo in",
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "foo -F -",
			stdin:   `{"objectID": "foo"},`,
			wantErr: "X Found 1 error (out of 1 objects) while parsing the file:\n  line 1: invalid character ',' after top-level value\n",
		},
		{
			name: "from stdin with invalid JSON (multiple objects)",
			cli:  "foo -F -",
			stdin: `{"objectID": "foo"},
			{"test": "bar"}`,
			wantErr: "X Found 2 errors (out of 2 objects) while parsing the file:\n  line 1: invalid character ',' after top-level value\n  line 2: objectID is required\n",
		},
		{
			name:    "from stdin with invalid JSON (1 object) with --continue-on-error",
			cli:     "foo -F - --continue-on-error",
			stdin:   `{"objectID": "foo"},`,
			wantErr: "X Found 1 error (out of 1 objects) while parsing the file:\n  line 1: invalid character ',' after top-level value\n",
		},
		{
			name: "from stdin with invalid JSON (2 objects) with --continue-on-error",
			cli:  "foo -F - --continue-on-error",
			stdin: `{"objectID": "foo"}
			{"test": "bar"}`,
			wantOut: "✓ Successfully updated 1 objects on foo in",
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
				r.Register(httpmock.REST("POST", "1/indexes/foo/batch"), httpmock.JSONResponse(search.BatchRes{}))
			}
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewUpdateCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.Contains(t, out.String(), tt.wantOut)
		})
	}
}
