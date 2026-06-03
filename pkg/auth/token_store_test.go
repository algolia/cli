package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
)

func TestSaveToken_PersistsUserIdentity(t *testing.T) {
	keyring.MockInit()

	resp := &dashboard.OAuthTokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    7200,
		CreatedAt:    time.Now().Unix(),
		Scope:        "scope:test",
		User: &dashboard.User{
			ID:    42,
			Email: "user@test.com",
			Name:  "Test User",
		},
	}

	require.NoError(t, SaveToken(resp))

	stored := LoadToken()
	require.NotNil(t, stored)
	assert.Equal(t, "access-token", stored.AccessToken)
	assert.Equal(t, "42", stored.UserID)
	assert.Equal(t, "user@test.com", stored.Email)
	assert.Equal(t, "Test User", stored.Name)
}

func TestSaveToken_WithoutUserLeavesIdentityEmpty(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "access-token",
		ExpiresIn:   7200,
		CreatedAt:   time.Now().Unix(),
	}))

	stored := LoadToken()
	require.NotNil(t, stored)
	assert.Empty(t, stored.UserID)
	assert.Empty(t, stored.Email)
	assert.Empty(t, stored.Name)
}

// newRefreshServer returns a dashboard client whose token endpoint replies with
// the given response, mimicking POST /2/oauth/token for grant_type=refresh_token.
func newRefreshServer(t *testing.T, resp dashboard.OAuthTokenResponse) *dashboard.Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(resp))
	}))
	t.Cleanup(srv.Close)

	client := dashboard.NewClientWithHTTPClient("test-client-id", srv.Client())
	client.DashboardURL = srv.URL
	return client
}

func TestGetValidToken_PreservesIdentityWhenRefreshOmitsUser(t *testing.T) {
	keyring.MockInit()

	// An existing session that already knows its user, but whose access token
	// has expired so a refresh is required.
	require.NoError(t, persistToken(StoredToken{
		AccessToken:  "old-access",
		RefreshToken: "refresh-1",
		ExpiresAt:    time.Now().Unix() - 60,
		UserID:       "42",
		Email:        "user@test.com",
		Name:         "Test User",
	}))

	client := newRefreshServer(t, dashboard.OAuthTokenResponse{
		AccessToken:  "new-access",
		RefreshToken: "refresh-2",
		ExpiresIn:    7200,
		CreatedAt:    time.Now().Unix(),
	})

	token, err := GetValidToken(client)
	require.NoError(t, err)
	assert.Equal(t, "new-access", token)

	stored := LoadToken()
	require.NotNil(t, stored)
	assert.Equal(t, "new-access", stored.AccessToken)
	assert.Equal(t, "refresh-2", stored.RefreshToken)
	// Identity survives a refresh response that doesn't echo the user back.
	assert.Equal(t, "42", stored.UserID)
	assert.Equal(t, "user@test.com", stored.Email)
	assert.Equal(t, "Test User", stored.Name)
}

func TestGetValidToken_SelfHealsIdentityFromRefresh(t *testing.T) {
	keyring.MockInit()

	// A session created before identity was persisted: no user fields yet.
	require.NoError(t, persistToken(StoredToken{
		AccessToken:  "old-access",
		RefreshToken: "refresh-1",
		ExpiresAt:    time.Now().Unix() - 60,
	}))

	client := newRefreshServer(t, dashboard.OAuthTokenResponse{
		AccessToken:  "new-access",
		RefreshToken: "refresh-2",
		ExpiresIn:    7200,
		CreatedAt:    time.Now().Unix(),
		User: &dashboard.User{
			ID:    7,
			Email: "healed@test.com",
			Name:  "Healed User",
		},
	})

	_, err := GetValidToken(client)
	require.NoError(t, err)

	stored := LoadToken()
	require.NotNil(t, stored)
	// Identity is back-filled from the refresh response.
	assert.Equal(t, "7", stored.UserID)
	assert.Equal(t, "healed@test.com", stored.Email)
	assert.Equal(t, "Healed User", stored.Name)
}
