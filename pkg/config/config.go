package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/config/state"
)

type IConfig interface {
	InitConfig()

	ConfiguredProfiles() []*Profile
	ProfileNames() []string

	ProfileExists(name string) bool
	RemoveProfile(name string) error
	SetDefaultProfile(name string) error

	ApplicationIDExists(appID string) (bool, string)
	ApplicationIDForProfile(profileName string) (bool, string)

	SetCrawlerAuth(profileName, crawlerUserID, crawlerAPIKey string) error

	Profile() *Profile
	Default() *Profile
}

// Config handles all overall configuration for the CLI
type Config struct {
	ApplicationName string

	CurrentProfile Profile

	File string
}

// InitConfig reads in profiles file and ENV variables if set.
func (c *Config) InitConfig() {
	// state.toml is the source of truth for credential resolution; config.toml
	// is only read as a legacy fallback (removed in CLI v2.0).
	c.CurrentProfile.statePath = state.DefaultPath()

	if c.File != "" {
		viper.SetConfigFile(c.File)
	} else {
		configFolder := c.GetConfigFolder(os.Getenv("XDG_CONFIG_HOME"))
		configFile := filepath.Join(configFolder, "config.toml")
		c.File = configFile
		viper.SetConfigType("toml")
		viper.SetConfigFile(configFile)
		viper.SetConfigPermissions(os.FileMode(0o600))

		// Try to change permissions manually, because we used to create files
		// with default permissions (0644)
		err := os.Chmod(configFile, os.FileMode(0o600))
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("%s", err)
		}
	}

	_ = viper.ReadInConfig()
}

// GetConfigFolder retrieves the folder where the configuration file is stored
// It searches for the xdg environment path first and will secondarily
// place it in the home directory
func (c *Config) GetConfigFolder(xdgPath string) string {
	configPath := xdgPath
	if configPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configPath = filepath.Join(home, ".config")
	}

	return filepath.Join(configPath, "algolia")
}

// loadState loads state.toml for read operations. When no state path is
// configured (e.g. a bare Config in unit tests) it returns an empty state so
// reads stay deterministic and never touch a real home directory.
func (c *Config) loadState() *state.State {
	if c.CurrentProfile.statePath == "" {
		return state.New()
	}
	s, err := state.Load(c.CurrentProfile.statePath)
	if err != nil {
		return state.New()
	}
	return s
}

// writeStatePath returns the path used for state writes, defaulting to the
// well-known location when not explicitly configured (always set in production
// by InitConfig).
func (c *Config) writeStatePath() string {
	if c.CurrentProfile.statePath != "" {
		return c.CurrentProfile.statePath
	}
	return state.DefaultPath()
}

// profileName returns the alias for an application, falling back to its ID.
func profileName(app *state.ApplicationState) string {
	if app.Alias != "" {
		return app.Alias
	}
	return app.ApplicationID
}

// resolveApp resolves a profile name to its application state, matching the
// stored alias first and then the application ID.
func (c *Config) resolveApp(st *state.State, name string) *state.ApplicationState {
	if app := st.AppByAlias(name); app != nil {
		return app
	}
	return st.App(name)
}

// ConfiguredProfiles returns the configured applications from state.toml. Only
// non-secret metadata is populated; secrets stay in the keychain and are read
// lazily by callers that need them.
func (c *Config) ConfiguredProfiles() []*Profile {
	s := c.loadState()

	appIDs := make([]string, 0, len(s.Applications))
	for appID := range s.Applications {
		appIDs = append(appIDs, appID)
	}
	sort.Strings(appIDs)

	profiles := make([]*Profile, 0, len(appIDs))
	for _, appID := range appIDs {
		app := s.Applications[appID]
		profiles = append(profiles, &Profile{
			Name:          profileName(app),
			ApplicationID: app.ApplicationID,
			APIKeyUUID:    app.APIKeyUUID,
			SearchHosts:   app.SearchHosts,
			Default:       appID == s.CurrentApplicationID,
			statePath:     c.CurrentProfile.statePath,
		})
	}

	return profiles
}

// Profile returns the current profile
func (c *Config) Profile() *Profile {
	return &c.CurrentProfile
}

// Default returns the default (current) profile
func (c *Config) Default() *Profile {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Default {
			return profile
		}
	}
	return nil
}

// ProfileNames returns the aliases of the configured profiles.
func (c *Config) ProfileNames() []string {
	profiles := c.ConfiguredProfiles()
	names := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	return names
}

// ProfileExists checks whether a profile with the given name (alias or
// application ID) exists.
func (c *Config) ProfileExists(name string) bool {
	return c.resolveApp(c.loadState(), name) != nil
}

// RemoveProfile removes a profile from state.toml and deletes its secrets from
// the keychain.
func (c *Config) RemoveProfile(name string) error {
	path := c.writeStatePath()
	s, err := state.Load(path)
	if err != nil {
		return err
	}

	app := c.resolveApp(s, name)
	if app == nil {
		return fmt.Errorf("profile '%s' not found", name)
	}
	appID := app.ApplicationID

	s.RemoveApp(appID)
	if err := s.Save(path); err != nil {
		return err
	}

	_ = state.DeleteSecret(appID, state.SecretAPIKey)
	_ = state.DeleteSecret(appID, state.SecretCrawlerAPIKey)
	return nil
}

// SetDefaultProfile marks the named profile's application as the current one.
func (c *Config) SetDefaultProfile(name string) error {
	path := c.writeStatePath()
	s, err := state.Load(path)
	if err != nil {
		return err
	}

	app := c.resolveApp(s, name)
	if app == nil {
		return fmt.Errorf("profile '%s' not found", name)
	}

	s.SetCurrent(app.ApplicationID)
	return s.Save(path)
}

// ApplicationIDExists checks whether an application ID is configured.
func (c *Config) ApplicationIDExists(appID string) (bool, string) {
	if app := c.loadState().App(appID); app != nil {
		return true, profileName(app)
	}
	return false, ""
}

// ApplicationIDForProfile returns the application ID for a given profile name.
func (c *Config) ApplicationIDForProfile(name string) (bool, string) {
	if app := c.resolveApp(c.loadState(), name); app != nil {
		return true, app.ApplicationID
	}
	return false, ""
}

// SetCrawlerAuth stores the crawler user ID in state.toml and the crawler API
// key in the keychain for the named profile.
func (c *Config) SetCrawlerAuth(profile, crawlerUserID, crawlerAPIKey string) error {
	path := c.writeStatePath()
	s, err := state.Load(path)
	if err != nil {
		return err
	}

	app := c.resolveApp(s, profile)
	if app == nil {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	app.CrawlerUserID = crawlerUserID
	s.SetApp(app)
	if err := s.Save(path); err != nil {
		return err
	}

	return state.SetSecret(app.ApplicationID, state.SecretCrawlerAPIKey, crawlerAPIKey)
}
