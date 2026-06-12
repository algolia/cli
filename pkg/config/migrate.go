package config

import (
	"os"

	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/keychain"
)

// ShouldMigrate reports whether the one-time config.toml → state.toml +
// keychain migration still has to run: a legacy config.toml exists and
// state.toml does not. state.toml is only written when the migration (or any
// new-model write) succeeds, so its absence doubles as the "not migrated yet"
// marker — an aborted migration naturally retries on the next run.
func (c *Config) ShouldMigrate() bool {
	if c.File == "" {
		return false
	}
	if _, err := os.Stat(c.File); err != nil {
		return false
	}
	return !c.StateFileExists()
}

// Migrate moves the legacy config.toml profiles into the new model (state.toml
// + OS keychain). config.toml itself is never modified.
//
// Secrets go to the keychain first; state.toml is only written — atomically,
// via State.Save's temp + rename — once every profile's keys are stored. A
// keychain failure mid-run therefore leaves state.toml absent and the whole
// migration retries on the next command; entries already written are simply
// rewritten then. With nothing to migrate an empty state.toml still gets
// written, so ShouldMigrate stops firing on every command.
func (c *Config) Migrate() error {
	state := &State{Applications: map[string]ApplicationState{}}

	for _, profile := range c.ConfiguredProfiles() {
		secrets := keychain.AppSecrets{
			APIKey:        profile.APIKey,
			CrawlerAPIKey: viper.GetString(profile.GetFieldName("crawler_api_key")),
		}
		if err := keychain.SaveAppSecrets(profile.ApplicationID, secrets); err != nil {
			return err
		}

		// api_key_uuid is unknown for legacy keys: left empty until an API
		// lookup can backfill it.
		state.UpsertApplication(profile.ApplicationID, ApplicationState{Alias: profile.Name})
		if profile.Default {
			state.SetCurrentApplication(profile.ApplicationID)
		}
	}

	return state.Save(c.StateFile)
}
