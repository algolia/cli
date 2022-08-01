package validators

import (
	"fmt"
	"os"

	"github.com/algolia/cli/pkg/config"
)

// ProfileNameExists validates that a string is a valid profile name.
func ProfileNameExists(cfg config.IConfig) func(profileName interface{}) error {
	return func(profileName interface{}) error {
		if cfg.ProfileExists(profileName.(string)) {
			return fmt.Errorf("profile '%s' already exists", profileName)
		}
		return nil
	}
}

// ApplicationIDExists validates that a string is a valid Application ID.
func ApplicationIDExists(cfg config.IConfig) func(appID interface{}) error {
	return func(appID interface{}) error {
		appIDExists, profile := cfg.ApplicationIDExists(appID.(string))
		if appIDExists {
			return fmt.Errorf("application ID '%s' already exists in profile '%s'", appID, profile)
		}
		return nil
	}
}

// PathExists validates that a string is a path that exists.
func PathExists(input string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("the provided path %s does not exist", input)
	}
	return nil
}
