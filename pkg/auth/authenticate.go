package auth

import (
	"errors"
	"fmt"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/iostreams"
)

// EnsureAuthenticated returns a valid access token from the stored session.
// If no valid session exists and the terminal is interactive, it triggers the
// browser-based OAuth login flow automatically.
func EnsureAuthenticated(
	io *iostreams.IOStreams,
	client *dashboard.Client,
) (string, error) {
	accessToken, err := GetValidToken(client)
	if err == nil {
		return accessToken, nil
	}

	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "%s %s\n", cs.WarningIcon(), err)

	return RunInteractiveOAuth(io, client, false, nil)
}

// ReauthenticateIfExpired checks if err is a session-expired error from the API.
// If so, it clears the invalid token and triggers the login flow.
func ReauthenticateIfExpired(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	err error,
) (string, error) {
	if !errors.Is(err, dashboard.ErrSessionExpired) {
		return "", err
	}

	cs := io.ColorScheme()
	ClearToken()
	fmt.Fprintf(io.Out, "%s Session expired.\n", cs.WarningIcon())

	return RunInteractiveOAuth(io, client, false, nil)
}
