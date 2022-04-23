package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config handles all overall configuration for the CLI
type Config struct {
	ApplicationName string

	Application Application

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

// ConfiguredApplications return the applications in the configuration file
func (c *Config) ConfiguredApplications() []*Application {
	configs := viper.AllSettings()
	applications := make([]*Application, 0, len(configs))
	for appName := range configs {
		app := &Application{
			Name: appName,
		}
		if err := viper.UnmarshalKey(appName, app); err != nil {
			log.Fatalf("%s", err)
		}
		applications = append(applications, app)
	}

	return applications
}

// ApplicationNames returns the list of name of the configured applications
func (c *Config) ApplicationNames() []string {
	return viper.AllKeys()
}

// ApplicationExists check if a given application exists
func (c *Config) AppExists(appName string) bool {
	return viper.IsSet(appName)
}
