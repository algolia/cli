package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type Application struct {
	Name string
	ID   string

	AdminAPIKey      string
	UsageAPIKey      string
	MonitoringAPIKey string
}

// AddApp add an application to the configuration
func (a *Application) AddApp() error {
	writeErr := a.writeApp(viper.GetViper())
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// GetField return the configuration field for the current application
func (a *Application) GetField(field string) string {
	return a.Name + "." + field
}

// GetFieldValue return the configuration field value for the current application
func (a *Application) GetFieldValue(field string) string {
	return viper.GetString(a.GetField(field))
}

// writeApp writes the application parameters to the configuration file
func (a *Application) writeApp(runtimeViper *viper.Viper) error {
	profilesFile := viper.ConfigFileUsed()

	err := utils.MakePath(profilesFile)
	if err != nil {
		return err
	}

	if a.ID != "" {
		runtimeViper.Set(a.GetField("id"), strings.TrimSpace(a.ID))
	}

	if a.AdminAPIKey != "" {
		runtimeViper.Set(a.GetField("admin_api_key"), strings.TrimSpace(a.AdminAPIKey))
	}

	if a.UsageAPIKey != "" {
		runtimeViper.Set(a.GetField("usage_api_key"), strings.TrimSpace(a.UsageAPIKey))
	}

	if a.MonitoringAPIKey != "" {
		runtimeViper.Set(a.GetField("monitoring_api_key"), strings.TrimSpace(a.MonitoringAPIKey))
	}

	runtimeViper.SetConfigFile(profilesFile)
	runtimeViper.SetConfigType(filepath.Ext(profilesFile))

	err = runtimeViper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// IsAuthenticated return true if the profile is "authenticated"
// It just checks if the Application ID and admin API key are present
func (a *Application) IsAuthenticated() bool {
	if a.ID != "" && a.AdminAPIKey != "" {
		return true
	}

	return false
}

// Validate return an error if the application credentials are invalid
func (a *Application) Validate() error {
	client := search.NewClient(a.ID, a.AdminAPIKey)
	_, err := client.ListAPIKeys(a.AdminAPIKey)
	if err != nil {
		return errors.New("could not validate the provided application credentials. It can either be a server or a network error or wrong appID/key credentials")
	}
	// TODO: Check if the key is an admin one (created date is not set

	return nil
}

// GetAdminAPIKey return the admin API key for the current application
func (a *Application) GetAdminAPIKey() (string, error) {
	if err := viper.ReadInConfig(); err == nil {
		key := a.GetFieldValue("admin_api_key")

		err := validators.AdminAPIKey(key)
		if err != nil {
			return "", err
		}

		return key, nil
	}

	return "", validators.ErrAdminAPIKeyNotConfigured
}

// GetID return the ID for the current application
func (a *Application) GetID() (string, error) {
	if err := viper.ReadInConfig(); err == nil {
		id := a.GetFieldValue("application_id")

		return id, nil
	}

	return "", validators.ErrApplicationIDNotConfigured
}
