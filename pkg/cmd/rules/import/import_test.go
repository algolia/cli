package importrules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock/v4"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test/v4"
)

func TestNewImportCmd(t *testing.T) {
	file := filepath.Join(t.TempDir(), "rules.ndjson")
	_ = os.WriteFile(file, []byte("{\"objectID\":\"test\"}"), 0o600)

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
				Index:              "index",
				ForwardToReplicas:  true,
				ClearExistingRules: false,
			},
		},
		{
			name: "forward to replicas",
			cli:  fmt.Sprintf("index -F %s -f=false", file),
			wantsOpts: ImportOptions{
				Index:              "index",
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
				Index:              "index",
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
			}
			require.NoError(t, err)

			assert.Equal(t, tt.wantsOpts.Index, opts.Index)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
			assert.Equal(t, tt.wantsOpts.ClearExistingRules, opts.ClearExistingRules)
		})
	}
}

func Test_runExportCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "rules.json")
	err := os.WriteFile(tmpFile, []byte("{\"objectID\":\"test\"}"), 0o600)
	require.NoError(t, err)

	var largeBatchBuilder strings.Builder
	for i := 0; i < 1001; i += 1 {
		largeBatchBuilder.Write([]byte("{\"objectID\":\"test\"}\n"))
	}

	tests := []struct {
		name    string
		cli     string
		stdin   string
		wantOut string
		wantErr string
		setup   func(*httpmock.Registry)
	}{
		{
			name:    "from stdin",
			cli:     "foo -F -",
			stdin:   `{"objectID":"test"}`,
			wantOut: "✓ Successfully imported 1 rules to foo\n",
			setup: func(r *httpmock.Registry) {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/rules/batch"),
					httpmock.JSONResponse(search.UpdatedAtResponse{}),
				)
			},
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			wantOut: "✓ Successfully imported 1 rules to foo\n",
			setup: func(r *httpmock.Registry) {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/rules/batch"),
					httpmock.JSONResponse(search.UpdatedAtResponse{}),
				)
			},
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "foo -F -",
			stdin:   `{"objectID", "test"},`,
			wantErr: "failed to parse JSON rule on line 0: invalid character ',' after object key",
			setup:   func(r *httpmock.Registry) {},
		},
		{
			name:    "from empty batch with clear existing",
			cli:     "foo -c -y -F -",
			stdin:   ``,
			wantOut: "✓ Successfully imported 0 rules to foo\n",
			setup: func(r *httpmock.Registry) {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/rules/clear"),
					httpmock.JSONResponse(search.UpdatedAtResponse{}),
				)
			},
		},
		{
			name:    "from empty batch without clear existing",
			cli:     "foo -F -",
			stdin:   ``,
			wantOut: "✓ Successfully imported 0 rules to foo\n",
			setup:   func(r *httpmock.Registry) {},
		},
		{
			name:    "from large batch clear existing",
			cli:     "foo -c -y -F -",
			stdin:   largeBatchBuilder.String(),
			wantOut: "✓ Successfully imported 1001 rules to foo\n",
			setup: func(r *httpmock.Registry) {
				r.Register(httpmock.Matcher(func(req *http.Request) bool {
					return httpmock.REST("POST", "1/indexes/foo/rules/batch")(req) &&
						req.URL.Query().Get("clearExistingRules") == "true"
				}), httpmock.JSONResponse(search.UpdatedAtResponse{}))
				r.Register(httpmock.Matcher(func(req *http.Request) bool {
					return httpmock.REST("POST", "1/indexes/foo/rules/batch")(req) &&
						req.URL.Query().Get("clearExistingRules") == ""
				}), httpmock.JSONResponse(search.UpdatedAtResponse{}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.setup != nil {
				tt.setup(&r)
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
