package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/algolia/cli/pkg/config/state"
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

	// APIKeyUUID is the resource ID of the API key stored in the keychain.
	// It is persisted to state.toml (not config.toml) so the key can later
	// be rotated/revoked without re-reading the secret.
	APIKeyUUID string `mapstructure:"-"`

	Default bool `mapstructure:"default"`

	// statePath points at state.toml; it is set by Config.InitConfig. When
	// empty (e.g. in unit tests) it falls back to state.DefaultPath() for
	// writes and disables state-based reads.
	statePath string
}

// loadState resolves the application state that applies to this profile from
// state.toml. The deprecated `--profile <name>` flag is resolved against the
// stored alias; when the name is not a known alias the current application is
// used instead for consistency. Returns nil when there is no state.toml or no
// applicable application entry.
func (p *Profile) loadState() *state.ApplicationState {
	if p.statePath == "" {
		return nil
	}
	s, err := state.Load(p.statePath)
	if err != nil {
		return nil
	}
	if p.Name != "" {
		if app := s.AppByAlias(p.Name); app != nil {
			return app
		}
	}
	return s.ResolveCurrent()
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

	if app := p.loadState(); app != nil && app.ApplicationID != "" {
		return app.ApplicationID, nil
	}

	// Legacy fallback: config.toml (read-only, scheduled for removal in v2.0).
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

	if app := p.loadState(); app != nil && app.ApplicationID != "" {
		if key, err := state.GetSecret(app.ApplicationID, state.SecretAPIKey); err == nil &&
			key != "" {
			return key, nil
		}
	}

	// Legacy fallback: config.toml (read-only, scheduled for removal in v2.0).
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

	if app := p.loadState(); app != nil && len(app.SearchHosts) > 0 {
		return app.SearchHosts
	}

	// Legacy fallback: config.toml (read-only, scheduled for removal in v2.0).
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

	if app := p.loadState(); app != nil && app.CrawlerUserID != "" {
		return app.CrawlerUserID, nil
	}

	// Legacy fallback: config.toml (read-only, scheduled for removal in v2.0).
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

	if app := p.loadState(); app != nil && app.ApplicationID != "" {
		if key, err := state.GetSecret(app.ApplicationID, state.SecretCrawlerAPIKey); err == nil &&
			key != "" {
			return key, nil
		}
	}

	// Legacy fallback: config.toml (read-only, scheduled for removal in v2.0).
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

// Add persists the profile to the new store: non-secret metadata (application
// ID, alias, search hosts, API key UUID) goes to state.toml and the API key
// itself goes to the OS keychain. config.toml is no longer written.
//
// Existing non-secret fields for the application (e.g. crawler_user_id) are
// preserved. The profile becomes the current application when it is flagged as
// default, or when no application is current yet (so the first configured
// application is immediately usable).
func (p *Profile) Add() error {
	if p.ApplicationID == "" {
		return ErrApplicationIDNotConfigured
	}

	path := p.statePath
	if path == "" {
		path = state.DefaultPath()
	}

	s, err := state.Load(path)
	if err != nil {
		return err
	}

	app := s.App(p.ApplicationID)
	if app == nil {
		app = &state.ApplicationState{ApplicationID: p.ApplicationID}
	}
	if p.Name != "" {
		app.Alias = p.Name
	}
	if len(p.SearchHosts) > 0 {
		app.SearchHosts = p.SearchHosts
	}
	if p.APIKeyUUID != "" {
		app.APIKeyUUID = p.APIKeyUUID
	}
	s.SetApp(app)

	if p.Default || s.CurrentApplicationID == "" {
		s.SetCurrent(p.ApplicationID)
	}

	if err := s.Save(path); err != nil {
		return err
	}

	if p.APIKey != "" {
		return state.SetSecret(p.ApplicationID, state.SecretAPIKey, p.APIKey)
	}

	return nil
}
