package config

import "os"

// ShouldMigrate reports whether the one-time config.toml → state.toml +
// keychain migration still has to run: a legacy config.toml exists and
// state.toml does not. state.toml is only written when the migration (or any
// new-model write) succeeds, so its absence doubles as the "not migrated yet"
// marker — an aborted migration naturally retries on the next run.
func (c *Config) ShouldMigrate() bool {
	if c.File == "" {
		return false
	}
	if _, err := os.Stat(c.File); err != nil {
		return false
	}
	return !c.StateFileExists()
}

// Migrate moves the legacy config.toml profiles into the new model (state.toml
// + OS keychain). config.toml itself is never modified.
//
// The migration body lands in follow-up PRs; until then this is a no-op that
// deliberately does NOT write state.toml, so ShouldMigrate keeps returning
// true and the real migration will run once shipped.
func (c *Config) Migrate() error {
	return nil
}
