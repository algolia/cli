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

func TestConfig_LoadStateCachesAndToleratesMissing(t *testing.T) {
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	st := cfg.loadState()
	require.NotNil(t, st)
	assert.Empty(t, st.CurrentApplicationID)

	// Second call returns the same cached pointer (loaded once).
	assert.Same(t, st, cfg.loadState())
}

func TestConfig_LoadStateReadsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(
		path,
		[]byte(
			"current_application_id = \"APP1\"\n\n[applications.APP1]\nalias = \"prod\"\n",
		),
		0o600,
	))

	cfg := &Config{StateFile: path}
	st := cfg.loadState()
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
}

func TestConfig_ActiveApplicationID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(
		path,
		[]byte(
			"current_application_id = \"CURRENT\"\n\n[applications.ALIASED]\nalias = \"prod\"\n",
		),
		0o600,
	))

	t.Run("env wins", func(t *testing.T) {
		t.Setenv("ALGOLIA_APPLICATION_ID", "ENV_APP")
		cfg := &Config{StateFile: path}
		assert.Equal(t, "ENV_APP", cfg.activeApplicationID())
	})

	t.Run("profile alias resolves", func(t *testing.T) {
		cfg := &Config{StateFile: path}
		cfg.CurrentProfile.Name = "prod"
		assert.Equal(t, "ALIASED", cfg.activeApplicationID())
	})

	t.Run("falls back to current_application_id", func(t *testing.T) {
		cfg := &Config{StateFile: path}
		assert.Equal(t, "CURRENT", cfg.activeApplicationID())
	})

	t.Run("unknown profile alias defers to config.toml, not current", func(t *testing.T) {
		cfg := &Config{StateFile: path}
		cfg.CurrentProfile.Name = "nope"
		assert.Empty(t, cfg.activeApplicationID()) // "" → legacy profile-by-name, not CURRENT
	})
}

func TestConfig_AppSecretsForCaches(t *testing.T) {
	keyring.MockInit()
	require.NoError(t, keychain.SaveAppSecrets("APP1", keychain.AppSecrets{APIKey: "key-1"}))

	cfg := &Config{}
	got := cfg.appSecretsFor("APP1")
	require.NotNil(t, got)
	assert.Equal(t, "key-1", got.APIKey)

	// Missing app → nil, and a keychain error must not panic.
	assert.Nil(t, cfg.appSecretsFor("MISSING"))

	// The nil result is cached: a later keychain write isn't picked up mid-command.
	require.NoError(t, keychain.SaveAppSecrets("MISSING", keychain.AppSecrets{APIKey: "late"}))
	assert.Nil(t, cfg.appSecretsFor("MISSING"))
}

func TestProfile_GetApplicationID_NewModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte("current_application_id = \"APP1\"\n"), 0o600))

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	appID, err := cfg.Profile().GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP1", appID)
}

func TestProfile_GetAPIKey_FromKeychain(t *testing.T) {
	keyring.MockInit()
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte("current_application_id = \"APP1\"\n"), 0o600))
	require.NoError(t, keychain.SaveAppSecrets("APP1", keychain.AppSecrets{APIKey: "secret-key"}))

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	key, err := cfg.Profile().GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "secret-key", key)
}

func TestProfile_GetCrawlerAPIKey_FromKeychain(t *testing.T) {
	keyring.MockInit()
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte("current_application_id = \"APP1\"\n"), 0o600))
	require.NoError(t, keychain.SaveAppSecrets("APP1",
		keychain.AppSecrets{APIKey: "k", CrawlerAPIKey: "crawler-key"}))

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	key, err := cfg.Profile().GetCrawlerAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "crawler-key", key)
}

func TestProfile_EnvWinsOverKeychain(t *testing.T) {
	keyring.MockInit()
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte("current_application_id = \"APP1\"\n"), 0o600))
	require.NoError(t, keychain.SaveAppSecrets("APP1", keychain.AppSecrets{APIKey: "keychain-key"}))
	t.Setenv("ALGOLIA_API_KEY", "env-key")

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	key, err := cfg.Profile().GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "env-key", key) // env beats the keychain
}

