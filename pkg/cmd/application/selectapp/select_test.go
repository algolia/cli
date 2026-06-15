package selectapp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/test"
)

func TestApplicationConfigured(t *testing.T) {
	cfg := test.NewDefaultConfigStub()
	cfg.SavedApps = map[string]test.SavedApplication{
		"STATE_APP": {Alias: "prod"},
	}
	// "default" is the config.toml profile's application ID (legacy fallback).

	t.Run("in state.toml", func(t *testing.T) {
		assert.True(t, applicationConfigured(cfg, "STATE_APP"))
	})

	t.Run("only in legacy config.toml", func(t *testing.T) {
		assert.True(t, applicationConfigured(cfg, "default"))
	})

	t.Run("unknown application", func(t *testing.T) {
		assert.False(t, applicationConfigured(cfg, "UNKNOWN"))
	})
}
