package list

import (
	"testing"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogsCmd(t *testing.T) {
	testIndexName := "foo"
	tests := []struct {
		name      string
		cli       string
		wantsErr  bool
		wantsOpts LogOptions
	}{
		{
			name:     "with default options",
			cli:      "",
			wantsErr: false,
			wantsOpts: LogOptions{
				Entries:   5,
				Start:     1,
				LogType:   "all",
				IndexName: nil,
			},
		},
		{
			name:     "with 69 entries, starting at 420, type query, filtered by index foo",
			cli:      "--entries 69 --start 420 --type query --index foo",
			wantsErr: false,
			wantsOpts: LogOptions{
				Entries:   69,
				Start:     420,
				LogType:   "query",
				IndexName: &testIndexName,
			},
		},
	}

	for _, tt := range tests {
		f := &cmdutil.Factory{}
		var opts *LogOptions
		cmd := NewListCmd(f, func(o *LogOptions) error {
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
		assert.Equal(t, tt.wantsOpts.Entries, opts.Entries)
		assert.Equal(t, tt.wantsOpts.Start, opts.Start)
		assert.Equal(t, tt.wantsOpts.LogType, opts.LogType)
		assert.Equal(t, tt.wantsOpts.IndexName, opts.IndexName)
	}
}
