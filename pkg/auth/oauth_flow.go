package auth

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// OAuthOptions configures the behaviour of RunOAuthFlow.
type OAuthOptions struct {
	// Code, if non-empty, skips the interactive browser+prompt step and
	// exchanges this authorization code directly. The caller must also
	// supply CodeVerifier so the PKCE handshake succeeds.
	Code         string
	CodeVerifier string
}

// RunInteractiveOAuth runs the browser-based OAuth PKCE flow and returns
// a valid access token. This is the authentication-only portion — it does
// not handle application selection or profile setup.
// If signup is true, the browser opens to the sign-up page.
func RunInteractiveOAuth(io *iostreams.IOStreams, client *dashboard.Client, signup bool, opts *OAuthOptions) (string, error) {
	if opts == nil {
		opts = &OAuthOptions{}
	}

	// Fast path: caller already has an authorization code + verifier.
	if opts.Code != "" {
		return exchangeCode(io, client, opts.Code, opts.CodeVerifier)
	}

	if !io.CanPrompt() {
		return "", fmt.Errorf("not logged in — run `algolia auth login` interactively first, or use --print-url and --code for non-interactive mode")
	}

	cs := io.ColorScheme()

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}
	codeChallenge := CodeChallenge(codeVerifier)

	var authorizeURL string
	if signup {
		authorizeURL = client.SignupAuthorizeURL(codeChallenge)
		fmt.Fprintf(io.Out, "Opening browser to create an account...\n")
	} else {
		authorizeURL = client.AuthorizeURL(codeChallenge)
		fmt.Fprintf(io.Out, "Opening browser to sign in...\n")
	}
	fmt.Fprintf(io.Out, "If the browser doesn't open, visit:\n  %s\n\n", cs.Bold(authorizeURL))
	_ = OpenBrowser(authorizeURL)

	var code string
	err = prompt.SurveyAskOne(
		&survey.Input{Message: "Paste the authorization code:"},
		&code,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return "", err
	}
	code = strings.TrimSpace(code)

	return exchangeCode(io, client, code, codeVerifier)
}

// PrintAuthorizeURL generates a PKCE challenge, prints the authorize URL,
// and returns the code verifier so the caller can later exchange the code
// in a separate invocation via --code / --code-verifier.
func PrintAuthorizeURL(io *iostreams.IOStreams, client *dashboard.Client, signup bool) (codeVerifier string, err error) {
	cs := io.ColorScheme()

	codeVerifier, err = GenerateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}
	codeChallenge := CodeChallenge(codeVerifier)

	var authorizeURL string
	if signup {
		authorizeURL = client.SignupAuthorizeURL(codeChallenge)
	} else {
		authorizeURL = client.AuthorizeURL(codeChallenge)
	}

	fmt.Fprintf(io.Out, "Open this URL in your browser to authenticate:\n\n  %s\n\n", cs.Bold(authorizeURL))
	fmt.Fprintf(io.Out, "Then run:\n\n  algolia auth login --code <AUTHORIZATION_CODE> --code-verifier %s\n\n", codeVerifier)
	return codeVerifier, nil
}

func exchangeCode(io *iostreams.IOStreams, client *dashboard.Client, code, codeVerifier string) (string, error) {
	cs := io.ColorScheme()

	io.StartProgressIndicatorWithLabel("Exchanging code for tokens")
	tokenResp, err := client.AuthorizationCodeGrant(code, codeVerifier)
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
