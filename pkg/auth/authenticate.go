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
	fmt.Fprintf(io.ErrOut, "%s %s\n", cs.WarningIcon(), err)

	if !io.CanPrompt() {
		return "", fmt.Errorf(
			"no usable session and authenticating requires a terminal: run %s",
			cs.Bold("algolia auth login"),
		)
	}

	// No flow tracker: this re-authentication belongs to the calling flow,
	// not to an `auth login` funnel.
	return RunOAuth(io, client, false, true, nil)
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
	fmt.Fprintf(io.ErrOut, "%s Session expired.\n", cs.WarningIcon())

	if !io.CanPrompt() {
		return "", fmt.Errorf(
			"your session expired and authenticating requires a terminal: run %s",
			cs.Bold("algolia auth login"),
		)
	}

	// No flow tracker: this re-authentication belongs to the calling flow,
	// not to an `auth login` funnel.
	return RunOAuth(io, client, false, true, nil)
}
