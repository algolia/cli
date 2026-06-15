package apputil

import "github.com/algolia/cli/pkg/config"

// ProfileApplicationIDs returns the set of application IDs backed by a legacy
// config.toml profile. Built once by the caller so a per-application loop tests
// membership in O(1) instead of re-parsing config.toml on every iteration.
func ProfileApplicationIDs(profiles []*config.Profile) map[string]bool {
	ids := make(map[string]bool, len(profiles))
	for _, p := range profiles {
		if p.ApplicationID != "" {
			ids[p.ApplicationID] = true
		}
	}
	return ids
}

// ApplicationConfigured reports whether an application is already known to the
// CLI. state.toml is the source of truth (an O(1) cached lookup); profileApps
// is the legacy config.toml fallback while config.toml is still supported
// (remove once it's gone).
func ApplicationConfigured(cfg config.IConfig, profileApps map[string]bool, appID string) bool {
	return cfg.ApplicationInState(appID) || profileApps[appID]
}
