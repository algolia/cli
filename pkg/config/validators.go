package config

import (
	"errors"
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
