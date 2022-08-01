package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
	if c.File != "" {
		viper.SetConfigFile(c.File)
	} else {
		configFolder := c.GetConfigFolder(os.Getenv("XDG_CONFIG_HOME"))
		configFile := filepath.Join(configFolder, "config.toml")
		c.File = configFile
		viper.SetConfigType("toml")
		viper.SetConfigFile(configFile)
		viper.SetConfigPermissions(os.FileMode(0600))

		// Try to change permissions manually, because we used to create files
		// with default permissions (0644)
		err := os.Chmod(configFile, os.FileMode(0600))
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
	runtimeViper := viper.GetViper()
	configs := runtimeViper.AllSettings()

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

	return c.write(runtimeViper)
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
