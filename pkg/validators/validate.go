package validators

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrAdminAPIKeyNotConfigured is the error returned when the loaded profile is missing the admin_api_key property
	ErrAdminAPIKeyNotConfigured = errors.New("you have not configured your admin API key yet")
	// ErrApplicationIDNotConfigured is the error returned when the loaded profile is missing the application_id property
	ErrApplicationIDNotConfigured = errors.New("you have not configured your Application ID yet")
)

// AdminAPIKey validates that a string looks like an Admin API key.
func AdminAPIKey(input string) error {
	if len(input) == 0 {
		return ErrAdminAPIKeyNotConfigured
	} else if len(input) != 32 {
		return errors.New("the provided Admin API key looks wrong, it must be 32 characters long")
	}
	return nil
}

// PathExists validates that a string is a path that exists.
func PathExists(input string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("the provided path %s does not exist", input)
	}
	return nil
}
