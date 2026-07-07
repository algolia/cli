package config

import (
	"github.com/algolia/cli/pkg/keychain"
)

// ActiveApplicationID exposes the application resolved by the new model
// (env → flag → state.toml alias → current application). Empty when only the
// legacy config.toml could answer.
func (c *Config) ActiveApplicationID() string {
	return c.activeApplicationID()
}

// APIKeyUUID returns the CLI-managed API key UUID recorded in state.toml for
// the given application, and whether one is stored. A configured application
// with no UUID (a legacy setup) returns ("", false).
func (c *Config) APIKeyUUID(appID string) (string, bool) {
	app, ok := c.loadState().Applications[appID]
	if !ok || app.APIKeyUUID == "" {
		return "", false
	}
	return app.APIKeyUUID, true
}

// ApplicationIDByAlias returns the application ID carrying the given alias in
// state.toml, and whether one was found.
func (c *Config) ApplicationIDByAlias(alias string) (string, bool) {
	return c.loadState().ApplicationByAlias(alias)
}

// SaveApplication persists an application's credentials in the new model.
// The keychain is written first so a failure never leaves state.toml pointing
// at a key that was not stored. Empty alias/apiKeyUUID preserve the values
// already in state.toml, and an existing crawler key is preserved.
//
// Note: a command that already resolved its active application keeps that
// resolution for the rest of the command (per-command cache).
func (c *Config) SaveApplication(appID, alias, apiKeyUUID, apiKey string, setCurrent bool) error {
	secrets, err := keychain.LoadAppSecrets(appID)
	if err != nil {
		return err
	}
	if secrets == nil {
		secrets = &keychain.AppSecrets{}
	}
	secrets.APIKey = apiKey
	if err := keychain.SaveAppSecrets(appID, *secrets); err != nil {
		return err
	}
	c.cacheSecrets(appID, secrets)

	st := c.loadState()
	app := st.Applications[appID]
	if alias != "" {
		app.Alias = alias
	}
	if apiKeyUUID != "" {
		app.APIKeyUUID = apiKeyUUID
	}
	st.UpsertApplication(appID, app)
	if setCurrent {
		st.SetCurrentApplication(appID)
	}

	return st.Save(c.StateFile)
}

// SetCrawlerAPIKey stores the crawler API key in the keychain entry of the
// given application, preserving the search API key (load-modify-save).
func (c *Config) SetCrawlerAPIKey(appID, crawlerAPIKey string) error {
	secrets, err := keychain.LoadAppSecrets(appID)
	if err != nil {
		return err
	}
	if secrets == nil {
		secrets = &keychain.AppSecrets{}
	}
	secrets.CrawlerAPIKey = crawlerAPIKey
	if err := keychain.SaveAppSecrets(appID, *secrets); err != nil {
		return err
	}
	c.cacheSecrets(appID, secrets)

	return nil
}

// cacheSecrets refreshes the per-command secrets cache after a write so reads
// in the same command observe the new values.
func (c *Config) cacheSecrets(appID string, secrets *keychain.AppSecrets) {
	c.secretsMu.Lock()
	defer c.secretsMu.Unlock()
	if c.secretsCache == nil {
		c.secretsCache = map[string]*keychain.AppSecrets{}
	}
	c.secretsCache[appID] = secrets
}
