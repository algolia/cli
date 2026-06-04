package config

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/config/state"
)

// legacyProfile is a flattened view of a single profile read from config.toml,
// including the crawler fields that are not part of the Profile struct.
type legacyProfile struct {
	Name          string
	ApplicationID string
	APIKey        string
	AdminAPIKey   string
	SearchHosts   []string
	Default       bool
	CrawlerUserID string
	CrawlerAPIKey string
}

// MigrateIfNeeded performs the one-time migration of config.toml into
// state.toml (non-secrets) + the OS keychain (secrets). It runs at most once:
// the presence of state.toml marks the migration as done. Migration is
// best-effort and never fatal — on any failure config.toml is left untouched so
// it remains usable as a read-only fallback (until removal in CLI v2.0).
//
// Skip rules:
//   - profiles with an empty api_key are skipped (logged with a next step);
//   - profiles sharing an application_id keep the one marked default = true,
//     the others are logged as conflicts and skipped;
//   - admin_api_key is never migrated; a one-line notice points the user at
//     ALGOLIA_ADMIN_API_KEY or the --api-key flag.
func (c *Config) MigrateIfNeeded(stderr io.Writer) {
	statePath := c.CurrentProfile.statePath
	if statePath == "" {
		statePath = state.DefaultPath()
	}

	// Already migrated: state.toml exists.
	if _, err := os.Stat(statePath); err == nil {
		return
	}

	configFile := c.File
	if configFile == "" {
		return
	}
	if _, err := os.Stat(configFile); err != nil {
		return // nothing legacy to migrate (fresh install)
	}

	legacy, err := readLegacyProfiles(configFile)
	if err != nil {
		fmt.Fprintf(stderr, "Could not read %s for migration: %s\n", configFile, err)
		return
	}
	if len(legacy) == 0 {
		return
	}

	st := migrateProfiles(legacy, stderr)

	if err := st.Save(statePath); err != nil {
		fmt.Fprintf(stderr, "Could not write %s during migration: %s\n", statePath, err)
		return
	}
}

// migrateProfiles builds the new State from legacy profiles, storing secrets in
// the keychain and logging skipped profiles and the admin-key notice.
func migrateProfiles(legacy []legacyProfile, stderr io.Writer) *state.State {
	// Default profiles first, then alphabetical, for deterministic conflict
	// resolution (the default among a shared application_id wins).
	sort.SliceStable(legacy, func(i, j int) bool {
		if legacy[i].Default != legacy[j].Default {
			return legacy[i].Default
		}
		return legacy[i].Name < legacy[j].Name
	})

	st := state.New()
	claimedBy := map[string]string{} // application_id -> alias that claimed it
	adminNoticeShown := false

	for _, p := range legacy {
		if p.AdminAPIKey != "" && !adminNoticeShown {
			fmt.Fprintf(
				stderr,
				"Note: admin API keys are not migrated. Set ALGOLIA_ADMIN_API_KEY or pass --api-key when a command needs one.\n",
			)
			adminNoticeShown = true
		}

		if p.APIKey == "" {
			fmt.Fprintf(
				stderr,
				"Skipped profile %q during migration: no API key stored. Run `algolia application select` to configure it.\n",
				p.Name,
			)
			continue
		}
		if p.ApplicationID == "" {
			fmt.Fprintf(
				stderr,
				"Skipped profile %q during migration: no application ID.\n",
				p.Name,
			)
			continue
		}
		if owner, taken := claimedBy[p.ApplicationID]; taken {
			fmt.Fprintf(
				stderr,
				"Skipped profile %q during migration: application %s is already configured by profile %q.\n",
				p.Name,
				p.ApplicationID,
				owner,
			)
			continue
		}

		st.SetApp(&state.ApplicationState{
			ApplicationID: p.ApplicationID,
			Alias:         p.Name,
			SearchHosts:   p.SearchHosts,
			CrawlerUserID: p.CrawlerUserID,
		})
		claimedBy[p.ApplicationID] = p.Name

		if err := state.SetSecret(p.ApplicationID, state.SecretAPIKey, p.APIKey); err != nil {
			fmt.Fprintf(stderr, "Could not store API key for profile %q: %s\n", p.Name, err)
		}
		if p.CrawlerAPIKey != "" {
			if err := state.SetSecret(p.ApplicationID, state.SecretCrawlerAPIKey, p.CrawlerAPIKey); err != nil {
				fmt.Fprintf(
					stderr,
					"Could not store crawler API key for profile %q: %s\n",
					p.Name,
					err,
				)
			}
		}

		if p.Default && st.CurrentApplicationID == "" {
			st.SetCurrent(p.ApplicationID)
		}
	}

	// No explicit default but a single application migrated: make it current
	// so commands work without --profile.
	if st.CurrentApplicationID == "" && len(claimedBy) == 1 {
		for appID := range claimedBy {
			st.SetCurrent(appID)
		}
	}

	return st
}

// readLegacyProfiles reads config.toml from a dedicated viper instance (so it
// does not depend on or mutate global viper state) and returns the profiles it
// contains, including their crawler credentials.
func readLegacyProfiles(configFile string) ([]legacyProfile, error) {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	settings := v.AllSettings()
	profiles := make([]legacyProfile, 0, len(settings))
	for name := range settings {
		var p Profile
		if err := v.UnmarshalKey(name, &p); err != nil {
			return nil, err
		}
		profiles = append(profiles, legacyProfile{
			Name:          name,
			ApplicationID: p.ApplicationID,
			APIKey:        p.APIKey,
			AdminAPIKey:   p.AdminAPIKey,
			SearchHosts:   p.SearchHosts,
			Default:       p.Default,
			CrawlerUserID: v.GetString(name + ".crawler_user_id"),
			CrawlerAPIKey: v.GetString(name + ".crawler_api_key"),
		})
	}

	return profiles, nil
}
