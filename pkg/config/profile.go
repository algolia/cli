package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/algolia/algolia-cli/pkg/validators"
	"github.com/spf13/viper"
)

// Profile handles all things related to managing the project specific configurations
type Profile struct {
	ProfileName      string
	ApplicationID    string
	AdminAPIKey      string
	UsageAPIKey      string
	MonitoringAPIKey string
}

// CreateProfile creates a profile when logging in
func (p *Profile) CreateProfile() error {
	writeErr := p.writeProfile(viper.GetViper())
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// GetConfigField returns the configuration field for the specific profile
func (p *Profile) GetConfigField(field string) string {
	return p.ProfileName + "." + field
}

func (p *Profile) GetFieldValue(field string) string {
	return viper.GetString(p.GetConfigField(field))
}

func (p *Profile) writeProfile(runtimeViper *viper.Viper) error {
	profilesFile := viper.ConfigFileUsed()

	err := makePath(profilesFile)
	if err != nil {
		return err
	}

	if p.ApplicationID != "" {
		runtimeViper.Set(p.GetConfigField("application_id"), strings.TrimSpace(p.ApplicationID))
	}

	if p.AdminAPIKey != "" {
		runtimeViper.Set(p.GetConfigField("admin_api_key"), strings.TrimSpace(p.AdminAPIKey))
	}

	if p.UsageAPIKey != "" {
		runtimeViper.Set(p.GetConfigField("usage_api_key"), strings.TrimSpace(p.UsageAPIKey))
	}

	if p.MonitoringAPIKey != "" {
		runtimeViper.Set(p.GetConfigField("monitoring_api_key"), strings.TrimSpace(p.MonitoringAPIKey))
	}

	runtimeViper.SetConfigFile(profilesFile)

	// Ensure we preserve the config file type
	runtimeViper.SetConfigType(filepath.Ext(profilesFile))

	err = runtimeViper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// IsAuthenticated will return true if the profile has been authenticated
func (p *Profile) IsAuthenticated() bool {
	if p.ApplicationID != "" && p.AdminAPIKey != "" {
		return true
	}

	return false
}

// GetAdminAPIKey will return the existing admin api key for the given profile
func (p *Profile) GetAdminAPIKey() (string, error) {
	envKey := os.Getenv("ALGOLIA_ADMIN_API_KEY")
	if envKey != "" {
		err := validators.ErrAdminAPIKeyNotConfigured
		if err != nil {
			return "", err
		}

		return envKey, nil
	}

	if p.AdminAPIKey != "" {
		err := validators.AdminAPIKey(p.AdminAPIKey)
		if err != nil {
			return "", err
		}

		return p.AdminAPIKey, nil
	}

	// Try to fetch the API key from the configuration file
	if err := viper.ReadInConfig(); err == nil {
		key := viper.GetString(p.GetConfigField("admin_api_key"))

		err := validators.AdminAPIKey(key)
		if err != nil {
			return "", err
		}

		return key, nil
	}

	return "", validators.ErrAdminAPIKeyNotConfigured
}

// GetAppID will return the existing application ID for the given profile
func (p *Profile) GetApplicationID() (string, error) {
	envKey := os.Getenv("ALGOLIA_APP_ID")
	if envKey != "" {
		err := validators.ErrApplicationIDNotConfigured
		if err != nil {
			return "", err
		}

		return envKey, nil
	}

	if p.AdminAPIKey != "" {
		// TODO: Validate the application ID?
		return p.AdminAPIKey, nil
	}

	// Try to fetch the API key from the configuration file
	if err := viper.ReadInConfig(); err == nil {
		key := viper.GetString(p.GetConfigField("application_id"))

		// TODO: Validate the application ID?
		return key, nil
	}

	return "", validators.ErrApplicationIDNotConfigured
}
