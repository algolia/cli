package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
)

const (
	keyringService = "algolia-cli"
	keyringUser    = "oauth-token"
)

// StoredToken represents the persisted OAuth tokens.
type StoredToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Scope        string `json:"scope,omitempty"`
}

// IsExpired returns true if the access token has expired (with a 60s buffer).
func (t *StoredToken) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt-60
}

// SaveToken persists tokens from an OAuthTokenResponse to the OS keychain.
func SaveToken(resp *dashboard.OAuthTokenResponse) error {
	expiresAt := resp.CreatedAt + int64(resp.ExpiresIn)
	if expiresAt == 0 {
		expiresAt = time.Now().Unix() + int64(resp.ExpiresIn)
	}

	stored := StoredToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
		Scope:        resp.Scope,
	}

	data, err := json.Marshal(stored)
	if err != nil {
		return err
	}

	return keyring.Set(keyringService, keyringUser, string(data))
}

// LoadToken reads the stored token from the OS keychain. Returns nil if not found.
func LoadToken() *StoredToken {
	secret, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil
	}

	var stored StoredToken
	if err := json.Unmarshal([]byte(secret), &stored); err != nil {
		return nil
	}

	if stored.AccessToken == "" {
		return nil
	}

	return &stored
}

// ClearToken removes the stored token from the OS keychain.
func ClearToken() {
	_ = keyring.Delete(keyringService, keyringUser)
}

// GetValidToken returns a valid access token, refreshing if necessary.
func GetValidToken(client *dashboard.Client) (string, error) {
	stored := LoadToken()
	if stored == nil {
		return "", fmt.Errorf("not logged in — run `algolia auth login` first")
	}

	if !stored.IsExpired() {
		return stored.AccessToken, nil
	}

	if stored.RefreshToken == "" {
		ClearToken()
		return "", fmt.Errorf("session expired — run `algolia auth login` to re-authenticate")
	}

	tokenResp, err := client.RefreshToken(stored.RefreshToken)
	if err != nil {
		ClearToken()
		return "", fmt.Errorf("session expired and refresh failed — run `algolia auth login` to re-authenticate: %w", err)
	}

	if err := SaveToken(tokenResp); err != nil {
		return tokenResp.AccessToken, nil
	}

	return tokenResp.AccessToken, nil
}
