package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/algolia/cli/pkg/utils"
	"github.com/spf13/viper"
)

// DefaultSearchHosts can be set at build time via ldflags, e.g.
// -X github.com/algolia/cli/pkg/config.DefaultSearchHosts=host1,host2
var DefaultSearchHosts string

type Profile struct {
	Name string

	ApplicationID string   `mapstructure:"application_id"`
	APIKey        string   `mapstructure:"api_key"`
	AdminAPIKey   string   `mapstructure:"admin_api_key"`
	SearchHosts   []string `mapstructure:"search_hosts"`

	Default bool `mapstructure:"default"`

	// config back-references the owning Config for new-model (state.toml +
	// keychain) resolution. nil for standalone profiles (e.g. those returned by
	// ConfiguredProfiles), which then resolve from config.toml only.
	config *Config

	// nameFromDefault records that Name was filled by LoadDefault rather than
	// by an explicit --profile flag, so the new-model resolver doesn't let the
	// legacy default profile shadow state.toml's current application.
	nameFromDefault bool
}

func (p *Profile) GetFieldName(field string) string {
	return p.Name + "." + field
}

func (p *Profile) LoadDefault() {
	configs := viper.AllSettings()
	for appName := range configs {
		if viper.GetBool(appName + ".default") {
			p.Name = appName
			p.nameFromDefault = true
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

	// New model: state.toml current/selected application.
	if p.config != nil {
		if appID := p.config.activeApplicationID(); appID != "" {
			return appID, nil
		}
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		appID := viper.GetString(p.GetFieldName("application_id"))
		if appID != "" {
			return appID, nil
		}
	}

	return "", ErrApplicationIDNotConfigured
}

func (p *Profile) GetAPIKey() (string, error) {
	if os.Getenv("ALGOLIA_API_KEY") != "" {
		return os.Getenv("ALGOLIA_API_KEY"), nil
	}

	if p.APIKey != "" {
		return p.APIKey, nil
	}

	// New model: once an application is resolved, its key comes only from that
	// application's keychain entry — never a different profile's config.toml key.
	if p.config != nil {
		if appID := p.config.activeApplicationID(); appID != "" {
			if secrets := p.config.appSecretsFor(appID); secrets != nil && secrets.APIKey != "" {
				return secrets.APIKey, nil
			}
			// The application is set but its key isn't in this machine's
			// keychain (e.g. state.toml synced across machines without it).
			return "", fmt.Errorf(
				"%w %q; run `algolia application select` to store one, or set ALGOLIA_API_KEY",
				ErrAPIKeyMissingFromKeychain,
				appID,
			)
		}
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		apiKey := viper.GetString(p.GetFieldName("api_key"))
		if apiKey != "" {
			return apiKey, nil
		}
	}

	// Fallback on legacy admin API key
	return p.GetAdminAPIKey()
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

	return "", ErrAPIKeyNotConfigured
}

func (p *Profile) GetSearchHosts() []string {
	envHosts := os.Getenv("ALGOLIA_SEARCH_HOSTS")
	if envHosts != "" {
		return strings.Split(envHosts, ",")
	}

	if p.SearchHosts != nil {
		return p.SearchHosts
	}

	// New model: hosts recorded for the resolved application. Empty falls
	// through to the legacy config.toml lookup while both models coexist.
	if p.config != nil {
		if appID := p.config.activeApplicationID(); appID != "" {
			if hosts := p.config.loadState().Applications[appID].SearchHosts; len(hosts) > 0 {
				return hosts
			}
		}
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		hosts := viper.GetStringSlice(p.GetFieldName("search_hosts"))
		if hosts != nil {
			return hosts
		}
	}

	if DefaultSearchHosts != "" {
		return strings.Split(DefaultSearchHosts, ",")
	}

	return nil
}

// GetCrawlerUserID returns the Crawler user ID.
func (p *Profile) GetCrawlerUserID() (string, error) {
	if os.Getenv("ALGOLIA_CRAWLER_USER_ID") != "" {
		return os.Getenv("ALGOLIA_CRAWLER_USER_ID"), nil
	}

	// New model: the user ID recorded for the resolved application. Empty
	// falls through to the legacy config.toml lookup.
	if p.config != nil {
		if appID := p.config.activeApplicationID(); appID != "" {
			if userID := p.config.loadState().Applications[appID].CrawlerUserID; userID != "" {
				return userID, nil
			}
		}
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		userID := viper.GetString(p.GetFieldName("crawler_user_id"))
		if userID != "" {
			return userID, nil
		}
	}

	return "", ErrCrawlerUserIDNotConfigured
}

// GetCrawlerAPIKey returns the Crawler API key.
func (p *Profile) GetCrawlerAPIKey() (string, error) {
	if os.Getenv("ALGOLIA_CRAWLER_API_KEY") != "" {
		return os.Getenv("ALGOLIA_CRAWLER_API_KEY"), nil
	}

	// New model: once an application is resolved, its crawler key comes only from
	// that application's keychain entry — never a different profile's.
	if p.config != nil {
		if appID := p.config.activeApplicationID(); appID != "" {
			if secrets := p.config.appSecretsFor(appID); secrets != nil &&
				secrets.CrawlerAPIKey != "" {
				return secrets.CrawlerAPIKey, nil
			}
			// The application is set but its crawler key isn't in this
			// machine's keychain.
			return "", fmt.Errorf(
				"no Crawler API key stored in your keychain for the current application %q; run `algolia auth crawler` to store one, or set ALGOLIA_CRAWLER_API_KEY",
				appID,
			)
		}
	}

	if p.Name == "" {
		p.LoadDefault()
	}

	if err := viper.ReadInConfig(); err == nil {
		apiKey := viper.GetString(p.GetFieldName("crawler_api_key"))
		if apiKey != "" {
			return apiKey, nil
		}
	}

	return "", ErrCrawlerAPIKeyNotConfigured
}

// Add adds a profile to the configuration, preserving any existing profiles.
func (p *Profile) Add() error {
	runtimeViper := viper.GetViper()
	runtimeViper.Set(p.GetFieldName("application_id"), p.ApplicationID)
	runtimeViper.Set(p.GetFieldName("api_key"), p.APIKey)

	return p.write(runtimeViper)
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
