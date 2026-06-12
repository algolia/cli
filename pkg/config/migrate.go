package config

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
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

	for _, profile := range c.migratableProfiles() {
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

// migratableProfiles applies the migration skip rules to the config.toml
// profiles before any keychain write happens:
//
//   - admin_api_key never moves to the new model: one log line points to its
//     replacements, whether the profile migrates or not.
//   - A profile without application_id or with an empty api_key has nothing
//     usable to migrate: skipped with a log line.
//   - Profiles sharing the same application_id would overwrite each other's
//     keychain entry: the default = true profile wins, the others are logged
//     as conflicts and skipped.
//
// Profiles are processed in name order so conflict resolution and logs stay
// deterministic (ConfiguredProfiles iterates a map).
func (c *Config) migratableProfiles() []*Profile {
	profiles := c.ConfiguredProfiles()
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })

	selected := make([]*Profile, 0, len(profiles))
	owner := map[string]int{} // application ID → index in selected

	for _, profile := range profiles {
		if profile.AdminAPIKey != "" {
			log.Warnf(
				"config migration: profile %q: admin_api_key is not migrated, use ALGOLIA_ADMIN_API_KEY or --api-key instead",
				profile.Name,
			)
		}
		if profile.ApplicationID == "" {
			log.Warnf("config migration: skipping profile %q: no application_id", profile.Name)
			continue
		}
		if profile.APIKey == "" {
			log.Warnf("config migration: skipping profile %q: empty api_key", profile.Name)
			continue
		}
		if i, ok := owner[profile.ApplicationID]; ok {
			kept, dropped := selected[i], profile
			if profile.Default && !kept.Default {
				selected[i] = profile
				kept, dropped = profile, kept
			}
			log.Warnf(
				"config migration: skipping profile %q: application %q already migrated from profile %q",
				dropped.Name,
				dropped.ApplicationID,
				kept.Name,
			)
			continue
		}
		owner[profile.ApplicationID] = len(selected)
		selected = append(selected, profile)
	}

	return selected
}
