package apputil

import "github.com/algolia/cli/pkg/config"

// ApplicationConfigured reports whether an application is already known to the
// CLI. state.toml is the source of truth; the legacy config.toml profiles are
// a fallback while config.toml is still supported (remove once it's gone).
func ApplicationConfigured(cfg config.IConfig, appID string) bool {
	if cfg.ApplicationInState(appID) {
		return true
	}
	exists, _ := cfg.ApplicationIDExists(appID)
	return exists
}
