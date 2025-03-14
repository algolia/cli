package create

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewCreateCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(tmpFile, []byte("{\"enableReRanking\":false}"), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts CreateOptions
	}{
		{
			name:     "no tty",
			cli:      fmt.Sprintf("my-crawler -F '%s'", tmpFile),
			tty:      false,
			wantsErr: false,
			wantsOpts: CreateOptions{
				Name:   "my-crawler",
				config: crawler.Config{},
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

			assert.Equal(t, tt.wantsOpts.Name, opts.Name)
			assert.Equal(t, tt.wantsOpts.config, opts.config)
		})
	}
}

func Test_runCreateCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(tmpFile, []byte("{\"enableReRanking\":false}"), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cli     string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "no tty",
			cli:     fmt.Sprintf("my-crawler -F '%s'", tmpFile),
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "tty",
			cli:     fmt.Sprintf("my-crawler -F '%s'", tmpFile),
			isTTY:   true,
			wantOut: "✓ Crawler my-crawler created: crawler-id\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			res := crawler.Crawler{ID: "crawler-id"}
			r.Register(httpmock.REST("POST", "api/1/crawlers"), httpmock.JSONResponse(res))
			defer r.Verify(t)

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
