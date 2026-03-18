package run

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewRunCmd(t *testing.T) {
	tests := []struct {
		name      string
		cli       string
		wantsErr  bool
		wantsOpts RunOptions
	}{
		{
			name:     "dry run",
			cli:      "my-crawler --dry-run",
			wantsErr: false,
			wantsOpts: RunOptions{
				ID:     "my-crawler",
				DryRun: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()

			f := &cmdutil.Factory{IOStreams: io}

			var opts *RunOptions
			cmd := NewRunCmd(f, func(o *RunOptions) error {
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
			}

			require.NoError(t, err)
			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())
			assert.Equal(t, tt.wantsOpts.ID, opts.ID)
			assert.Equal(t, tt.wantsOpts.DryRun, opts.DryRun)
		})
	}
}

func Test_runRunCmd_dryRunJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewRunCmd(f, nil)

	out, err := test.Execute(cmd, "my-crawler --dry-run --output json", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `"action":"run_crawler"`)
	assert.Contains(t, out.String(), `"id":"my-crawler"`)
	assert.Contains(t, out.String(), `"dryRun":true`)
}
