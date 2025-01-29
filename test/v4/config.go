package test

import (
	"github.com/algolia/cli/pkg/config"
)

type ConfigStub struct {
	CurrentProfile config.Profile
	profiles       []*config.Profile
}

func (c *ConfigStub) InitConfig() {}

func (c *ConfigStub) Profile() *config.Profile {
	return &c.CurrentProfile
}

func (c *ConfigStub) Default() *config.Profile {
	return &c.CurrentProfile
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
