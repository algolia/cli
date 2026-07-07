package apputil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestApplicationStatus(t *testing.T) {
	cfg := test.NewDefaultConfigStub()
	cfg.SavedApps = map[string]test.SavedApplication{
		"STATE_APP": {Alias: "prod"},
	}
	// "default" is the config.toml profile's application ID (legacy fallback).
	profileApps := ProfileApplicationIDs(cfg.ConfiguredProfiles())

	t.Run("in state.toml", func(t *testing.T) {
		assert.Equal(t, StatusConfigured, ApplicationStatus(cfg, profileApps, "STATE_APP"))
	})

	t.Run("only in legacy config.toml", func(t *testing.T) {
		assert.Equal(t, StatusOutOfSync, ApplicationStatus(cfg, profileApps, "default"))
	})

	t.Run("unknown application", func(t *testing.T) {
		assert.Equal(t, StatusUnknown, ApplicationStatus(cfg, profileApps, "UNKNOWN"))
	})
}

func TestAppOptionLabel(t *testing.T) {
	cfg := test.NewDefaultConfigStub()
	cfg.SavedApps = map[string]test.SavedApplication{
		"STATE_APP": {Alias: "prod"},
	}
	profileApps := ProfileApplicationIDs(cfg.ConfiguredProfiles())
	io, _, _, _ := iostreams.Test()
	cs := io.ColorScheme()

	t.Run("configured", func(t *testing.T) {
		app := dashboard.Application{ID: "STATE_APP", Name: "Prod"}
		assert.Equal(t, "STATE_APP (Prod)  (configured)", AppOptionLabel(cfg, profileApps, cs, app))
	})

	t.Run("out of sync", func(t *testing.T) {
		app := dashboard.Application{ID: "default", Name: "Legacy"}
		assert.Equal(t, "default (Legacy)  (select to sync)", AppOptionLabel(cfg, profileApps, cs, app))
	})

	t.Run("unknown", func(t *testing.T) {
		app := dashboard.Application{ID: "NEW_APP", Name: "New"}
		assert.Equal(t, "NEW_APP (New)", AppOptionLabel(cfg, profileApps, cs, app))
	})
}

func TestProfileApplicationIDs(t *testing.T) {
	profiles := []*config.Profile{
		{Name: "prod", ApplicationID: "APP1"},
		{Name: "dev", ApplicationID: "APP2"},
		{Name: "broken", ApplicationID: ""}, // skipped: no app ID
	}

	ids := ProfileApplicationIDs(profiles)

	assert.True(t, ids["APP1"])
	assert.True(t, ids["APP2"])
	assert.False(t, ids[""]) // empty IDs never become a member
	assert.Len(t, ids, 2)
}
