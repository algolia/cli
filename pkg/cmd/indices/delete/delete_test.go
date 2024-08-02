package delete

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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
			cli:         "foo --confirm --includeReplicas",
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
				// First settings call with `Exists()`
				r.Register(
					httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/settings", index)),
					httpmock.JSONResponse(search.SettingsResponse{}),
				)
				if tt.hasReplicas {
					// Settings calls for the primary index and its replicas
					r.Register(
						httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/settings", index)),
						httpmock.JSONResponse(search.SettingsResponse{
							Replicas: []string{"bar"},
						}),
					)
					r.Register(
						httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/settings", index)),
						httpmock.JSONResponse(search.SettingsResponse{
							Replicas: []string{"bar"},
						}),
					)
					r.Register(
						httpmock.REST("GET", "1/indexes/bar/settings"),
						httpmock.JSONResponse(search.SettingsResponse{
							Primary: test.Pointer("foo"),
						}),
					)
					// Additional DELETE calls for the replicas
					r.Register(
						httpmock.REST("DELETE", "1/indexes/bar"),
						httpmock.JSONResponse(search.DeletedAtResponse{}),
					)
				}
				if tt.isReplica {
					// We want the first `Delete()` call to fail
					r.Register(
						httpmock.REST("DELETE", fmt.Sprintf("1/indexes/%s", index)),
						httpmock.ErrorResponse(),
					)
					// Second settings call to fetch the primary index name
					r.Register(
						httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/settings", index)),
						httpmock.JSONResponse(search.SettingsResponse{
							Primary: test.Pointer("bar"),
						}),
					)
					// Third settings call to fetch the primary index settings
					r.Register(
						httpmock.REST("GET", "1/indexes/bar/settings"),
						httpmock.JSONResponse(search.SettingsResponse{
							Replicas: []string{index},
						}),
					)
					// Fourth settings call to update the primary settings
					r.Register(
						httpmock.REST("PUT", "1/indexes/bar/settings"),
						httpmock.JSONResponse(search.UpdatedAtResponse{}),
					)
				}
				// Final `Delete()` call
				r.Register(
					httpmock.REST("DELETE", fmt.Sprintf("1/indexes/%s", index)),
					httpmock.JSONResponse(search.DeletedAtResponse{}),
				)
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
