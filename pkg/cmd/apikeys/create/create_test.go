package create

import (
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewCreateCmd(t *testing.T) {
	oneHour, _ := time.ParseDuration("1h")

	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts CreateOptions
	}{
		{
			name:     "all the flags",
			cli:      "-i foo,bar --acl search,browse -r \"http://foo.com\" -u 1h -d \"description\"",
			tty:      false,
			wantsErr: false,
			wantsOpts: CreateOptions{
				ACL:         []string{"search", "browse"},
				Indices:     []string{"foo", "bar"},
				Description: "description",
				Referers:    []string{"http://foo.com"},
				Validity:    oneHour,
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

			var opts *CreateOptions
			cmd := NewCreateCmd(f, func(o *CreateOptions) error {
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

			assert.Equal(t, tt.wantsOpts.ACL, opts.ACL)
			assert.Equal(t, tt.wantsOpts.Indices, opts.Indices)
			assert.Equal(t, tt.wantsOpts.Description, opts.Description)
			assert.Equal(t, tt.wantsOpts.Referers, opts.Referers)
			assert.Equal(t, tt.wantsOpts.Validity, opts.Validity)
		})
	}
}

func Test_runCreateCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "no TTY",
			cli:     "",
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY",
			cli:     "",
			isTTY:   true,
			wantOut: "✓ API key created: foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/keys"),
				httpmock.JSONResponse(search.AddApiKeyResponse{Key: "foo"}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewCreateCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
