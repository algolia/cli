package clear

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func TestNewClearCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts ClearOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "foo",
			tty:      false,
			wantsErr: true,
			wantsOpts: ClearOptions{
				DoConfirm: true,
				Index:     "foo",
			},
		},
		{
			name:     "--confirm without tty",
			cli:      "foo --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: ClearOptions{
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

			var opts *ClearOptions
			cmd := NewClearCmd(f, func(o *ClearOptions) error {
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

func runCommand(http *httpmock.Registry, isTTY bool, cli string) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.Test()
	io.SetStdoutTTY(isTTY)
	io.SetStdinTTY(isTTY)
	io.SetStderrTTY(isTTY)

	client := search.NewClientWithConfig(search.Configuration{
		Requester: http,
	})

	factory := &cmdutil.Factory{
		IOStreams: io,
		SearchClient: func() (*search.Client, error) {
			return client, nil
		},
	}

	cmd := NewClearCmd(factory, nil)

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

func Test_runCreateCmd(t *testing.T) {
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
			wantOut: "✓ Cleared index foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("POST", fmt.Sprintf("1/indexes/%s/clear", tt.index)), httpmock.JSONResponse(search.CreateKeyRes{Key: "foo"}))
			defer r.Verify(t)

			out, err := runCommand(&r, tt.isTTY, tt.cli)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
