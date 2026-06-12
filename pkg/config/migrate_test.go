package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/keychain"
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

// migrationConfig writes a config.toml with the given content, points the
// global viper at it (ConfiguredProfiles reads through viper) and returns a
// Config ready to migrate.
func migrationConfig(t *testing.T, content string) *Config {
	t.Helper()

	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.toml")
	require.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))

	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	require.NoError(t, viper.ReadInConfig())
	t.Cleanup(viper.Reset)

	return &Config{
		File:      configFile,
		StateFile: filepath.Join(dir, "state.toml"),
	}
}

func TestConfig_Migrate(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, `[prod]
application_id = "APP1"
api_key = "key-1"
crawler_api_key = "crawler-1"
default = true

[dev]
application_id = "APP2"
api_key = "key-2"
`)

	require.NoError(t, cfg.Migrate())

	// Secrets land in the keychain, crawler key included when set.
	prod, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, prod)
	assert.Equal(t, "key-1", prod.APIKey)
	assert.Equal(t, "crawler-1", prod.CrawlerAPIKey)

	dev, err := keychain.LoadAppSecrets("APP2")
	require.NoError(t, err)
	require.NotNil(t, dev)
	assert.Equal(t, "key-2", dev.APIKey)
	assert.Empty(t, dev.CrawlerAPIKey)

	// state.toml: one entry per application, alias = profile name, current
	// application = the default profile's one.
	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
	assert.Equal(t, "dev", st.Applications["APP2"].Alias)
	assert.Empty(t, st.Applications["APP1"].APIKeyUUID) // unknown for legacy keys

	// state.toml now exists: the trigger turns off.
	assert.False(t, cfg.ShouldMigrate())
}

func TestConfig_Migrate_NoDefaultProfileLeavesCurrentEmpty(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, `[dev]
application_id = "APP2"
api_key = "key-2"
`)

	require.NoError(t, cfg.Migrate())

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Empty(t, st.CurrentApplicationID)
	assert.Equal(t, "dev", st.Applications["APP2"].Alias)
}

func TestConfig_Migrate_EmptyConfigStillWritesState(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, "")

	require.NoError(t, cfg.Migrate())

	// An empty state.toml must exist, otherwise the migration would re-run
	// (and re-log) on every command.
	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Empty(t, st.CurrentApplicationID)
	assert.Empty(t, st.Applications)
	assert.False(t, cfg.ShouldMigrate())
}

func TestConfig_Migrate_KeychainFailureLeavesStateAbsent(t *testing.T) {
	keyring.MockInitWithError(keyring.ErrUnsupportedPlatform)
	cfg := migrationConfig(t, `[prod]
application_id = "APP1"
api_key = "key-1"
`)

	require.Error(t, cfg.Migrate())

	// state.toml untouched: ShouldMigrate keeps firing so the migration
	// retries on the next run.
	assert.NoFileExists(t, cfg.StateFile)
	assert.True(t, cfg.ShouldMigrate())
}
