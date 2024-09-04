package copy

import (
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

func TestNewCopyCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts CopyOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "foo bar",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "--confirm without tty",
			cli:      "foo bar --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: CopyOptions{
				DoConfirm:        false,
				SourceIndex:      "foo",
				DestinationIndex: "bar",
				Scope:            []string{},
				Wait:             true,
			},
		},
		{
			name:     "specifying scopes",
			cli:      "foo bar --scope settings",
			tty:      true,
			wantsErr: false,
			wantsOpts: CopyOptions{
				DoConfirm:        true,
				SourceIndex:      "foo",
				DestinationIndex: "bar",
				Scope:            []string{"settings"},
				Wait:             true,
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

			var opts *CopyOptions
			cmd := NewCopyCmd(f, func(o *CopyOptions) error {
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

			assert.Equal(t, tt.wantsOpts.SourceIndex, opts.SourceIndex)
			assert.Equal(t, tt.wantsOpts.DestinationIndex, opts.DestinationIndex)
			assert.Equal(t, tt.wantsOpts.Scope, opts.Scope)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runCreateCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		source  string
		dest    string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "no TTY",
			cli:     "foo bar --confirm",
			source:  "foo",
			dest:    "bar",
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY",
			cli:     "foo bar --confirm",
			source:  "foo",
			dest:    "bar",
			isTTY:   true,
			wantOut: "✓ Copied records, settings, synonyms, and rules from foo to bar\n",
		},
		{
			name:    "Specifying scopes",
			cli:     "foo bar --confirm --scope settings",
			source:  "foo",
			dest:    "bar",
			isTTY:   true,
			wantOut: "✓ Copied settings from foo to bar\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/indexes/foo/operation"),
				httpmock.JSONResponse(search.UpdatedAtResponse{}),
			)
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewCopyCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
