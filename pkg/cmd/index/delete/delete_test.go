package delete

import (
	"fmt"
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

func TestNewDeleteCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts DeleteOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "foo",
			tty:      false,
			wantsErr: true,
			wantsOpts: DeleteOptions{
				DoConfirm: true,
				Index:     "foo",
			},
		},
		{
			name:     "--confirm without tty",
			cli:      "foo --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			if tt.tty {
				io.SetStdinTTY(tt.tty)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *DeleteOptions
			cmd := NewDeleteCmd(f, func(o *DeleteOptions) error {
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

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.Index, opts.Index)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		index   string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "no TTY",
			cli:     "foo --confirm",
			index:   "foo",
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY",
			cli:     "foo --confirm",
			index:   "foo",
			isTTY:   true,
			wantOut: "✓ Deleted index foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("DELETE", fmt.Sprintf("1/indexes/%s", tt.index)), httpmock.JSONResponse(search.DeleteKeyRes{}))
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewDeleteCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
