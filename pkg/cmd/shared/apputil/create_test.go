package apputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/keychain"
	"github.com/algolia/cli/test"
)

func TestConfigureProfile_SavesApplication(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	cfg := test.NewDefaultConfigStub()
	app := &dashboard.Application{
		ID: "APP1", Name: "My App", APIKey: "key-1", APIKeyUUID: "uuid-1",
	}

	require.NoError(t, ConfigureProfile(io, cfg, app, "", true))

	assert.Equal(t,
		test.SavedApplication{Alias: "my app", APIKeyUUID: "uuid-1", APIKey: "key-1"},
		cfg.SavedApps["APP1"])
	assert.Equal(t, "APP1", cfg.CurrentAppID)
}

func TestConfigureProfile_ExplicitNameAndNoDefault(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	cfg := test.NewDefaultConfigStub()
	app := &dashboard.Application{ID: "APP1", Name: "My App", APIKey: "key-1"}

	require.NoError(t, ConfigureProfile(io, cfg, app, "Prod", false))

	assert.Equal(t, "prod", cfg.SavedApps["APP1"].Alias)
	assert.Empty(t, cfg.CurrentAppID)
}

func TestConfigureProfile_AliasCollisionDerivesUniqueAlias(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	cfg := test.NewDefaultConfigStub()
	require.NoError(t, cfg.SaveApplication("OTHER", "my app", "", "other-key", false))

	app := &dashboard.Application{ID: "APP1", Name: "My App", APIKey: "key-1"}
	require.NoError(t, ConfigureProfile(io, cfg, app, "", true))

	assert.Equal(t, "my app-app1", cfg.SavedApps["APP1"].Alias)
}

func TestReuseExistingAPIKey_FromKeychain(t *testing.T) {
	keyring.MockInit()
	require.NoError(t, keychain.SaveAppSecrets("APP1",
		keychain.AppSecrets{APIKey: "stored-key"}))
	cfg := test.NewDefaultConfigStub()
	app := &dashboard.Application{ID: "APP1"}

	assert.True(t, ReuseExistingAPIKey(cfg, app))
	assert.Equal(t, "stored-key", app.APIKey)
}

func TestReuseExistingAPIKey_FromLegacyProfile(t *testing.T) {
	keyring.MockInit() // empty keychain → falls through to config.toml profiles
	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "legacy", ApplicationID: "APP1", APIKey: "legacy-key"},
	})
	app := &dashboard.Application{ID: "APP1"}

	assert.True(t, ReuseExistingAPIKey(cfg, app))
	assert.Equal(t, "legacy-key", app.APIKey)
}

func TestReuseExistingAPIKey_NotFound(t *testing.T) {
	keyring.MockInit()
	cfg := test.NewDefaultConfigStub()
	app := &dashboard.Application{ID: "APP1"}

	assert.False(t, ReuseExistingAPIKey(cfg, app))
	assert.Empty(t, app.APIKey)
}
