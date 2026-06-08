package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/algolia/cli/pkg/utils"
)

// ApplicationState holds the non-secret, per-application data persisted in
// state.toml. Secrets (API keys) live in the OS keychain, not here.
type ApplicationState struct {
	APIKeyUUID string `toml:"api_key_uuid"`
	Alias      string `toml:"alias"`
}

// State is the in-memory representation of state.toml, the new source of truth
// for non-secret CLI configuration.
type State struct {
	CurrentApplicationID string                      `toml:"current_application_id"`
	Applications         map[string]ApplicationState `toml:"applications"`
}

// LoadState reads state.toml from path. A missing file is not an error: it
// returns an empty, ready-to-use State.
func LoadState(path string) (*State, error) {
	state := &State{Applications: map[string]ApplicationState{}}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return state, nil
	}

	if _, err := toml.DecodeFile(path, state); err != nil {
		return nil, err
	}

	if state.Applications == nil {
		state.Applications = map[string]ApplicationState{}
	}

	return state, nil
}

// Save writes the State to path atomically: it encodes to a temporary file in
// the same directory, then renames it over the target so a crash mid-write can
// never leave a half-written state.toml.
func (s *State) Save(path string) error {
	if err := utils.MakePath(path); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".state-*.toml")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op once the rename below succeeds

	if err := toml.NewEncoder(tmp).Encode(s); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o600); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}
