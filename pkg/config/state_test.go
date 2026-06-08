package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_LoadMissingFileReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")

	state, err := LoadState(path)
	require.NoError(t, err)
	require.NotNil(t, state)
	assert.Empty(t, state.CurrentApplicationID)
	assert.Empty(t, state.Applications)
}

func TestState_SaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")

	state := &State{
		CurrentApplicationID: "APP1",
		Applications: map[string]ApplicationState{
			"APP1": {APIKeyUUID: "uuid-1", Alias: "prod"},
		},
	}
	require.NoError(t, state.Save(path))

	loaded, err := LoadState(path)
	require.NoError(t, err)
	assert.Equal(t, "APP1", loaded.CurrentApplicationID)
	assert.Equal(t, "uuid-1", loaded.Applications["APP1"].APIKeyUUID)
	assert.Equal(t, "prod", loaded.Applications["APP1"].Alias)
}

func TestState_MutatorsPersist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")

	state := &State{}
	state.UpsertApplication("APP1", ApplicationState{APIKeyUUID: "uuid-1", Alias: "prod"})
	state.SetCurrentApplication("APP1")
	require.NoError(t, state.Save(path))

	loaded, err := LoadState(path)
	require.NoError(t, err)
	assert.Equal(t, "APP1", loaded.CurrentApplicationID)
	assert.Equal(t, "prod", loaded.Applications["APP1"].Alias)

	// Upsert replaces an existing entry rather than duplicating it.
	loaded.UpsertApplication("APP1", ApplicationState{APIKeyUUID: "uuid-2", Alias: "staging"})
	assert.Len(t, loaded.Applications, 1)
	assert.Equal(t, "uuid-2", loaded.Applications["APP1"].APIKeyUUID)
	assert.Equal(t, "staging", loaded.Applications["APP1"].Alias)
}

func TestState_ApplicationByAlias(t *testing.T) {
	state := &State{
		Applications: map[string]ApplicationState{
			"APP1": {Alias: "prod"},
			"APP2": {Alias: "staging"},
		},
	}

	appID, found := state.ApplicationByAlias("staging")
	assert.True(t, found)
	assert.Equal(t, "APP2", appID)

	_, found = state.ApplicationByAlias("missing")
	assert.False(t, found)
}

func TestState_LoadCorruptFileReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	require.NoError(t, os.WriteFile(path, []byte(`key = "unterminated`), 0o600))

	state, err := LoadState(path)
	require.Error(t, err)
	assert.Nil(t, state)
}
