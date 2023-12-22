package config

import (
	"errors"
)

var (
	// ErrAPIKeyNotConfigured is the error returned when the loaded profile is missing the api_key property
	ErrAPIKeyNotConfigured = errors.New("you have not configured your API key yet")
	// ErrApplicationIDNotConfigured is the error returned when the loaded profile is missing the application_id property
	ErrApplicationIDNotConfigured = errors.New("you have not configured your Application ID yet")

	// ErrCrawlerAPIKeyNotConfigured is the error returned when the loaded profile is missing the crawler_api_key property
	ErrCrawlerAPIKeyNotConfigured = errors.New("you have not configured your Crawler API key yet")
	// ErrCrawlerUserIDNotConfigured is the error returned when the loaded profile is missing the crawler_user_id property
	ErrCrawlerUserIDNotConfigured = errors.New("you have not configured your Crawler user ID yet")
)

// AdminAPIKey validates that a string looks like an Admin API key.
func AdminAPIKey(input string) error {
	if len(input) == 0 {
		return ErrAPIKeyNotConfigured
	} else if len(input) != 32 {
		return errors.New("the provided API key looks wrong, it must be 32 characters long")
	}
	return nil
}
