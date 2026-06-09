package config

import (
	"os"
	"path/filepath"
	"testing"

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
