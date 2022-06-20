package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type Application struct {
	Name string

	ID          string `mapstructure:"application_id"`
	AdminAPIKey string `mapstructure:"admin_api_key"`

	Default bool `mapstructure:"default"`
}

func (a *Application) GetFieldName(field string) string {
	return a.Name + "." + field
}

func (a *Application) LoadDefault() {
	configs := viper.AllSettings()
	for appName := range configs {
		if viper.GetBool(appName + ".default") {
			a.Name = appName
		}
	}
}

func (a Application) GetID() (string, error) {
	if os.Getenv("ALGOLIA_APPLICATION_ID") != "" {
		return os.Getenv("ALGOLIA_APPLICATION_ID"), nil
	}

	if a.ID != "" {
		return a.ID, nil
	}

	if a.Name == "" {
		a.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		appId := viper.GetString(a.GetFieldName("application_id"))
		if appId != "" {
			return appId, nil
		}
	}

	return "", validators.ErrApplicationIDNotConfigured
}

func (a *Application) GetAdminAPIKey() (string, error) {
	if os.Getenv("ALGOLIA_ADMIN_API_KEY") != "" {
		return os.Getenv("ALGOLIA_ADMIN_API_KEY"), nil
	}

	if a.AdminAPIKey != "" {
		return a.AdminAPIKey, nil
	}

	if a.Name == "" {
		a.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		adminAPIKey := viper.GetString(a.GetFieldName("admin_api_key"))
		if adminAPIKey != "" {
			return adminAPIKey, nil
		}
	}

	return "", validators.ErrAdminAPIKeyNotConfigured
}

// Add add an application to the configuration
func (a *Application) Add() error {
	runtimeViper := viper.GetViper()
	runtimeViper.Set(a.GetFieldName("application_id"), a.ID)
	runtimeViper.Set(a.GetFieldName("admin_api_key"), a.AdminAPIKey)

	err := a.write(runtimeViper)
	if err != nil {
		return err
	}

	if a.Default {
		err := a.SetDefault()
		if err != nil {
			return err
		}
	}

	return nil
}

// SetDefault set the default application
func (a *Application) SetDefault() error {
	runtimeViper := viper.GetViper()
	configs := runtimeViper.AllSettings()

	found := false

	for appName := range configs {
		runtimeViper := viper.GetViper()
		runtimeViper.Set(appName+".default", false)

		if appName == a.Name {
			found = true
			runtimeViper.Set(appName+".default", true)
		}
	}

	if !found {
		return fmt.Errorf("application %s not found", a.Name)
	}

	return a.write(runtimeViper)
}

// Remove remove an application from the configuration
func (a *Application) Remove() error {
	runtimeViper := viper.GetViper()
	configMap := runtimeViper.AllSettings()
	delete(configMap, a.Name)

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

	return a.write(nv)
}

// write writes the application parameters to the configuration file
func (a *Application) write(runtimeViper *viper.Viper) error {
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
