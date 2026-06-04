package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_MissingFileReturnsEmptyState(t *testing.T) {
	s, err := Load(filepath.Join(t.TempDir(), "does-not-exist.toml"))
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.Empty(t, s.CurrentApplicationID)
	assert.Empty(t, s.Applications)
}

func TestLoad_EmptyPathReturnsEmptyState(t *testing.T) {
	s, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.NotNil(t, s.Applications)
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")

	s := New()
	s.SetApp(&ApplicationState{
		ApplicationID:   "APP1",
		Alias:           "prod",
		APIKeyUUID:      "uuid-1",
		ApplicationName: "Production",
		CrawlerUserID:   "crawler-1",
		Region:          "us",
		SearchHosts:     []string{"host1", "host2"},
	})
	s.SetCurrent("APP1")
	require.NoError(t, s.Save(path))

	got, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "APP1", got.CurrentApplicationID)

	app := got.App("APP1")
	require.NotNil(t, app)
	assert.Equal(t, "APP1", app.ApplicationID)
	assert.Equal(t, "prod", app.Alias)
	assert.Equal(t, "uuid-1", app.APIKeyUUID)
	assert.Equal(t, "Production", app.ApplicationName)
	assert.Equal(t, "crawler-1", app.CrawlerUserID)
	assert.Equal(t, "us", app.Region)
	assert.Equal(t, []string{"host1", "host2"}, app.SearchHosts)
}

func TestSave_FilePermissionsAre0600(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.toml")
	s := New()
	s.SetCurrent("APP1")
	require.NoError(t, s.Save(path))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestSave_CreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "deeper", "state.toml")
	s := New()
	s.SetCurrent("APP1")
	require.NoError(t, s.Save(path))

	_, err := os.Stat(path)
	require.NoError(t, err)
}

func TestSave_OverwritesAtomically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.toml")

	first := New()
	first.SetApp(&ApplicationState{ApplicationID: "APP1", Alias: "one"})
	first.SetCurrent("APP1")
	require.NoError(t, first.Save(path))

	second := New()
	second.SetApp(&ApplicationState{ApplicationID: "APP2", Alias: "two"})
	second.SetCurrent("APP2")
	require.NoError(t, second.Save(path))

	got, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "APP2", got.CurrentApplicationID)
	assert.Nil(t, got.App("APP1"))
	assert.NotNil(t, got.App("APP2"))

	// No leftover temp files.
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestSave_EmptyPathErrors(t *testing.T) {
	require.Error(t, New().Save(""))
}

func TestResolveCurrent(t *testing.T) {
	s := New()
	assert.Nil(t, s.ResolveCurrent(), "no current set")

	s.SetApp(&ApplicationState{ApplicationID: "APP1", Alias: "a"})
	s.CurrentApplicationID = "MISSING"
	assert.Nil(t, s.ResolveCurrent(), "current points at missing app")

	s.SetCurrent("APP1")
	require.NotNil(t, s.ResolveCurrent())
	assert.Equal(t, "APP1", s.ResolveCurrent().ApplicationID)
}

func TestSetCurrent_CreatesEntry(t *testing.T) {
	s := New()
	s.SetCurrent("APP1")
	assert.Equal(t, "APP1", s.CurrentApplicationID)
	require.NotNil(t, s.App("APP1"))
	assert.Equal(t, "APP1", s.App("APP1").ApplicationID)
}

func TestAppByAlias(t *testing.T) {
	s := New()
	s.SetApp(&ApplicationState{ApplicationID: "APP1", Alias: "prod"})
	s.SetApp(&ApplicationState{ApplicationID: "APP2", Alias: "staging"})

	require.NotNil(t, s.AppByAlias("staging"))
	assert.Equal(t, "APP2", s.AppByAlias("staging").ApplicationID)
	assert.Nil(t, s.AppByAlias("unknown"))
	assert.Nil(t, s.AppByAlias(""), "empty alias must never match")
}

func TestAPIKeyUUID(t *testing.T) {
	s := New()
	assert.Empty(t, s.APIKeyUUID("APP1"))

	s.PutAPIKeyUUID("APP1", "uuid-123")
	assert.Equal(t, "uuid-123", s.APIKeyUUID("APP1"))

	// PutAPIKeyUUID must not clobber other fields on an existing app.
	s.SetApp(&ApplicationState{ApplicationID: "APP2", Alias: "two"})
	s.PutAPIKeyUUID("APP2", "uuid-456")
	assert.Equal(t, "two", s.App("APP2").Alias)
	assert.Equal(t, "uuid-456", s.App("APP2").APIKeyUUID)
}

func TestSetApp_IgnoresInvalid(t *testing.T) {
	s := New()
	s.SetApp(nil)
	s.SetApp(&ApplicationState{ApplicationID: ""})
	assert.Empty(t, s.Applications)
}
