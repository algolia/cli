package auth

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
)

const (
	keyringService = "algolia-cli"
	keyringUser    = "oauth-token"
)

// StoredToken represents the persisted OAuth tokens and the identity of the
// authenticated user.
type StoredToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Scope        string `json:"scope,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	Email        string `json:"email,omitempty"`
	Name         string `json:"name,omitempty"`
}

// IsExpired returns true if the access token has expired (with a 60s buffer).
func (t *StoredToken) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt-60
}

// SaveToken persists tokens (and the user identity, when present) from an
// OAuthTokenResponse to the OS keychain.
func SaveToken(resp *dashboard.OAuthTokenResponse) error {
	return persistToken(storedTokenFromResponse(resp))
}

// storedTokenFromResponse builds a StoredToken from an OAuth token response,
// including the user identity when the response carries a user object.
func storedTokenFromResponse(resp *dashboard.OAuthTokenResponse) StoredToken {
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

	if resp.User != nil && resp.User.ID != 0 {
		stored.UserID = strconv.Itoa(resp.User.ID)
		stored.Email = resp.User.Email
		stored.Name = resp.User.Name
	}

	return stored
}

// persistToken marshals and writes a StoredToken to the OS keychain.
func persistToken(stored StoredToken) error {
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

	refreshed := storedTokenFromResponse(tokenResp)
	if refreshed.UserID == "" {
		refreshed.UserID = stored.UserID
		refreshed.Email = stored.Email
		refreshed.Name = stored.Name
	}

	_ = persistToken(refreshed)

	return refreshed.AccessToken, nil
}
