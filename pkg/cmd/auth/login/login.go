package login

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// DefaultOAuthClientID is injected at build time via ldflags.
// Override with ALGOLIA_OAUTH_CLIENT_ID environment variable for local development.
var DefaultOAuthClientID = ""

// OAuthClientID returns the OAuth client ID, preferring the ALGOLIA_OAUTH_CLIENT_ID
// environment variable over the compiled-in default (set via ldflags).
func OAuthClientID() string {
	if v := os.Getenv("ALGOLIA_OAUTH_CLIENT_ID"); v != "" {
		return v
	}
	if DefaultOAuthClientID == "" {
		fmt.Fprintln(os.Stderr, "fatal: ALGOLIA_OAUTH_CLIENT_ID is not set and no default was compiled in")
		os.Exit(1)
	}
	return DefaultOAuthClientID
}

// LoginOptions holds all options for the login command.
type LoginOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	AppName     string
	ProfileName string
	Default     bool

	// Non-interactive OAuth fields
	PrintURL     bool
	Code         string
	CodeVerifier string

	NewDashboardClient func(clientID string) *dashboard.Client
}

// NewLoginCmd returns a new instance of the login command.
func NewLoginCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Sign in to your Algolia account",
		Long: heredoc.Doc(`
			Authenticate with your Algolia account via the browser.
			Opens the Algolia Dashboard for sign-in (or sign-up), then exchanges
			the authorization code for API tokens using OAuth 2.0 with PKCE.

			For non-interactive environments (CI, sandboxed terminals), use the
			two-step flow:

			  1. algolia auth login --print-url
			     → prints the authorize URL and a PKCE code-verifier
			  2. algolia auth login --code <CODE> --code-verifier <VERIFIER>
			     → exchanges the code and sets up your profile
		`),
		Example: heredoc.Doc(`
			# Sign in interactively (opens browser)
			$ algolia auth login

			# Auto-select an application by name
			$ algolia auth login --app-name "My App" --default

			# Non-interactive: step 1 — get the URL
			$ algolia auth login --print-url

			# Non-interactive: step 2 — exchange the code
			$ algolia auth login --code <AUTH_CODE> --code-verifier <VERIFIER>
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Code != "" && opts.CodeVerifier == "" {
				return fmt.Errorf("--code-verifier is required when using --code")
			}
			return runLoginCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.AppName, "app-name", "", "Auto-select application by name")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Name for the CLI profile (defaults to application name)")
	cmd.Flags().BoolVar(&opts.Default, "default", true, "Set the profile as the default")
	cmd.Flags().BoolVar(&opts.PrintURL, "print-url", false, "Print the authorize URL and PKCE verifier, then exit (for non-interactive flows)")
	cmd.Flags().StringVar(&opts.Code, "code", "", "Authorization code obtained from the authorize URL")
	cmd.Flags().StringVar(&opts.CodeVerifier, "code-verifier", "", "PKCE code verifier from --print-url (required with --code)")

	return cmd
}

func runLoginCmd(opts *LoginOptions) error {
	return RunOAuthFlow(opts, false)
}

// RunOAuthFlow runs the full browser-based OAuth + profile setup flow.
// If signup is true, the browser opens to the sign-up page instead of sign-in.
func RunOAuthFlow(opts *LoginOptions, signup bool) error {
	cs := opts.IO.ColorScheme()
	client := opts.NewDashboardClient(OAuthClientID())

	if opts.PrintURL {
		_, err := auth.PrintAuthorizeURL(opts.IO, client, signup)
		return err
	}

	oauthOpts := &auth.OAuthOptions{
		Code:         opts.Code,
		CodeVerifier: opts.CodeVerifier,
	}

	accessToken, err := auth.RunInteractiveOAuth(opts.IO, client, signup, oauthOpts)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
	apps, err := client.ListApplications(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	var appDetails *dashboard.Application

	if len(apps) == 0 {
		fmt.Fprintf(opts.IO.Out, "\n%s No applications found. Let's create one.\n", cs.WarningIcon())

		appDetails, err = apputil.CreateAndFetchApplication(opts.IO, client, accessToken, "", opts.AppName)
		if err != nil {
			return err
		}
	} else {
		interactive := opts.IO.CanPrompt()
		app, err := selectApplication(opts, apps, interactive)
		if err != nil {
			return err
		}

		opts.IO.StartProgressIndicatorWithLabel("Fetching application details")
		appDetails, err = client.GetApplication(accessToken, app.ID)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
	}

	profileName := opts.ProfileName
	if profileName == "" {
		profileName = appDetails.Name
	}

	return apputil.ConfigureProfile(opts.IO, opts.Config, appDetails, profileName, opts.Default)
}

func selectApplication(opts *LoginOptions, apps []dashboard.Application, interactive bool) (*dashboard.Application, error) {
	if opts.AppName != "" {
		for i := range apps {
			if apps[i].Name == opts.AppName {
				return &apps[i], nil
			}
		}
		return nil, fmt.Errorf("application %q not found", opts.AppName)
	}

	if len(apps) == 1 {
		return &apps[0], nil
	}

	if !interactive {
		return nil, fmt.Errorf("multiple applications found — use --app-name to select one")
	}

	appNames := make([]string, len(apps))
	for i, app := range apps {
		appNames[i] = fmt.Sprintf("%s (%s)", app.ID, app.Name)
	}

	var selected int
	err := prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select an application:",
			Options: appNames,
		},
		&selected,
	)
	if err != nil {
		return nil, err
	}

	return &apps[selected], nil
}
