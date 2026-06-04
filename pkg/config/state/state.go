// Package state stores the CLI's non-secret application state in state.toml,
// alongside the secrets kept in the OS keychain (see secret.go). It is the
// source of truth introduced to replace the legacy config.toml; config.toml is
// only read as a migration/fallback source and is scheduled for removal in
// CLI v2.0.
package state

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

// ApplicationState holds the non-secret state for a single Algolia
// application. Secrets (api_key, crawler_api_key) live in the OS keychain and
// are intentionally absent here.
type ApplicationState struct {
	// ApplicationID is the Algolia application ID and the map key under
	// [applications.<id>].
	ApplicationID string `toml:"application_id"`
	// Alias preserves the legacy profile name so that `--profile <name>`
	// keeps working until CLI v2.0.
	Alias string `toml:"alias,omitempty"`
	// APIKeyUUID identifies the API key stored in the keychain, so it can be
	// rotated/revoked without re-reading the secret.
	APIKeyUUID string `toml:"api_key_uuid,omitempty"`

	// The following are lazily fetched when needed and cached here.
	ApplicationName string   `toml:"application_name,omitempty"`
	CrawlerUserID   string   `toml:"crawler_user_id,omitempty"`
	Region          string   `toml:"region,omitempty"`
	SearchHosts     []string `toml:"search_hosts,omitempty"`
}

// State is the on-disk representation of state.toml.
type State struct {
	// CurrentApplicationID points at the active application in Applications.
	CurrentApplicationID string `toml:"current_application_id,omitempty"`
	// Applications is keyed by application ID.
	Applications map[string]*ApplicationState `toml:"applications,omitempty"`
}

// New returns an empty, ready-to-use State.
func New() *State {
	return &State{Applications: map[string]*ApplicationState{}}
}

// DefaultPath returns the default location of state.toml, mirroring the config
// folder resolution used for config.toml (XDG_CONFIG_HOME, then ~/.config).
func DefaultPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := homedir.Dir()
		if err != nil {
			// Last resort: relative path in the working directory.
			return filepath.Join("algolia", "state.toml")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "algolia", "state.toml")
}

// Load reads state.toml from path. A missing file yields an empty State and no
// error, so callers can treat "not migrated yet" the same as "empty state".
func Load(path string) (*State, error) {
	s := New()
	if path == "" {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}

	if _, err := toml.Decode(string(data), s); err != nil {
		return nil, err
	}
	if s.Applications == nil {
		s.Applications = map[string]*ApplicationState{}
	}
	return s, nil
}

// Save writes state.toml atomically with 0600 permissions, creating parent
// directories as needed. It writes to a temp file in the destination directory
// and renames it into place to avoid leaving a half-written file.
func (s *State) Save(path string) error {
	if path == "" {
		return errors.New("state: empty path")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(s); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".state-*.toml")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}

// ResolveCurrent returns the current application's state, or nil when no
// current application is set or it is missing from Applications.
func (s *State) ResolveCurrent() *ApplicationState {
	if s == nil || s.CurrentApplicationID == "" {
		return nil
	}
	return s.App(s.CurrentApplicationID)
}

// SetCurrent marks appID as the current application, creating its entry when
// missing.
func (s *State) SetCurrent(appID string) {
	if appID == "" {
		return
	}
	s.ensureApp(appID)
	s.CurrentApplicationID = appID
}

// App returns the application state for appID, or nil when absent.
func (s *State) App(appID string) *ApplicationState {
	if s == nil || s.Applications == nil {
		return nil
	}
	return s.Applications[appID]
}

// AppByAlias returns the application state whose alias matches, or nil. It is
// used to resolve the deprecated `--profile <name>` flag.
func (s *State) AppByAlias(alias string) *ApplicationState {
	if s == nil || alias == "" {
		return nil
	}
	for _, app := range s.Applications {
		if app.Alias == alias {
			return app
		}
	}
	return nil
}

// SetApp inserts or replaces the state for an application (keyed by its ID).
func (s *State) SetApp(app *ApplicationState) {
	if app == nil || app.ApplicationID == "" {
		return
	}
	if s.Applications == nil {
		s.Applications = map[string]*ApplicationState{}
	}
	s.Applications[app.ApplicationID] = app
}

// PutAPIKeyUUID records the keychain API key UUID for appID, creating the
// entry when missing.
func (s *State) PutAPIKeyUUID(appID, uuid string) {
	if appID == "" {
		return
	}
	s.ensureApp(appID).APIKeyUUID = uuid
}

// APIKeyUUID returns the stored API key UUID for appID, or "".
func (s *State) APIKeyUUID(appID string) string {
	if app := s.App(appID); app != nil {
		return app.APIKeyUUID
	}
	return ""
}

// ensureApp returns the existing entry for appID or creates a new one.
func (s *State) ensureApp(appID string) *ApplicationState {
	if s.Applications == nil {
		s.Applications = map[string]*ApplicationState{}
	}
	app, ok := s.Applications[appID]
	if !ok {
		app = &ApplicationState{ApplicationID: appID}
		s.Applications[appID] = app
	}
	return app
}
