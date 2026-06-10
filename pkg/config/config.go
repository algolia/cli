package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/algolia/cli/pkg/keychain"
	"github.com/algolia/cli/pkg/utils"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// Config handles all overall configuration for the CLI.
//
// It must not be copied after InitConfig: it holds sync primitives (govet
// copylocks) and CurrentProfile holds a back-pointer to it. Pass it by pointer.
type Config struct {
	ApplicationName string

	CurrentProfile Profile

	File      string
	StateFile string

	// state is the parsed state.toml, loaded once per command via loadState.
	stateOnce sync.Once
	state     *State

	// activeApp is the resolved current application ID, computed once.
	activeAppOnce sync.Once
	activeApp     string

	// secretsCache memoizes per-application keychain lookups (guarded by secretsMu).
	secretsMu    sync.Mutex
	secretsCache map[string]*keychain.AppSecrets
}

// InitConfig reads in profiles file and ENV variables if set.
func (c *Config) InitConfig() {
	if c.File != "" {
		viper.SetConfigFile(c.File)
	} else {
		configFolder := c.GetConfigFolder(os.Getenv("XDG_CONFIG_HOME"))
		configFile := filepath.Join(configFolder, "config.toml")
		c.File = configFile
		c.StateFile = filepath.Join(configFolder, "state.toml")
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

	c.CurrentProfile.config = c

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

// loadState reads state.toml once per command and caches it. A missing or
// corrupt file degrades to an empty State so resolution can fall back to
// config.toml rather than crash.
func (c *Config) loadState() *State {
	c.stateOnce.Do(func() {
		st, err := LoadState(c.StateFile)
		if err != nil {
			log.Warnf("ignoring unreadable state file %q: %s", c.StateFile, err)
			st = &State{Applications: map[string]ApplicationState{}}
		}
		c.state = st
	})
	return c.state
}

// activeApplicationID resolves (once per command) which application the new
// model should read against. Returns "" when no new-model app applies, so the
// caller falls back to config.toml.
func (c *Config) activeApplicationID() string {
	c.activeAppOnce.Do(func() {
		c.activeApp = c.resolveActiveApplicationID()
	})
	return c.activeApp
}

func (c *Config) resolveActiveApplicationID() string {
	if v := os.Getenv("ALGOLIA_APPLICATION_ID"); v != "" {
		return v
	}

	p := &c.CurrentProfile
	if p.ApplicationID != "" { // --application-id flag
		return p.ApplicationID
	}

	st := c.loadState()
	// Only a Name set explicitly (--profile flag) counts here: a name filled by
	// LoadDefault must not shadow the state.toml current application.
	if p.Name != "" && !p.nameFromDefault {
		if appID, ok := st.ApplicationByAlias(p.Name); ok {
			return appID
		}
		return "" // unknown alias → let the legacy config.toml profile-by-name handle it
	}

	return st.CurrentApplicationID
}

// appSecretsFor returns the cached keychain secrets for an application, loading
// them once. A missing entry or a keychain failure yields nil (also cached, so
// a single command never hits the keychain twice for the same app). The mutex
// keeps the cache safe if a getter is ever called concurrently; resolution runs
// on the main goroutine today.
func (c *Config) appSecretsFor(appID string) *keychain.AppSecrets {
	c.secretsMu.Lock()
	defer c.secretsMu.Unlock()

	if c.secretsCache == nil {
		c.secretsCache = map[string]*keychain.AppSecrets{}
	}
	if cached, ok := c.secretsCache[appID]; ok {
		return cached
	}

	secrets, err := keychain.LoadAppSecrets(appID)
	if err != nil {
		log.Warnf("ignoring keychain error for application %q: %s", appID, err)
		secrets = nil
	}
	c.secretsCache[appID] = secrets
	return secrets
}

// ConfiguredProfiles return the profiles in the configuration file
func (c *Config) ConfiguredProfiles() []*Profile {
	configs := viper.AllSettings()
	applications := make([]*Profile, 0, len(configs))
	for appName := range configs {
		app := &Profile{
			Name: appName,
		}
		if err := viper.UnmarshalKey(appName, app); err != nil {
			log.Fatalf("%s", err)
		}
		applications = append(applications, app)
	}

	return applications
}

// Profile returns the current profile
func (c *Config) Profile() *Profile {
	return &c.CurrentProfile
}

// Default returns the default profile
func (c *Config) Default() *Profile {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Default {
			return profile
		}
	}
	return nil
}

// ProfileNames returns the list of name of the configured profiles
func (c *Config) ProfileNames() []string {
	return viper.AllKeys()
}

// ProfileExists check if a profile with the given name exists
func (c *Config) ProfileExists(appName string) bool {
	return viper.IsSet(appName)
}

// RemoveProfile remove a profile from the configuration
func (c *Config) RemoveProfile(name string) error {
	runtimeViper := viper.GetViper()
	configMap := runtimeViper.AllSettings()
	delete(configMap, name)

	buf := new(bytes.Buffer)

	encodeErr := toml.NewEncoder(buf).Encode(configMap)
	if encodeErr != nil {
		return encodeErr
	}

	nv := viper.New()
	nv.SetConfigType("toml") // hint to viper that we've encoded the data as toml

	err := nv.ReadConfig(buf)
	if err != nil {
		return err
	}

	return c.write(nv)
}

// SetDefaultProfile set the default profile
func (c *Config) SetDefaultProfile(name string) error {
	configuration, err := c.read()
	if err != nil {
		return err
	}

	configs := configuration.AllSettings()

	found := false

	for profileName := range configs {
		runtimeViper := viper.GetViper()
		runtimeViper.Set(profileName+".default", false)

		if profileName == name {
			found = true
			runtimeViper.Set(profileName+".default", true)
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' not found", name)
	}

	return c.write(configuration)
}

// ApplicationIDExists check if an application ID exists in any profiles
func (c *Config) ApplicationIDExists(appID string) (bool, string) {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.ApplicationID == appID {
			return true, profile.Name
		}
	}

	return false, ""
}

// ApplicationIDForProfile returns the application ID for a given profile name.
func (c *Config) ApplicationIDForProfile(profileName string) (bool, string) {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Name == profileName {
			return true, profile.ApplicationID
		}
	}

	return false, ""
}

// SetCrawlerAuth sets the config properties for crawler public api
func (c *Config) SetCrawlerAuth(profile, crawlerUserID, crawlerAPIKey string) error {
	configuration, err := c.read()
	if err != nil {
		return err
	}

	profiles := configuration.AllSettings()

	if _, exists := profiles[profile]; !exists {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	configuration.Set(profile+".crawler_user_id", crawlerUserID)
	configuration.Set(profile+".crawler_api_key", crawlerAPIKey)

	return c.write(configuration)
}

// read reads the configuration file and returns its runtime
func (c *Config) read() (*viper.Viper, error) {
	runtimeViper := viper.GetViper()

	runtimeViper.SetConfigType("toml")
	err := runtimeViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return runtimeViper, nil
}

// write writes the configuration file
func (c *Config) write(runtimeViper *viper.Viper) error {
	configFile := viper.ConfigFileUsed()
	err := utils.MakePath(configFile)
	if err != nil {
		return err
	}
	runtimeViper.SetConfigFile(configFile)
	runtimeViper.SetConfigType(filepath.Ext(configFile))

	err = runtimeViper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
