package state

import (
	"errors"

	"github.com/zalando/go-keyring"
)

// secretService is the keychain service under which per-application secrets are
// stored. It matches the service used for the OAuth token (pkg/auth) but the
// account names are namespaced per application and secret kind, so there is no
// collision.
const secretService = "algolia-cli"

// Secret kinds stored in the OS keychain, namespaced per application.
const (
	SecretAPIKey        = "api_key"
	SecretCrawlerAPIKey = "crawler_api_key"
)

// secretAccount builds the keychain account name for an application secret,
// e.g. "ABCDEF1234:api_key".
func secretAccount(appID, kind string) string {
	return appID + ":" + kind
}

// GetSecret reads a secret of the given kind for appID from the OS keychain.
// It returns an empty string and no error when the secret is not set.
func GetSecret(appID, kind string) (string, error) {
	value, err := keyring.Get(secretService, secretAccount(appID, kind))
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// SetSecret stores a secret of the given kind for appID in the OS keychain.
func SetSecret(appID, kind, value string) error {
	return keyring.Set(secretService, secretAccount(appID, kind), value)
}

// DeleteSecret removes a secret of the given kind for appID from the OS
// keychain. Removing a missing secret is not an error.
func DeleteSecret(appID, kind string) error {
	err := keyring.Delete(secretService, secretAccount(appID, kind))
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}
	return nil
}
