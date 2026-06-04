package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/config/state"
)

// newMigrateConfig writes config.toml with the given content and returns a
// Config wired to a temp config.toml + state.toml plus the state path.
func newMigrateConfig(t *testing.T, configTOML string) (*Config, string) {
	t.Helper()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	statePath := filepath.Join(dir, "state.toml")

	if configTOML != "" {
		require.NoError(t, os.WriteFile(configPath, []byte(configTOML), 0o600))
	}

	c := &Config{File: configPath}
	c.CurrentProfile.statePath = statePath
	return c, statePath
}

func TestMigrateIfNeeded_HappyPath(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[prod]
application_id = "APP_PROD"
api_key = "key-prod"
default = true

[staging]
application_id = "APP_STAGING"
api_key = "key-staging"
search_hosts = ["host1", "host2"]
`)

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	st, err := state.Load(statePath)
	require.NoError(t, err)

	assert.Equal(t, "APP_PROD", st.CurrentApplicationID, "default profile becomes current")

	prod := st.App("APP_PROD")
	require.NotNil(t, prod)
	assert.Equal(t, "prod", prod.Alias)

	staging := st.App("APP_STAGING")
	require.NotNil(t, staging)
	assert.Equal(t, "staging", staging.Alias)
	assert.Equal(t, []string{"host1", "host2"}, staging.SearchHosts)

	// Secrets land in the keychain, not in state.toml.
	prodKey, err := state.GetSecret("APP_PROD", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "key-prod", prodKey)
	stagingKey, err := state.GetSecret("APP_STAGING", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "key-staging", stagingKey)

	raw, err := os.ReadFile(statePath)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), "key-prod", "API keys must never be written to state.toml")
}

func TestMigrateIfNeeded_SkipsEmptyAPIKey(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[good]
application_id = "APP_GOOD"
api_key = "key-good"
default = true

[empty]
application_id = "APP_EMPTY"
`)

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	st, err := state.Load(statePath)
	require.NoError(t, err)
	assert.NotNil(t, st.App("APP_GOOD"))
	assert.Nil(t, st.App("APP_EMPTY"), "profile without API key is skipped")
	assert.Contains(t, stderr.String(), `Skipped profile "empty"`)
	assert.Contains(t, stderr.String(), "algolia application select")
}

func TestMigrateIfNeeded_ConflictKeepsDefault(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[secondary]
application_id = "SHARED_APP"
api_key = "key-secondary"

[primary]
application_id = "SHARED_APP"
api_key = "key-primary"
default = true
`)

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	st, err := state.Load(statePath)
	require.NoError(t, err)

	app := st.App("SHARED_APP")
	require.NotNil(t, app)
	assert.Equal(t, "primary", app.Alias, "the default profile wins the shared application_id")
	assert.Equal(t, "SHARED_APP", st.CurrentApplicationID)

	key, err := state.GetSecret("SHARED_APP", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "key-primary", key)

	assert.Contains(t, stderr.String(), `Skipped profile "secondary"`)
	assert.Contains(t, stderr.String(), "already configured by profile")
}

func TestMigrateIfNeeded_AdminKeyNoticeAndNotMigrated(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[legacy]
application_id = "APP_LEGACY"
admin_api_key = "admin-secret"
`)

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	assert.Contains(t, stderr.String(), "ALGOLIA_ADMIN_API_KEY")

	// admin_api_key alone (no api_key) means the profile is skipped entirely.
	st, err := state.Load(statePath)
	require.NoError(t, err)
	assert.Nil(t, st.App("APP_LEGACY"))

	stored, err := state.GetSecret("APP_LEGACY", state.SecretAPIKey)
	require.NoError(t, err)
	assert.Empty(t, stored, "admin keys are never stored")
}

func TestMigrateIfNeeded_CrawlerCredentials(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[prod]
application_id = "APP_PROD"
api_key = "key-prod"
crawler_user_id = "crawler-user"
crawler_api_key = "crawler-secret"
default = true
`)

	c.MigrateIfNeeded(&bytes.Buffer{})

	st, err := state.Load(statePath)
	require.NoError(t, err)

	app := st.App("APP_PROD")
	require.NotNil(t, app)
	assert.Equal(t, "crawler-user", app.CrawlerUserID, "crawler user ID is a non-secret in state.toml")

	crawlerKey, err := state.GetSecret("APP_PROD", state.SecretCrawlerAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "crawler-secret", crawlerKey)

	raw, err := os.ReadFile(statePath)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), "crawler-secret")
}

func TestMigrateIfNeeded_SingleProfileWithoutDefaultBecomesCurrent(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[only]
application_id = "APP_ONLY"
api_key = "key-only"
`)

	c.MigrateIfNeeded(&bytes.Buffer{})

	st, err := state.Load(statePath)
	require.NoError(t, err)
	assert.Equal(t, "APP_ONLY", st.CurrentApplicationID)
}

func TestMigrateIfNeeded_IdempotentWhenStateExists(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, `
[prod]
application_id = "APP_PROD"
api_key = "key-prod"
default = true
`)

	// Pre-existing state.toml means the migration is already done.
	existing := state.New()
	existing.SetApp(&state.ApplicationState{ApplicationID: "EXISTING", Alias: "existing"})
	existing.SetCurrent("EXISTING")
	require.NoError(t, existing.Save(statePath))

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	st, err := state.Load(statePath)
	require.NoError(t, err)
	assert.Equal(t, "EXISTING", st.CurrentApplicationID, "existing state.toml is left untouched")
	assert.Nil(t, st.App("APP_PROD"))
	assert.Empty(t, stderr.String())
}

func TestMigrateIfNeeded_NoConfigFileIsNoop(t *testing.T) {
	keyring.MockInit()

	c, statePath := newMigrateConfig(t, "") // no config.toml written

	var stderr bytes.Buffer
	c.MigrateIfNeeded(&stderr)

	_, err := os.Stat(statePath)
	assert.True(t, os.IsNotExist(err), "no state.toml is created without a legacy config")
	assert.Empty(t, stderr.String())
}
