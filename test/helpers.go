package test

import (
	"bytes"
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
)

type CmdInOut struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (c CmdInOut) String() string {
	return c.OutBuf.String()
}

func (c CmdInOut) Stderr() string {
	return c.ErrBuf.String()
}

type OutputStub struct {
	Out   []byte
	Error error
}

func (s OutputStub) Output() ([]byte, error) {
	if s.Error != nil {
		return s.Out, s.Error
	}
	return s.Out, nil
}

func (s OutputStub) Run() error {
	if s.Error != nil {
		return s.Error
	}
	return nil
}

func NewFactory(isTTY bool, http *httpmock.Registry, cfg config.IConfig, in string) (*cmdutil.Factory, *CmdInOut) {
	io, stdin, stdout, stderr := iostreams.Test()
	io.SetStdoutTTY(isTTY)
	io.SetStdinTTY(isTTY)
	io.SetStderrTTY(isTTY)

	if in != "" {
		stdin.WriteString(in)
	}

	f := &cmdutil.Factory{
		IOStreams: io,
	}

	if http != nil {
		f.SearchClient = func() (*search.Client, error) {
			return search.NewClientWithConfig(search.Configuration{
				Requester: http,
			}), nil
		}
	}

	if cfg != nil {
		f.Config = cfg
	} else {
		f.Config = &config.Config{}
	}

	return f, &CmdInOut{
		InBuf:  stdin,
		OutBuf: stdout,
		ErrBuf: stderr,
	}
}

func Execute(cmd *cobra.Command, cli string, inOut *CmdInOut) (*CmdInOut, error) {
	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	if inOut.InBuf != nil {
		cmd.SetIn(inOut.InBuf)
	} else {
		cmd.SetIn(&bytes.Buffer{})
	}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	_, err = cmd.ExecuteC()
	if err != nil {
		return nil, err
	}

	return inOut, nil
}
