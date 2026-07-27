package create

import (
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

const unroutableAPIURL = "http://127.0.0.1:1"

func withoutSession(t *testing.T) {
	t.Helper()
	t.Setenv("ALGOLIA_API_KEY", "")
	t.Setenv("ALGOLIA_ADMIN_API_KEY", "")
	t.Setenv("ALGOLIA_APPLICATION_ID", "")
	t.Setenv("ALGOLIA_API_URL", unroutableAPIURL)
	keyring.MockInit()
}

func explicitKeyConfig() *test.ConfigStub {
	return &test.ConfigStub{
		CurrentProfile: config.Profile{ApplicationID: "APP1", APIKey: "adm"},
	}
}

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
			cli:     "--acl search",
			isTTY:   false,
			wantOut: "foo\n",
		},
		{
			name:    "TTY",
			cli:     "--acl search",
			isTTY:   true,
			wantOut: "✓ API key created: foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withoutSession(t)

			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/keys"),
				httpmock.JSONResponse(search.AddApiKeyResponse{Key: "foo"}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, explicitKeyConfig(), "")
			cmd := NewCreateCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
