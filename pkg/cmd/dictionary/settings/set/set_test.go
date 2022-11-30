package set

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

func TestSetCmd(t *testing.T) {
	tests := []struct {
		name        string
		tty         bool
		cli         string
		wantsErr    bool
		wantsErrMsg string
		wantsOpts   SetOptions
	}{
		{
			name:        "no flags",
			tty:         true,
			cli:         "",
			wantsErr:    true,
			wantsErrMsg: "Either --disable-standard-entries and/or --enable-standard-entries or --reset-standard-entries must be set",
		},
		{
			name:        "enable, disable and clear",
			tty:         true,
			cli:         "--disable-standard-entries en --enable-standard-entries fr --reset-standard-entries",
			wantsErr:    true,
			wantsErrMsg: "You cannot reset standard entries and disable or enable standard entries at the same time",
		},
		{
			name:        "same language for disable and enable entries",
			cli:         "--disable-standard-entries en --enable-standard-entries en",
			tty:         false,
			wantsErr:    true,
			wantsErrMsg: "You cannot disable and enable standard entries for the same language: en",
		},
		{
			name: "disable standard entries",
			cli:  "--disable-standard-entries en",
			tty:  false,
			wantsOpts: SetOptions{
				DisableStandardEntries: []string{"en"},
				EnableStandardEntries:  []string{},
			},
		},
		{
			name: "enable standard entries",
			cli:  "--enable-standard-entries en",
			tty:  false,
			wantsOpts: SetOptions{
				EnableStandardEntries:  []string{"en"},
				DisableStandardEntries: []string{},
			},
		},
		{
			name: "reset standard entries",
			cli:  "--reset-standard-entries",
			tty:  false,
			wantsOpts: SetOptions{
				ResetStandardEntries:   true,
				EnableStandardEntries:  []string{},
				DisableStandardEntries: []string{},
			},
		},
		{
			name: "disable and enable standard entries",
			cli:  "--disable-standard-entries en --enable-standard-entries fr",
			tty:  false,
			wantsOpts: SetOptions{
				DisableStandardEntries: []string{"en"},
				EnableStandardEntries:  []string{"fr"},
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

			var opts *SetOptions
			cmd := NewSetCmd(f, func(o *SetOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantsErrMsg, err.Error())
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.DisableStandardEntries, opts.DisableStandardEntries)
			assert.Equal(t, tt.wantsOpts.EnableStandardEntries, opts.EnableStandardEntries)
			assert.Equal(t, tt.wantsOpts.ResetStandardEntries, opts.ResetStandardEntries)
		})
	}
}
