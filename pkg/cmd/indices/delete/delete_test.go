package delete

import (
	"fmt"
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

func TestNewDeleteCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts DeleteOptions
	}{
		{
			name:     "single index, no --confirm, without tty",
			cli:      "foo",
			tty:      false,
			wantsErr: true,
			wantsOpts: DeleteOptions{
				DoConfirm: true,
				Indices:   []string{"foo"},
			},
		},
		{
			name:     "single index, --confirm, without tty",
			cli:      "foo --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Indices:   []string{"foo"},
			},
		},
		{
			name:     "multiple indices, --confirm, without tty",
			cli:      "foo bar baz --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Indices:   []string{"foo", "bar", "baz"},
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

			assert.Equal(t, tt.wantsOpts.Indices, opts.Indices)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
	tests := []struct {
		name        string
		cli         string
		indices     []string
		isReplica   bool
		hasReplicas bool
		isTTY       bool
		wantOut     string
	}{
		{
			name:    "no TTY",
			cli:     "foo --confirm",
			indices: []string{"foo"},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY",
			cli:     "foo --confirm",
			indices: []string{"foo"},
			isTTY:   true,
			wantOut: "✓ Deleted index foo\n",
		},
		{
			name:    "no TTY, multiple indices",
			cli:     "foo bar --confirm",
			indices: []string{"foo", "bar"},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY, multiple indices",
			cli:     "foo bar --confirm",
			indices: []string{"foo", "bar"},
			isTTY:   true,
			wantOut: "✓ Deleted indices foo, bar\n",
		},
		{
			name:      "TTY, replica indices",
			cli:       "foo --confirm",
			indices:   []string{"foo"},
			isReplica: true,
			isTTY:     true,
			wantOut:   "✓ Deleted index foo\n",
		},
		{
			name:        "TTY, has replica indices",
			cli:         "foo --confirm --include-replicas",
			indices:     []string{"foo"},
			hasReplicas: true,
			isTTY:       true,
			wantOut:     "✓ Deleted index foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			for _, index := range tt.indices {
				// GetSettings request to check if index exists
				r.Register(
					httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/settings", index)),
					httpmock.JSONResponse(search.SettingsResponse{}),
				)
				// DeleteIndex request for the primary index
				r.Register(
					httpmock.REST("DELETE", fmt.Sprintf("1/indexes/%s", index)),
					httpmock.JSONResponse(search.DeletedAtResponse{}),
				)
				// if tt.hasReplicas {
				// 	r.Register(
				// 		httpmock.REST("DELETE", "1/indexes/bar"),
				// 		httpmock.JSONResponse(search.DeletedAtResponse{}),
				// 	)
				// }
			}
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
