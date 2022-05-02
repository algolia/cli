package add

import (
	"testing"

	"github.com/google/shlex"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAddCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts AddOptions
	}{
		{
			name:     "not interactive, missing flags",
			cli:      "",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "not interactive, all flags",
			cli:      "--name my-app --app-id my-app-id --admin-api-key my-admin-api-key",
			tty:      false,
			wantsErr: false,
			wantsOpts: AddOptions{
				Application: config.Application{
					Name:        "my-app",
					ID:          "my-app-id",
					AdminAPIKey: "my-admin-api-key",
				},
			},
		},
		{
			name:     "interactive, no flags",
			cli:      "",
			tty:      true,
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			io.SetStdinTTY(tt.tty)
			io.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *AddOptions
			cmd := NewAddCmd(f, func(o *AddOptions) error {
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

			assert.Equal(t, tt.wantsOpts.Application.Name, opts.Application.Name)
			assert.Equal(t, tt.wantsOpts.Application.ID, opts.Application.ID)
			assert.Equal(t, tt.wantsOpts.Application.AdminAPIKey, opts.Application.AdminAPIKey)
			assert.Equal(t, tt.wantsOpts.Application.Default, opts.Application.Default)
		})
	}
}
