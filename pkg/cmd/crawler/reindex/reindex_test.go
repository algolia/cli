package reindex

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewReindexCmd(t *testing.T) {
	tests := []struct {
		name      string
		cli       string
		wantsErr  bool
		wantsOpts ReindexOptions
	}{
		{
			name:     "dry run",
			cli:      "crawler-1 crawler-2 --dry-run",
			wantsErr: false,
			wantsOpts: ReindexOptions{
				IDs:    []string{"crawler-1", "crawler-2"},
				DryRun: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()

			f := &cmdutil.Factory{IOStreams: io}

			var opts *ReindexOptions
			cmd := NewReindexCmd(f, func(o *ReindexOptions) error {
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
			assert.Equal(t, tt.wantsOpts.IDs, opts.IDs)
			assert.Equal(t, tt.wantsOpts.DryRun, opts.DryRun)
		})
	}
}

func Test_runReindexCmd_dryRunJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewReindexCmd(f, nil)

	out, err := test.Execute(cmd, "crawler-1 crawler-2 --dry-run --output json", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `"action":"reindex_crawlers"`)
	assert.Contains(t, out.String(), `"ids":["crawler-1","crawler-2"]`)
	assert.Contains(t, out.String(), `"dryRun":true`)
}
