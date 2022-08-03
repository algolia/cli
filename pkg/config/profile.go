package config

import (
	"os"
	"path/filepath"

	"github.com/algolia/cli/pkg/utils"
	"github.com/spf13/viper"
)

type Profile struct {
	Name string

	ApplicationID string `mapstructure:"application_id"`
	AdminAPIKey   string `mapstructure:"admin_api_key"`

	Default bool `mapstructure:"default"`
}

func (p *Profile) GetFieldName(field string) string {
	return p.Name + "." + field
}

func (p *Profile) LoadDefault() {
	configs := viper.AllSettings()
	for appName := range configs {
		if viper.GetBool(appName + ".default") {
			p.Name = appName
		}
	}
}

func (p *Profile) GetApplicationID() (string, error) {
	if os.Getenv("ALGOLIA_APPLICATION_ID") != "" {
		return os.Getenv("ALGOLIA_APPLICATION_ID"), nil
	}

	if p.ApplicationID != "" {
		return p.ApplicationID, nil
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		appId := viper.GetString(p.GetFieldName("application_id"))
		if appId != "" {
			return appId, nil
		}
	}

	return "", ErrApplicationIDNotConfigured
}

func (p *Profile) GetAdminAPIKey() (string, error) {
	if os.Getenv("ALGOLIA_ADMIN_API_KEY") != "" {
		return os.Getenv("ALGOLIA_ADMIN_API_KEY"), nil
	}

	if p.AdminAPIKey != "" {
		return p.AdminAPIKey, nil
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		adminAPIKey := viper.GetString(p.GetFieldName("admin_api_key"))
		if adminAPIKey != "" {
			return adminAPIKey, nil
		}
	}

	return "", ErrAdminAPIKeyNotConfigured
}

// Add add a profile to the configuration
func (p *Profile) Add() error {
	runtimeViper := viper.GetViper()
	runtimeViper.Set(p.GetFieldName("application_id"), p.ApplicationID)
	runtimeViper.Set(p.GetFieldName("admin_api_key"), p.AdminAPIKey)

	err := p.write(runtimeViper)
	if err != nil {
		return err
	}

	return nil
}

// write writes the configuration file
func (p *Profile) write(runtimeViper *viper.Viper) error {
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
