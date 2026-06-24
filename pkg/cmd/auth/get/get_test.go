package get

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/test"
)

func TestGet_NotLoggedIn(t *testing.T) {
	keyring.MockInit()
	auth.ClearToken()

	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Equal(t, "you are not logged in — run `algolia auth login` first", err.Error())
}

func TestGet_Expired(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "secret-access",
		CreatedAt:   time.Now().Unix() - 7200,
		ExpiresIn:   3600,
		User: &dashboard.User{
			ID:    42,
			Email: "user@example.com",
			Name:  "Test User",
		},
	}))

	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Equal(t, "your session has expired — run `algolia auth login` again", err.Error())
}

func TestGet_PrintsIdentityWithoutTokens(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "secret-access",
		RefreshToken: "secret-refresh",
		CreatedAt:    time.Now().Unix(),
		ExpiresIn:    3600,
		User: &dashboard.User{
			ID:    42,
			Email: "user@example.com",
			Name:  "Test User",
		},
	}))

	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewGetCmd(f)
	out, err := test.Execute(cmd, "--output ndjson", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `"user_id":"42"`)
	assert.Contains(t, out.String(), `"email":"user@example.com"`)
	assert.Contains(t, out.String(), `"name":"Test User"`)
	assert.NotContains(t, out.String(), "secret-access")
	assert.NotContains(t, out.String(), "secret-refresh")
	assert.NotContains(t, out.String(), "access_token")
	assert.NotContains(t, out.String(), "refresh_token")
}
