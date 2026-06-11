package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ShouldMigrate(t *testing.T) {
	tests := []struct {
		name       string
		configFile bool
		stateFile  bool
		want       bool
	}{
		{
			name:       "config.toml only: migration pending",
			configFile: true,
			stateFile:  false,
			want:       true,
		},
		{
			name:       "both files: already migrated (or new model in use)",
			configFile: true,
			stateFile:  true,
			want:       false,
		},
		{
			name:       "state.toml only: nothing to migrate",
			configFile: false,
			stateFile:  true,
			want:       false,
		},
		{
			name:       "no files: fresh install",
			configFile: false,
			stateFile:  false,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			cfg := &Config{
				File:      filepath.Join(dir, "config.toml"),
				StateFile: filepath.Join(dir, "state.toml"),
			}
			if tt.configFile {
				require.NoError(t, os.WriteFile(cfg.File, []byte(""), 0o600))
			}
			if tt.stateFile {
				require.NoError(t, os.WriteFile(cfg.StateFile, []byte(""), 0o600))
			}

			assert.Equal(t, tt.want, cfg.ShouldMigrate())
		})
	}
}

func TestConfig_ShouldMigrate_unresolvedPaths(t *testing.T) {
	// InitConfig never ran: paths are empty, the trigger must stay off.
	cfg := &Config{}
	assert.False(t, cfg.ShouldMigrate())
}
