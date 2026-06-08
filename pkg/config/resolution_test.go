package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_LoadStateCachesAndToleratesMissing(t *testing.T) {
	cfg := &Config{StateFile: filepath.Join(t.TempDir(), "state.toml")}

	st := cfg.loadState()
	require.NotNil(t, st)
	assert.Empty(t, st.CurrentApplicationID)

	// Second call returns the same cached pointer (loaded once).
	assert.Same(t, st, cfg.loadState())
}

func TestConfig_LoadStateReadsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(
		path,
		[]byte(
			"current_application_id = \"APP1\"\n\n[applications.APP1]\nalias = \"prod\"\n",
		),
		0o600,
	))

	cfg := &Config{StateFile: path}
	st := cfg.loadState()
	assert.Equal(t, "APP1", st.CurrentApplicationID)
	assert.Equal(t, "prod", st.Applications["APP1"].Alias)
}
