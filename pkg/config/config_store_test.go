package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/config/state"
)

// newStateConfig returns a Config whose current profile points at statePath,
// mirroring what InitConfig does in production.
func newStateConfig(statePath string) *Config {
	c := &Config{}
	c.CurrentProfile.statePath = statePath
	return c
}

func TestConfig_ConfiguredProfiles_FromState(t *testing.T) {
	keyring.MockInit()

	statePath := seedState(
		t,
		"APP_B",
		&state.ApplicationState{
			ApplicationID: "APP_A",
			Alias:         "alpha",
			SearchHosts:   []string{"h1"},
		},
		&state.ApplicationState{ApplicationID: "APP_B", Alias: "beta"},
		&state.ApplicationState{ApplicationID: "APP_C"}, // no alias -> name falls back to ID
	)

	cfg := newStateConfig(statePath)
	profiles := cfg.ConfiguredProfiles()
	require.Len(t, profiles, 3)

	// Deterministic order: sorted by application ID.
	assert.Equal(t, "alpha", profiles[0].Name)
	assert.Equal(t, "APP_A", profiles[0].ApplicationID)
	assert.Equal(t, []string{"h1"}, profiles[0].SearchHosts)
	assert.False(t, profiles[0].Default)

	assert.Equal(t, "beta", profiles[1].Name)
	assert.True(t, profiles[1].Default, "current application is the default profile")

	assert.Equal(t, "APP_C", profiles[2].Name, "missing alias falls back to application ID")

	// Secrets are never populated by ConfiguredProfiles.
	for _, p := range profiles {
		assert.Empty(t, p.APIKey)
	}
}

func TestConfig_ConfiguredProfiles_EmptyWhenNoStatePath(t *testing.T) {
	cfg := &Config{} // bare config, as in unit tests without InitConfig
	assert.Empty(t, cfg.ConfiguredProfiles())
	assert.Nil(t, cfg.Default())
	assert.Empty(t, cfg.ProfileNames())
	assert.False(t, cfg.ProfileExists("anything"))
}

func TestConfig_SetDefaultProfile_ByAliasAndID(t *testing.T) {
	keyring.MockInit()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "alpha"},
		&state.ApplicationState{ApplicationID: "APP_B", Alias: "beta"},
	)
	cfg := newStateConfig(statePath)

	// Resolve by alias.
	require.NoError(t, cfg.SetDefaultProfile("beta"))
	assert.Equal(t, "beta", cfg.Default().Name)

	// Resolve by application ID.
	require.NoError(t, cfg.SetDefaultProfile("APP_A"))
	assert.Equal(t, "alpha", cfg.Default().Name)

	// Unknown profile errors.
	assert.Error(t, cfg.SetDefaultProfile("nope"))
}

func TestConfig_RemoveProfile_DeletesStateAndSecrets(t *testing.T) {
	keyring.MockInit()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "alpha"},
		&state.ApplicationState{ApplicationID: "APP_B", Alias: "beta"},
	)
	require.NoError(t, state.SetSecret("APP_A", state.SecretAPIKey, "key-a"))
	require.NoError(t, state.SetSecret("APP_A", state.SecretCrawlerAPIKey, "crawler-a"))

	cfg := newStateConfig(statePath)

	require.NoError(t, cfg.RemoveProfile("alpha"))

	// State entry removed and current pointer cleared (it referenced APP_A).
	profiles := cfg.ConfiguredProfiles()
	require.Len(t, profiles, 1)
	assert.Equal(t, "beta", profiles[0].Name)
	assert.Nil(t, cfg.Default(), "removing the current app clears the default")

	// Secrets are deleted from the keychain.
	apiKey, err := state.GetSecret("APP_A", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Empty(t, apiKey)
	crawlerKey, err := state.GetSecret("APP_A", state.SecretCrawlerAPIKey)
	require.NoError(t, err)
	assert.Empty(t, crawlerKey)

	// Removing an unknown profile errors.
	assert.Error(t, cfg.RemoveProfile("nope"))
}

func TestConfig_ApplicationIDLookups(t *testing.T) {
	keyring.MockInit()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "alpha"},
	)
	cfg := newStateConfig(statePath)

	exists, name := cfg.ApplicationIDExists("APP_A")
	assert.True(t, exists)
	assert.Equal(t, "alpha", name)

	exists, _ = cfg.ApplicationIDExists("UNKNOWN")
	assert.False(t, exists)

	exists, appID := cfg.ApplicationIDForProfile("alpha")
	assert.True(t, exists)
	assert.Equal(t, "APP_A", appID)

	// Resolvable by application ID too.
	exists, appID = cfg.ApplicationIDForProfile("APP_A")
	assert.True(t, exists)
	assert.Equal(t, "APP_A", appID)

	assert.True(t, cfg.ProfileExists("alpha"))
	assert.True(t, cfg.ProfileExists("APP_A"))
	assert.False(t, cfg.ProfileExists("nope"))
}

func TestConfig_SetCrawlerAuth_StateAndKeychain(t *testing.T) {
	keyring.MockInit()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "alpha"},
	)
	cfg := newStateConfig(statePath)

	require.NoError(t, cfg.SetCrawlerAuth("alpha", "crawler-user", "crawler-key"))

	// Non-secret crawler_user_id lands in state.toml.
	st, err := state.Load(statePath)
	require.NoError(t, err)
	assert.Equal(t, "crawler-user", st.App("APP_A").CrawlerUserID)

	// The crawler API key lands in the keychain.
	key, err := state.GetSecret("APP_A", state.SecretCrawlerAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "crawler-key", key)

	// Unknown profile errors.
	assert.Error(t, cfg.SetCrawlerAuth("nope", "u", "k"))
}
