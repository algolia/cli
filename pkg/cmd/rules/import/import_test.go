package importRules

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
	file := filepath.Join(t.TempDir(), "rules.ndjson")
	_ = ioutil.WriteFile(file, []byte("{\"objectID\":\"test\"}"), 0600)

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
				Indice:             "index",
				ForwardToReplicas:  true,
				ClearExistingRules: false,
			},
		},
		{
			name: "forward to replicas",
			cli:  fmt.Sprintf("index -F %s -f=false", file),
			wantsOpts: ImportOptions{
				Indice:             "index",
				ForwardToReplicas:  false,
				ClearExistingRules: false,
			},
		},
		{
			name:     "replace existing rules, no --confirm, no TTY",
			tty:      false,
			cli:      fmt.Sprintf("index -F %s -c", file),
			wantsErr: true,
		},
		{
			name: "clear existing rules, --confirm, no TTY",
			tty:  false,
			cli:  fmt.Sprintf("index -F %s -c --confirm", file),
			wantsOpts: ImportOptions{
				Indice:             "index",
				ForwardToReplicas:  true,
				ClearExistingRules: true,
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
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantsOpts.Indice, opts.Indice)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
			assert.Equal(t, tt.wantsOpts.ClearExistingRules, opts.ClearExistingRules)
		})
	}
}

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
