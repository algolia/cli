package delete

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

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
				APIKeys:   []string{"foo"},
			},
		},
		{
			name:     "--confirm without tty",
			cli:      "foo --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				APIKeys:   []string{"foo"},
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

			assert.Equal(t, tt.wantsOpts.APIKeys, opts.APIKeys)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
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

	cmd := NewDeleteCmd(factory, nil)

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

func TestDeleteAPIKeys(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_search.NewMockClientInterface(mockCtrl)

	client.EXPECT().GetAPIKey("test").Return(search.Key{}, nil).Times(1)
	client.EXPECT().DeleteAPIKey("test").Return(search.DeleteKeyRes{}, nil).Times(1)

	output, err := runCommand(client, false, "test --confirm")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ``, output.String())
	assert.Equal(t, ``, output.Stderr())
}

func TestDeleteAPIKeys_multiple(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_search.NewMockClientInterface(mockCtrl)

	gomock.InOrder(
		client.EXPECT().GetAPIKey("test").Return(search.Key{}, nil).Times(1),
		client.EXPECT().GetAPIKey("test1").Return(search.Key{}, nil).Times(1),
		client.EXPECT().DeleteAPIKey("test").Return(search.DeleteKeyRes{}, nil).Times(1),
		client.EXPECT().DeleteAPIKey("test1").Return(search.DeleteKeyRes{}, nil).Times(1),
	)

	output, err := runCommand(client, false, "test test1 --confirm")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ``, output.String())
	assert.Equal(t, ``, output.Stderr())
}

func TestDeleteAPIKeys_not_exists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_search.NewMockClientInterface(mockCtrl)

	client.EXPECT().GetAPIKey("test").Return(search.Key{}, errors.New("")).Times(1)

	output, err := runCommand(client, false, "test --confirm")
	if err.Error() != fmt.Errorf("API key \"test\" does not exist").Error() {
		t.Fatal(err)
	}

	assert.Equal(t, ``, output.String())
	assert.Equal(t, ``, output.Stderr())
}
