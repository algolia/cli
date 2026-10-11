package status

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, v := range []string{
		"ALGOLIA_APPLICATION_ID",
		"ALGOLIA_API_KEY",
		"ALGOLIA_ADMIN_API_KEY",
		"ALGOLIA_SEARCH_HOSTS",
	} {
		t.Setenv(v, "")
	}
}

func newTestOpts(cfg *test.ConfigStub, token *auth.StoredToken) (*StatusOptions, *iostreams.IOStreams, func() string) {
	io, _, out, _ := iostreams.Test()
	opts := &StatusOptions{
		IO:        io,
		Config:    cfg,
		LoadToken: func() *auth.StoredToken { return token },
	}
	return opts, io, out.String
}

func TestStatus_NotSignedInNoCredentials(t *testing.T) {
	clearEnv(t)
	opts, _, out := newTestOpts(&test.ConfigStub{}, nil)

	err := runStatusCmd(opts)

	require.ErrorIs(t, err, cmdutil.ErrSilent)
	assert.Contains(t, out(), "Not signed in")
	assert.Contains(t, out(), "No application selected")
	assert.Contains(t, out(), "algolia auth login")
}

func TestStatus_SignedInWithApplicationAndKey(t *testing.T) {
	clearEnv(t)
	cfg := test.NewDefaultConfigStub()
	cfg.CurrentProfile.ApplicationID = "APP1"
	cfg.CurrentProfile.APIKey = "key-1"
	token := &auth.StoredToken{
		AccessToken: "access",
		ExpiresAt:   time.Now().Unix() + 3600,
		Email:       "user@example.com",
	}
	opts, _, out := newTestOpts(cfg, token)

	err := runStatusCmd(opts)

	require.NoError(t, err)
	assert.Contains(t, out(), "Signed in as user@example.com")
	assert.Contains(t, out(), "Current application: APP1")
	assert.Contains(t, out(), "API key: available")
	assert.NotContains(t, out(), "key-1")
}

func TestStatus_ExpiredSessionWithoutRefreshToken(t *testing.T) {
	clearEnv(t)
	token := &auth.StoredToken{
		AccessToken: "access",
		ExpiresAt:   time.Now().Unix() - 3600,
		Email:       "user@example.com",
	}
	opts, _, out := newTestOpts(&test.ConfigStub{}, token)

	err := runStatusCmd(opts)

	require.ErrorIs(t, err, cmdutil.ErrSilent)
	assert.Contains(t, out(), "Session expired")
}

func TestStatus_EnvCredentialsWithoutSession(t *testing.T) {
	clearEnv(t)
	t.Setenv("ALGOLIA_APPLICATION_ID", "ENVAPP")
	t.Setenv("ALGOLIA_API_KEY", "env-key")
	t.Setenv("ALGOLIA_SEARCH_HOSTS", "c1-test-1.algolianet.com")
	opts, _, out := newTestOpts(&test.ConfigStub{}, nil)

	err := runStatusCmd(opts)

	require.NoError(t, err)
	assert.Contains(t, out(), "Not signed in")
	assert.Contains(t, out(), "Current application: ENVAPP")
	assert.Contains(t, out(), "ALGOLIA_APPLICATION_ID is set")
	assert.Contains(t, out(), "ALGOLIA_API_KEY is set")
	assert.Contains(t, out(), "ALGOLIA_SEARCH_HOSTS is set")
	assert.Contains(t, out(), "c1-test-1.algolianet.com")
	assert.NotContains(t, out(), "env-key")
}

func TestStatus_NeverTriggersLogin(t *testing.T) {
	clearEnv(t)
	keyring.MockInit()
	auth.ClearToken()
	f, out := test.NewFactory(false, nil, &test.ConfigStub{}, "")

	_, err := test.Execute(NewStatusCmd(f), "", out)

	require.Error(t, err)
	assert.NotContains(t, out.String(), "Opening browser")
}

func TestStatus_TokenNeedingRefreshStillSignedIn(t *testing.T) {
	clearEnv(t)
	cfg := test.NewDefaultConfigStub()
	cfg.CurrentProfile.ApplicationID = "APP1"
	cfg.CurrentProfile.APIKey = "key-1"
	token := &auth.StoredToken{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Unix() - 3600,
		Email:        "user@example.com",
	}
	opts, _, out := newTestOpts(cfg, token)

	err := runStatusCmd(opts)

	require.NoError(t, err)
	assert.Contains(t, out(), "Signed in as user@example.com")
	assert.Contains(t, out(), "refreshes automatically")
}
