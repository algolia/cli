package apputil

import (
	"fmt"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

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

type ConfiguredStatus int

const (
	StatusUnknown ConfiguredStatus = iota
	StatusOutOfSync
	StatusConfigured
)

func ApplicationStatus(cfg config.IConfig, profileApps map[string]bool, appID string) ConfiguredStatus {
	if cfg.ApplicationInState(appID) {
		return StatusConfigured
	}
	if profileApps[appID] {
		return StatusOutOfSync
	}
	return StatusUnknown
}

func AppOptionLabel(
	cfg config.IConfig,
	profileApps map[string]bool,
	cs *iostreams.ColorScheme,
	app dashboard.Application,
) string {
	label := fmt.Sprintf("%s (%s)", app.ID, app.Name)
	switch ApplicationStatus(cfg, profileApps, app.ID) {
	case StatusConfigured:
		label += "  " + cs.Green("(configured)")
	case StatusOutOfSync:
		label += "  " + cs.Yellow("(select to sync)")
	}
	return label
}
