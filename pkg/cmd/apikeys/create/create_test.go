package create

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/golang/mock/gomock"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/mock_search"
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

func runCommand(client search.ClientInterface, isTTY bool, cli string) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.Test()
	io.SetStdoutTTY(isTTY)
	io.SetStdinTTY(isTTY)
	io.SetStderrTTY(isTTY)

	factory := &cmdutil.Factory{
		IOStreams: io,
		SearchClient: func() (search.ClientInterface, error) {
			return client, nil
		},
	}

	cmd := NewCreateCmd(factory, nil)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	cmd.SetIn(&bytes.Buffer{})
	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	_, err = cmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func TestCreateAPIKeys(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_search.NewMockClientInterface(mockCtrl)

	client.EXPECT().AddAPIKey(search.Key{}).Return(search.CreateKeyRes{}, nil).Times(1)

	output, err := runCommand(client, false, "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ``, output.String())
	assert.Equal(t, ``, output.Stderr())
}

func TestCreateAPIKeys_tty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_search.NewMockClientInterface(mockCtrl)

	client.EXPECT().AddAPIKey(search.Key{}).Return(search.CreateKeyRes{Key: "foo"}, nil).Times(1)

	output, err := runCommand(client, true, "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "✓ API key created foo\n", output.String())
	assert.Equal(t, ``, output.Stderr())
}
