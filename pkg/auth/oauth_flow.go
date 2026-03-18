package auth

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/iostreams"
)

// RunOAuth runs the OAuth PKCE flow with a local callback server and returns
// a valid access token. A local HTTP server is started on a random port to
// receive the authorization code via redirect — no copy-paste required.
//
// When openBrowser is true the authorize URL is opened automatically;
// otherwise only the URL is printed (useful when the browser can't be
// launched, e.g. SSH / containers).
//
// If signup is true the browser opens to the sign-up page.
func RunOAuth(io *iostreams.IOStreams, client *dashboard.Client, signup, openBrowser bool) (string, error) {
	cs := io.ColorScheme()

	redirectURI, resultCh, err := StartCallbackServer()
	if err != nil {
		return "", err
	}

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}
	codeChallenge := CodeChallenge(codeVerifier)

	var authorizeURL string
	if signup {
		authorizeURL = client.SignupAuthorizeURL(codeChallenge, redirectURI)
	} else {
		authorizeURL = client.AuthorizeURL(codeChallenge, redirectURI)
	}

	if openBrowser {
		if signup {
			fmt.Fprintf(io.Out, "Opening browser to create an account...\n")
		} else {
			fmt.Fprintf(io.Out, "Opening browser to sign in...\n")
		}
		fmt.Fprintf(io.Out, "If the browser doesn't open, visit:\n  %s\n\n", cs.Bold(authorizeURL))
		_ = OpenBrowser(authorizeURL)
	} else {
		fmt.Fprintf(io.Out, "Open this URL in your browser to authenticate:\n\n  %s\n\n", cs.Bold(authorizeURL))
	}

	fmt.Fprintf(io.Out, "Waiting for authentication...\n")
	cbResult := <-resultCh

	if cbResult.Error != "" {
		return "", fmt.Errorf("authorization failed: %s", cbResult.Error)
	}
	if cbResult.Code == "" {
		return "", fmt.Errorf("no authorization code received")
	}

	io.StartProgressIndicatorWithLabel("Exchanging code for tokens")
	tokenResp, err := client.AuthorizationCodeGrant(cbResult.Code, codeVerifier, redirectURI)
	io.StopProgressIndicator()
	if err != nil {
		return "", err
	}

	if tokenResp.User != nil {
		fmt.Fprintf(io.Out, "%s Signed in as %s\n", cs.SuccessIcon(), cs.Bold(tokenResp.User.Email))
	}

	if err := SaveToken(tokenResp); err != nil {
		fmt.Fprintf(io.ErrOut, "%s Could not save auth token: %s\n", cs.WarningIcon(), err)
	}

	return tokenResp.AccessToken, nil
}

// OpenBrowser opens the given URL in the user's default browser.
func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