func TestProfile_FallsBackToConfigToml(t *testing.T) {
	// No state.toml entry and no keychain → resolve from the legacy config.toml.
	configFile := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(
		configFile,
		[]byte(
			"[legacy]\napplication_id = \"LEGACY_APP\"\napi_key = \"legacy-key\"\ndefault = true\n",
		),
		0o600,
	))
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	require.NoError(t, viper.ReadInConfig())
	t.Cleanup(viper.Reset)

	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "absent.toml")}
	cfg.CurrentProfile.config = cfg

	appID, err := cfg.Profile().GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "LEGACY_APP", appID)

	key, err := cfg.Profile().GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "legacy-key", key)
}

func TestProfile_DefaultedNameDoesNotMaskCurrentApplication(t *testing.T) {
	// CheckAuth calls LoadDefault() before any getter runs. The defaulted
	// legacy profile name must NOT be mistaken for an explicit --profile flag,
	// otherwise the legacy default application silently wins over state.toml's
	// current_application_id.
	keyring.MockInit()
	statePath := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(
		t,
		os.WriteFile(statePath, []byte("current_application_id = \"APP1\"\n"), 0o600),
	)
	require.NoError(t, keychain.SaveAppSecrets("APP1", keychain.AppSecrets{APIKey: "app1-key"}))

	configFile := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(
		configFile,
		[]byte("[legacy]\napplication_id = \"LEGACY\"\napi_key = \"legacy-key\"\ndefault = true\n"),
		0o600,
	))
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	require.NoError(t, viper.ReadInConfig())
	t.Cleanup(viper.Reset)

	cfg := &Config{StateFile: statePath}
	cfg.CurrentProfile.config = cfg

	cfg.CurrentProfile.LoadDefault() // what CheckAuth does when no --profile is given
	require.Equal(t, "legacy", cfg.CurrentProfile.Name)

	appID, err := cfg.Profile().GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP1", appID)

	key, err := cfg.Profile().GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "app1-key", key)
}

func TestProfile_GetAPIKey_ActiveAppWithoutKeyErrors(t *testing.T) {
	keyring.MockInit()
	statePath := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(
		t,
		os.WriteFile(statePath, []byte("current_application_id = \"APP1\"\n"), 0o600),
	)

	// A legacy default profile whose key must NOT leak for the resolved APP1.
	configFile := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(
		configFile,
		[]byte("[legacy]\napplication_id = \"LEGACY\"\napi_key = \"legacy-key\"\ndefault = true\n"),
		0o600,
	))
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	require.NoError(t, viper.ReadInConfig())
	t.Cleanup(viper.Reset)

	cfg := &Config{StateFile: statePath}
	cfg.CurrentProfile.config = cfg

	// APP1 resolved from state but no keychain key → error, never "legacy-key".
	_, err := cfg.Profile().GetAPIKey()
	require.Error(t, err)
}

func TestProfile_GetSearchHosts_FromState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte(
		"current_application_id = \"APP1\"\n\n[applications.APP1]\nalias = \"prod\"\nsearch_hosts = [\"h1\", \"h2\"]\n",
	), 0o600))

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	assert.Equal(t, []string{"h1", "h2"}, cfg.Profile().GetSearchHosts())
}

func TestProfile_GetSearchHosts_StateEmptyFallsBackToConfigToml(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(
		t,
		os.WriteFile(statePath, []byte("current_application_id = \"APP1\"\n"), 0o600),
	)

	configFile := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(
		configFile,
		[]byte(
			"[legacy]\napplication_id = \"APP1\"\nsearch_hosts = [\"legacy-host\"]\ndefault = true\n",
		),
		0o600,
	))
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	require.NoError(t, viper.ReadInConfig())
	t.Cleanup(viper.Reset)

	cfg := &Config{StateFile: statePath}
	cfg.CurrentProfile.config = cfg

	// No hosts in state.toml for APP1: the legacy lookup still answers while
	// config.toml exists.
	assert.Equal(t, []string{"legacy-host"}, cfg.Profile().GetSearchHosts())
}

func TestProfile_GetCrawlerUserID_FromState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte(
		"current_application_id = \"APP1\"\n\n[applications.APP1]\nalias = \"prod\"\ncrawler_user_id = \"crawler-user\"\n",
	), 0o600))

	cfg := &Config{StateFile: path}
	cfg.CurrentProfile.config = cfg

	userID, err := cfg.Profile().GetCrawlerUserID()
	require.NoError(t, err)
	assert.Equal(t, "crawler-user", userID)
}
