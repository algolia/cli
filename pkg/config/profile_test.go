package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTestViper(configFile string) {
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	_ = viper.ReadInConfig()
}

func TestAddProfile_PreservesExistingProfiles(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.toml")

	// First CLI invocation: add profile A.
	initTestViper(configFile)

	profileA := &Profile{
		Name:          "app-a",
		ApplicationID: "APP_A_ID",
		APIKey:        "key-a",
	}
	require.NoError(t, profileA.Add())

	// Second CLI invocation: add profile B.
	initTestViper(configFile)

	profileB := &Profile{
		Name:          "app-b",
		ApplicationID: "APP_B_ID",
		APIKey:        "key-b",
	}
	require.NoError(t, profileB.Add())

	// Third CLI invocation: read back and verify both profiles exist.
	initTestViper(configFile)

	cfg := &Config{}
	profiles := cfg.ConfiguredProfiles()

	profilesByName := make(map[string]*Profile)
	for _, p := range profiles {
		profilesByName[p.Name] = p
	}

	assert.Len(t, profiles, 2, "both profiles should be preserved on disk")
	assert.Equal(t, "APP_A_ID", profilesByName["app-a"].ApplicationID)
	assert.Equal(t, "APP_B_ID", profilesByName["app-b"].ApplicationID)
}
