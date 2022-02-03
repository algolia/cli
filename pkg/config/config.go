package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/algolia/cli/pkg/utils"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config handles all overall configuration for the CLI
type Config struct {
	ApplicationName string
	Applications    map[string]*Application

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

	// If a profiles file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&c.Applications)
	}
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

// GetApplications return the applications in the configuration file
func (c *Config) GetApplications() map[string]string {
	configs := viper.AllSettings()
	applications := make(map[string]string)
	for app := range configs {
		applications[app] = viper.GetStringMapString(app)["application_id"]
	}

	return applications
}

// ApplicationNames returns the list of application names
func (c *Config) ApplicationNames() []string {
	names := make([]string, 0, len(c.Applications))
	for name := range c.Applications {
		names = append(names, name)
	}
	return names
}

// GetCurrentApplication returns the current application
func (c *Config) GetCurrentApplication() (*Application, error) {
	if c.ApplicationName == "" {
		return nil, fmt.Errorf("no application name set")
	}
	return c.Applications[c.ApplicationName], nil
}

// ApplicationExists check if a given application exists
func (c *Config) AppExists(appName string) bool {
	return viper.IsSet(appName)
}

// GetApplicationField returns the configuration field for the specific application
func (c *Config) GetApplicationField(app *Application, field string) string {
	return app.Name + "." + field
}

// AddApplication add an application to the configuration
func (c *Config) AddApplication(app *Application) error {
	runtimeViper := viper.GetViper()
	runtimeViper.Set(c.GetApplicationField(app, "application_id"), app.ID)
	runtimeViper.Set(c.GetApplicationField(app, "admin_api_key"), app.AdminAPIKey)

	return c.write(viper.GetViper())
}

// writeApp writes the application parameters to the configuration file
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
