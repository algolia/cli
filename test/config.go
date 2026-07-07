package test

import (
	"fmt"

	"github.com/algolia/cli/pkg/config"
)

type CrawlerAuth struct {
	UserID string
	APIKey string
}

// SavedApplication records what SaveApplication stored for an application.
type SavedApplication struct {
	Alias      string
	APIKeyUUID string
	APIKey     string
}

type ConfigStub struct {
	CurrentProfile config.Profile
	profiles       []*config.Profile
	CrawlerAuth    map[string]CrawlerAuth

	ActiveAppID  string
	CurrentAppID string
	SavedApps    map[string]SavedApplication
	CrawlerKeys  map[string]string
	HasStateFile bool
}

func (c *ConfigStub) InitConfig() {}

func (c *ConfigStub) Profile() *config.Profile {
	return &c.CurrentProfile
}

func (c *ConfigStub) Default() *config.Profile {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Default {
			return profile
		}
	}

	return nil
}

func (c *ConfigStub) ConfiguredProfiles() []*config.Profile {
	return c.profiles
}

func (c *ConfigStub) ProfileNames() []string {
	names := make([]string, 0, len(c.ConfiguredProfiles()))
	for _, profile := range c.ConfiguredProfiles() {
		names = append(names, profile.Name)
	}
	return names
}

func (c *ConfigStub) ProfileExists(name string) bool {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Name == name {
			return true
		}
	}
	return false
}

func (c *ConfigStub) ApplicationIDExists(appID string) (bool, string) {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.ApplicationID == appID {
			return true, profile.Name
		}
	}
	return false, ""
}

func (c *ConfigStub) ApplicationIDForProfile(profileName string) (bool, string) {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Name == profileName {
			return true, profile.ApplicationID
		}
	}
	return false, ""
}

func (c *ConfigStub) RemoveProfile(name string) error {
	for i, profile := range c.ConfiguredProfiles() {
		if profile.Name == name {
			c.profiles = append(c.profiles[:i], c.profiles[i+1:]...)
			return nil
		}
	}
	return nil
}

func (c *ConfigStub) SetDefaultProfile(name string) error {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Name == name {
			profile.Default = true
		} else {
			profile.Default = false
		}
	}
	return nil
}

func (c *ConfigStub) SetCrawlerAuth(name, crawlerUserID, crawlerAPIKey string) error {
	for _, profile := range c.ConfiguredProfiles() {
		if profile.Name == name {
			if c.CrawlerAuth == nil {
				c.CrawlerAuth = map[string]CrawlerAuth{}
			}

			c.CrawlerAuth[name] = CrawlerAuth{
				UserID: crawlerUserID,
				APIKey: crawlerAPIKey,
			}

			return nil
		}
	}

	return fmt.Errorf("profile '%s' not found", name)
}

func NewConfigStubWithProfiles(p []*config.Profile) *ConfigStub {
	return &ConfigStub{
		CurrentProfile: *p[0],
		profiles:       p,
	}
}

func NewDefaultConfigStub() *ConfigStub {
	return NewConfigStubWithProfiles([]*config.Profile{
		{
			Name:          "default",
			ApplicationID: "default",
			AdminAPIKey:   "default",
			Default:       true,
		},
	})
}

func (c *ConfigStub) ActiveApplicationID() string {
	return c.ActiveAppID
}

func (c *ConfigStub) ApplicationInState(appID string) bool {
	_, ok := c.SavedApps[appID]
	return ok
}

func (c *ConfigStub) APIKeyUUID(appID string) (string, bool) {
	app, ok := c.SavedApps[appID]
	if !ok || app.APIKeyUUID == "" {
		return "", false
	}
	return app.APIKeyUUID, true
}

func (c *ConfigStub) ApplicationIDByAlias(alias string) (string, bool) {
	for appID, app := range c.SavedApps {
		if app.Alias == alias {
			return appID, true
		}
	}
	return "", false
}

func (c *ConfigStub) SaveApplication(
	appID, alias, apiKeyUUID, apiKey string,
	setCurrent bool,
) error {
	if c.SavedApps == nil {
		c.SavedApps = map[string]SavedApplication{}
	}
	saved := c.SavedApps[appID]
	if alias != "" {
		saved.Alias = alias
	}
	if apiKeyUUID != "" {
		saved.APIKeyUUID = apiKeyUUID
	}
	saved.APIKey = apiKey
	c.SavedApps[appID] = saved
	if setCurrent {
		c.CurrentAppID = appID
	}
	return nil
}

func (c *ConfigStub) StateFileExists() bool {
	return c.HasStateFile
}

func (c *ConfigStub) SetCrawlerAPIKey(appID, crawlerAPIKey string) error {
	if c.CrawlerKeys == nil {
		c.CrawlerKeys = map[string]string{}
	}
	c.CrawlerKeys[appID] = crawlerAPIKey
	return nil
}
