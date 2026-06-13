package config

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/keychain"
)

// ShouldMigrate reports whether the one-time config.toml → state.toml +
// keychain migration still has to run: config.toml exists and state.toml
// (only written on success, so doubling as the "migrated" marker) does not.
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
// + OS keychain); config.toml is never modified. Keychain first, state.toml
// last (atomic): a failure leaves state.toml absent, so the migration retries
// on the next run.
func (c *Config) Migrate() error {
	// An unparseable config.toml must not mark the migration as done: abort
	// before writing state.toml so it retries once the file is fixed.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	state := &State{Applications: map[string]ApplicationState{}}

	for _, profile := range c.migratableProfiles() {
		secrets := keychain.AppSecrets{
			APIKey:        profile.APIKey,
			CrawlerAPIKey: viper.GetString(profile.GetFieldName("crawler_api_key")),
		}
		if err := keychain.SaveAppSecrets(profile.ApplicationID, secrets); err != nil {
			return err
		}

		state.UpsertApplication(profile.ApplicationID, ApplicationState{
			Alias:         profile.Name,
			SearchHosts:   profile.SearchHosts,
			CrawlerUserID: viper.GetString(profile.GetFieldName("crawler_user_id")),
		})
		if profile.Default {
			state.SetCurrentApplication(profile.ApplicationID)
		}
	}

	return state.Save(c.StateFile)
}

// migratableProfiles applies the skip rules: profiles without application_id
// or api_key are dropped, admin_api_key never migrates, and the default
// profile wins when several share an application_id. Name order keeps the
// conflict resolution deterministic (ConfiguredProfiles iterates a map).
func (c *Config) migratableProfiles() []*Profile {
	// Decode the profiles here rather than through ConfiguredProfiles, whose
	// log.Fatalf on an undecodable entry would brick every command at startup.
	configs := viper.AllSettings()
	profiles := make([]*Profile, 0, len(configs))
	for name := range configs {
		profile := &Profile{Name: name}
		if err := viper.UnmarshalKey(name, profile); err != nil {
			log.Warnf("config migration: skipping profile %q: %s", name, err)
			continue
		}
		profiles = append(profiles, profile)
	}
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
