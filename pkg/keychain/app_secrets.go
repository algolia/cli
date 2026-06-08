package keychain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	// service is the OS keychain service the CLI stores secrets under. It is
	// shared with the OAuth token entry (pkg/auth); per-application secrets are
	// namespaced by their user key so they never collide.
	service = "algolia-cli"

	// appSecretsUserPrefix namespaces per-application keychain entries.
	appSecretsUserPrefix = "app:"
)

// AppSecrets holds the secret credentials for a single application, persisted
// as a JSON blob in the OS keychain.
type AppSecrets struct {
	APIKey        string `json:"api_key"`
	CrawlerAPIKey string `json:"crawler_api_key,omitempty"`
}

// appSecretsUser returns the keychain user key for a given application ID.
func appSecretsUser(appID string) string {
	return appSecretsUserPrefix + appID
}

// SaveAppSecrets persists the secrets for an application to the OS keychain.
func SaveAppSecrets(appID string, secrets AppSecrets) error {
	if appID == "" {
		return fmt.Errorf("appID is required")
	}

	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	return keyring.Set(service, appSecretsUser(appID), string(data))
}

// LoadAppSecrets reads an application's secrets from the OS keychain. A missing
// entry is not an error: it returns (nil, nil). Real failures (keychain
// unavailable, malformed data) return an error.
func LoadAppSecrets(appID string) (*AppSecrets, error) {
	if appID == "" {
		return nil, fmt.Errorf("appID is required")
	}

	secret, err := keyring.Get(service, appSecretsUser(appID))
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var secrets AppSecrets
	if err := json.Unmarshal([]byte(secret), &secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}
