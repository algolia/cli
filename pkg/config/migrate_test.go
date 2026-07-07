package config

import (
	"os"
	"path/filepath"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"
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

// migrationConfig writes a config.toml, points the global viper at it
// (ConfiguredProfiles reads through viper) and returns a Config ready to migrate.
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
crawler_user_id = "crawler-user"
search_hosts = ["h1.algolia.net", "h2.algolia.net"]
default = true

[dev]
application_id = "APP2"
api_key = "key-2"
`)

	require.NoError(t, cfg.Migrate())

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

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
	assert.Equal(t, "dev", st.Applications["APP2"].Alias)
	assert.Empty(t, st.Applications["APP1"].APIKeyUUID) // unknown for legacy keys
	assert.Equal(
		t,
		[]string{"h1.algolia.net", "h2.algolia.net"},
		st.Applications["APP1"].SearchHosts,
	)
	assert.Equal(t, "crawler-user", st.Applications["APP1"].CrawlerUserID)
	assert.Empty(t, st.Applications["APP2"].SearchHosts)
	assert.Empty(t, st.Applications["APP2"].CrawlerUserID)

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

	// An empty state.toml must exist, otherwise the migration re-runs forever.
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

	// state.toml untouched: the migration retries on the next run.
	assert.NoFileExists(t, cfg.StateFile)
	assert.True(t, cfg.ShouldMigrate())
}

func TestConfig_Migrate_SkipRules(t *testing.T) {
	keyring.MockInit()
	hook := logtest.NewGlobal()
	t.Cleanup(hook.Reset)

	cfg := migrationConfig(t, `[nokey]
application_id = "APP3"
api_key = ""

[noapp]
api_key = "key-x"

[adminonly]
application_id = "APP4"
admin_api_key = "admin-key"
`)

	require.NoError(t, cfg.Migrate())

	// Nothing migrated, but the trigger still turns off.
	for _, appID := range []string{"APP3", "APP4"} {
		secrets, err := keychain.LoadAppSecrets(appID)
		require.NoError(t, err)
		assert.Nil(t, secrets)
	}
	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Empty(t, st.Applications)
	assert.False(t, cfg.ShouldMigrate())

	logs := make([]string, 0, len(hook.AllEntries()))
	for _, entry := range hook.AllEntries() {
		logs = append(logs, entry.Message)
	}
	assert.Contains(t, logs,
		`config migration: skipping profile "nokey": empty api_key`)
	assert.Contains(t, logs,
		`config migration: skipping profile "noapp": no application_id`)
	assert.Contains(t, logs,
		`config migration: skipping profile "adminonly": empty api_key`)
	assert.Contains(
		t,
		logs,
		`config migration: profile "adminonly": admin_api_key is not migrated, use ALGOLIA_ADMIN_API_KEY or --api-key instead`,
	)
}

func TestConfig_Migrate_DuplicateApplicationKeepsDefault(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, `[backup]
application_id = "APP1"
api_key = "backup-key"

[prod]
application_id = "APP1"
api_key = "prod-key"
default = true
`)

	require.NoError(t, cfg.Migrate())

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	require.Len(t, st.Applications, 1)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
	assert.Equal(t, "APP1", st.CurrentApplicationID)

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, secrets)
	assert.Equal(t, "prod-key", secrets.APIKey)
}

func TestConfig_Migrate_AdminKeyAlongsideAPIKeyStillMigrates(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, `[prod]
application_id = "APP1"
api_key = "key-1"
admin_api_key = "admin-1"
default = true
`)

	require.NoError(t, cfg.Migrate())

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, secrets)
	assert.Equal(t, "key-1", secrets.APIKey)

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
}

func TestConfig_Migrate_UndecodableProfileSkipped(t *testing.T) {
	keyring.MockInit()
	cfg := migrationConfig(t, `telemetry = "off"

[prod]
application_id = "APP1"
api_key = "key-1"
default = true

[bad]
application_id = "APP2"
api_key = ["a", "b"]
`)

	// Undecodable entries (root scalar, wrong types) are skipped, not fatal.
	require.NoError(t, cfg.Migrate())

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	require.Len(t, st.Applications, 1)
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
}

func TestConfig_Migrate_UnreadableConfigAborts(t *testing.T) {
	keyring.MockInit()
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.toml")
	require.NoError(t, os.WriteFile(configFile, []byte("not [ valid ### toml\n"), 0o600))

	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	_ = viper.ReadInConfig() // swallowed, like InitConfig does
	t.Cleanup(viper.Reset)

	cfg := &Config{File: configFile, StateFile: filepath.Join(dir, "state.toml")}

	// No state.toml written: the migration retries once the file is fixed.
	require.Error(t, cfg.Migrate())
	assert.NoFileExists(t, cfg.StateFile)
	assert.True(t, cfg.ShouldMigrate())
}
