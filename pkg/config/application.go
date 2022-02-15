package config

import (
	"bytes"
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
}

func (a *Application) GetFieldName(field string) string {
	return a.Name + "." + field
}

func (a Application) GetID() (string, error) {
	if os.Getenv("ALGOLIA_APPLICATION_ID") != "" {
		return os.Getenv("ALGOLIA_APPLICATION_ID"), nil
	}

	if a.ID != "" {
		return a.ID, nil
	}

	if err := viper.ReadInConfig(); err == nil {
		return viper.GetString(a.GetFieldName("application_id")), nil
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

	if err := viper.ReadInConfig(); err == nil {
		return viper.GetString(a.GetFieldName("admin_api_key")), nil
	}

	return "", validators.ErrAdminAPIKeyNotConfigured
}

// Add add an application to the configuration
func (a *Application) Add() error {
	runtimeViper := viper.GetViper()
	runtimeViper.Set(a.GetFieldName("application_id"), a.ID)
	runtimeViper.Set(a.GetFieldName("admin_api_key"), a.AdminAPIKey)

	return a.write(viper.GetViper())
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
