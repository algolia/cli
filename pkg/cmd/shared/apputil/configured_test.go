package apputil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/test"
)

func TestApplicationConfigured(t *testing.T) {
	cfg := test.NewDefaultConfigStub()
	cfg.SavedApps = map[string]test.SavedApplication{
		"STATE_APP": {Alias: "prod"},
	}
	// "default" is the config.toml profile's application ID (legacy fallback).
	profileApps := ProfileApplicationIDs(cfg.ConfiguredProfiles())

	t.Run("in state.toml", func(t *testing.T) {
		assert.True(t, ApplicationConfigured(cfg, profileApps, "STATE_APP"))
	})

	t.Run("only in legacy config.toml", func(t *testing.T) {
		assert.True(t, ApplicationConfigured(cfg, profileApps, "default"))
	})

	t.Run("unknown application", func(t *testing.T) {
		assert.False(t, ApplicationConfigured(cfg, profileApps, "UNKNOWN"))
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
