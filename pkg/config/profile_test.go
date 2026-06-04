package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/config/state"
)

func initTestViper(configFile string) {
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	_ = viper.ReadInConfig()
}

// seedState writes a state.toml with the given applications and current app,
// and returns its path. The keychain must be mocked separately.
func seedState(t *testing.T, current string, apps ...*state.ApplicationState) string {
	t.Helper()
	statePath := filepath.Join(t.TempDir(), "state.toml")
	st := state.New()
	for _, app := range apps {
		st.SetApp(app)
	}
	if current != "" {
		st.SetCurrent(current)
	}
	require.NoError(t, st.Save(statePath))
	return statePath
}

func TestProfile_ResolvesFromStateAndKeychain(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := seedState(t, "APP_STATE", &state.ApplicationState{
		ApplicationID: "APP_STATE",
		Alias:         "prod",
		CrawlerUserID: "crawler-user",
		SearchHosts:   []string{"h1", "h2"},
	})
	require.NoError(t, state.SetSecret("APP_STATE", state.SecretAPIKey, "key-state"))
	require.NoError(t, state.SetSecret("APP_STATE", state.SecretCrawlerAPIKey, "crawler-key"))

	p := &Profile{statePath: statePath}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP_STATE", appID)

	apiKey, err := p.GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "key-state", apiKey)

	assert.Equal(t, []string{"h1", "h2"}, p.GetSearchHosts())

	uid, err := p.GetCrawlerUserID()
	require.NoError(t, err)
	assert.Equal(t, "crawler-user", uid)

	crawlerKey, err := p.GetCrawlerAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "crawler-key", crawlerKey)
}

func TestProfile_EnvWinsOverState(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := seedState(t, "APP_STATE", &state.ApplicationState{
		ApplicationID: "APP_STATE",
		Alias:         "prod",
	})
	require.NoError(t, state.SetSecret("APP_STATE", state.SecretAPIKey, "key-state"))

	t.Setenv("ALGOLIA_APPLICATION_ID", "ENV_APP")
	t.Setenv("ALGOLIA_API_KEY", "env-key")

	p := &Profile{statePath: statePath}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "ENV_APP", appID)

	apiKey, err := p.GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "env-key", apiKey)
}

func TestProfile_FlagWinsOverState(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := seedState(t, "APP_STATE", &state.ApplicationState{
		ApplicationID: "APP_STATE",
		Alias:         "prod",
	})
	require.NoError(t, state.SetSecret("APP_STATE", state.SecretAPIKey, "key-state"))

	p := &Profile{statePath: statePath, ApplicationID: "FLAG_APP", APIKey: "flag-key"}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "FLAG_APP", appID)

	apiKey, err := p.GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "flag-key", apiKey)
}

func TestProfile_ProfileFlagResolvedViaAlias(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "a"},
		&state.ApplicationState{ApplicationID: "APP_B", Alias: "b"},
	)
	require.NoError(t, state.SetSecret("APP_A", state.SecretAPIKey, "key-a"))
	require.NoError(t, state.SetSecret("APP_B", state.SecretAPIKey, "key-b"))

	// `--profile b` resolves against the stored alias, not the current app.
	p := &Profile{statePath: statePath, Name: "b"}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP_B", appID)

	apiKey, err := p.GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "key-b", apiKey)
}

func TestProfile_UnknownProfileFlagFallsBackToCurrent(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := seedState(t, "APP_A",
		&state.ApplicationState{ApplicationID: "APP_A", Alias: "a"},
	)
	require.NoError(t, state.SetSecret("APP_A", state.SecretAPIKey, "key-a"))

	// `--profile does-not-exist` is not a known alias: resolve to current.
	p := &Profile{statePath: statePath, Name: "does-not-exist"}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP_A", appID)
}

func TestProfile_FallsBackToConfigTomlWhenNoState(t *testing.T) {
	keyring.MockInit()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[prod]
application_id = "APP_CFG"
api_key = "key-cfg"
default = true
`), 0o600))
	initTestViper(configPath)

	// state.toml does not exist on disk → fall back to config.toml.
	p := &Profile{statePath: filepath.Join(dir, "state.toml")}

	appID, err := p.GetApplicationID()
	require.NoError(t, err)
	assert.Equal(t, "APP_CFG", appID)

	apiKey, err := p.GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "key-cfg", apiKey)
}

func TestAddProfile_WritesStateAndKeychain(t *testing.T) {
	keyring.MockInit()
	viper.Reset()

	statePath := filepath.Join(t.TempDir(), "state.toml")

	// First CLI invocation: add profile A. With no current app yet, it
	// becomes the current application.
	profileA := &Profile{
		statePath:     statePath,
		Name:          "app-a",
		ApplicationID: "APP_A_ID",
		APIKey:        "key-a",
		APIKeyUUID:    "uuid-a",
	}
	require.NoError(t, profileA.Add())

	// Second CLI invocation: add profile B (not default), preserving A.
	profileB := &Profile{
		statePath:     statePath,
		Name:          "app-b",
		ApplicationID: "APP_B_ID",
		APIKey:        "key-b",
	}
	require.NoError(t, profileB.Add())

	// Read back via Config pointed at the same state.toml.
	cfg := &Config{}
	cfg.CurrentProfile.statePath = statePath
	profiles := cfg.ConfiguredProfiles()

	profilesByName := make(map[string]*Profile)
	for _, p := range profiles {
		profilesByName[p.Name] = p
	}

	assert.Len(t, profiles, 2, "both profiles should be preserved on disk")
	assert.Equal(t, "APP_A_ID", profilesByName["app-a"].ApplicationID)
	assert.Equal(t, "APP_B_ID", profilesByName["app-b"].ApplicationID)
	assert.Equal(t, "uuid-a", profilesByName["app-a"].APIKeyUUID)

	// The first application added becomes the current (default) one.
	assert.True(t, profilesByName["app-a"].Default)
	assert.False(t, profilesByName["app-b"].Default)

	// config.toml is never written; secrets live in the keychain only.
	keyA, err := state.GetSecret("APP_A_ID", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "key-a", keyA)
	keyB, err := state.GetSecret("APP_B_ID", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "key-b", keyB)

	// ConfiguredProfiles must not leak secrets into the profile metadata.
	assert.Empty(t, profilesByName["app-a"].APIKey)
}
