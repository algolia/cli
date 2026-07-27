package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/iostreams"
)

func TestEnsureAuthenticated_WithoutATerminalDoesNotStartTheOAuthFlow(t *testing.T) {
	keyring.MockInit()
	ClearToken()

	io, _, stdout, stderr := iostreams.Test()
	require.False(t, io.CanPrompt())

	_, err := EnsureAuthenticated(io, dashboard.NewClient("test"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a terminal")
	assert.Contains(t, err.Error(), "algolia auth login")

	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "not logged in")
}

func TestReauthenticateIfExpired_WithoutATerminalDoesNotStartTheOAuthFlow(t *testing.T) {
	keyring.MockInit()

	io, _, stdout, stderr := iostreams.Test()

	_, err := ReauthenticateIfExpired(io, dashboard.NewClient("test"), dashboard.ErrSessionExpired)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a terminal")
	assert.Contains(t, err.Error(), "algolia auth login")

	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "Session expired")
}

func TestReauthenticateIfExpired_PassesThroughOtherErrors(t *testing.T) {
	io, _, stdout, stderr := iostreams.Test()

	other := assert.AnError
	_, err := ReauthenticateIfExpired(io, dashboard.NewClient("test"), other)
	assert.Same(t, other, err)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}
