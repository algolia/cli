package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/keychain"
)

func TestConfig_SaveApplication_WritesKeychainThenState(t *testing.T) {
	keyring.MockInit()
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	require.NoError(t, cfg.SaveApplication("APP1", "prod", "uuid-1", "key-1", true))

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, secrets)
	assert.Equal(t, "key-1", secrets.APIKey)

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t,
		ApplicationState{APIKeyUUID: "uuid-1", Alias: "prod"},
		st.Applications["APP1"])
}

func TestConfig_SaveApplication_PreservesExistingValues(t *testing.T) {
	keyring.MockInit()
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	require.NoError(t, keychain.SaveAppSecrets("APP1",
		keychain.AppSecrets{APIKey: "old-key", CrawlerAPIKey: "crawler-key"}))
	require.NoError(t, cfg.SaveApplication("APP1", "custom", "uuid-1", "old-key", false))

	// Empty alias and UUID keep the stored values; the crawler key survives.
	require.NoError(t, cfg.SaveApplication("APP1", "", "", "new-key", false))

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	assert.Equal(t, "new-key", secrets.APIKey)
	assert.Equal(t, "crawler-key", secrets.CrawlerAPIKey)

	st, err := LoadState(cfg.StateFile)
	require.NoError(t, err)
	assert.Equal(t,
		ApplicationState{APIKeyUUID: "uuid-1", Alias: "custom"},
		st.Applications["APP1"])
	assert.Empty(t, st.CurrentApplicationID) // setCurrent was always false
}

func TestConfig_SaveApplication_KeychainErrorAborts(t *testing.T) {
	keyring.MockInitWithError(errors.New("keychain locked"))
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	require.Error(t, cfg.SaveApplication("APP1", "prod", "uuid-1", "key-1", true))

	// Keychain-first: state.toml must not exist after a keychain failure.
	_, err := os.Stat(cfg.StateFile)
	assert.True(t, os.IsNotExist(err))
}

func TestConfig_SaveApplication_RefreshesSecretsCache(t *testing.T) {
	keyring.MockInit()
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	// Prime the negative cache: nothing stored yet.
	require.Nil(t, cfg.appSecretsFor("APP1"))

	require.NoError(t, cfg.SaveApplication("APP1", "prod", "uuid-1", "key-1", true))

	got := cfg.appSecretsFor("APP1")
	require.NotNil(t, got)
	assert.Equal(t, "key-1", got.APIKey)
}

func TestConfig_SetCrawlerAPIKey_PreservesSearchKey(t *testing.T) {
	keyring.MockInit()
	cfg := &Config{}

	require.NoError(t, keychain.SaveAppSecrets("APP1",
		keychain.AppSecrets{APIKey: "search-key"}))
	require.NoError(t, cfg.SetCrawlerAPIKey("APP1", "crawler-key"))

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	assert.Equal(t, "search-key", secrets.APIKey)
	assert.Equal(t, "crawler-key", secrets.CrawlerAPIKey)
}

func TestConfig_SetCrawlerAPIKey_CreatesEntryWhenMissing(t *testing.T) {
	keyring.MockInit()
	cfg := &Config{}

	require.NoError(t, cfg.SetCrawlerAPIKey("APP1", "crawler-key"))

	secrets, err := keychain.LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, secrets)
	assert.Equal(t, "crawler-key", secrets.CrawlerAPIKey)
	assert.Empty(t, secrets.APIKey)
}

func TestConfig_ActiveApplicationIDAndAliasAccessors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(
		path,
		[]byte("current_application_id = \"APP1\"\n\n[applications.APP1]\nalias = \"prod\"\n"),
		0o600,
	))
	cfg := &Config{StateFile: path}

	assert.Equal(t, "APP1", cfg.ActiveApplicationID())

	appID, ok := cfg.ApplicationIDByAlias("prod")
	assert.True(t, ok)
	assert.Equal(t, "APP1", appID)
}
