package importSynonyms

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewImportCmd(t *testing.T) {
	file := filepath.Join(t.TempDir(), "synonyms.ndjson")
	_ = ioutil.WriteFile(file, []byte("{\"objectID\":\"test\", \"type\": \"synonym\", \"synonyms\": [\"test\"]}"), 0600)

	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts ImportOptions
	}{
		{
			name:     "no file specified",
			cli:      "index",
			wantsErr: true,
		},
		{
			name:     "file not found",
			cli:      "index --file not-found",
			wantsErr: true,
		},
		{
			name: "file specified",
			cli:  fmt.Sprintf("index -F %s", file),
			wantsOpts: ImportOptions{
				Index:                   "index",
				ForwardToReplicas:       true,
				ReplaceExistingSynonyms: false,
			},
		},
		{
			name: "forward to replicas",
			cli:  fmt.Sprintf("index -F %s -f=false", file),
			wantsOpts: ImportOptions{
				Index:                   "index",
				ForwardToReplicas:       false,
				ReplaceExistingSynonyms: false,
			},
		},
		{
			name: "replace existing synonyms",
			cli:  fmt.Sprintf("index -F %s -r=true", file),
			wantsOpts: ImportOptions{
				Index:                   "index",
				ForwardToReplicas:       true,
				ReplaceExistingSynonyms: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			if tt.tty {
				io.SetStdinTTY(tt.tty)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *ImportOptions
			cmd := NewImportCmd(f, func(o *ImportOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.wantsOpts.Index, opts.Index)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
			assert.Equal(t, tt.wantsOpts.ReplaceExistingSynonyms, opts.ReplaceExistingSynonyms)
		})
	}
}

func Test_runExportCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "synonyms.json")
	err := ioutil.WriteFile(tmpFile, []byte("{\"objectID\":\"test\", \"type\": \"synonym\", \"synonyms\": [\"test\"]}"), 0600)
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
			stdin:   `{"objectID":"test", "type": "synonym", "synonyms": ["test"]}`,
			wantOut: "✓ Successfully imported 1 synonyms to foo\n",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully imported 1 synonyms to foo\n",
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "foo -F -",
			stdin:   `{"objectID", "test"},`,
			wantErr: "failed to parse JSON synonym on line 0: invalid character ',' after object key",
		},
		{
			name:    "from file with forward to replicas",
			cli:     fmt.Sprintf("foo -F '%s' -f", tmpFile),
			wantOut: "✓ Successfully imported 1 synonyms to foo\n",
		},
		{
			name:    "from file with replace existing synonyms",
			cli:     fmt.Sprintf("foo -F '%s' -r", tmpFile),
			wantOut: "✓ Successfully imported 1 synonyms to foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.wantErr == "" {
				r.Register(httpmock.REST("POST", "1/indexes/foo/synonyms/batch"), httpmock.JSONResponse(search.UpdateTaskRes{}))
			}
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewImportCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
