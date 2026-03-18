package login

import (
	"fmt"

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

// LoginOptions holds all options for the login command.
type LoginOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	AppName     string
	ProfileName string
	Default     bool

	// NoBrowser disables automatic browser opening; the authorize URL is
	// printed instead. The CLI still starts a local callback server and
	// waits for the redirect.
	NoBrowser bool

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

			A local HTTP server is started to receive the OAuth redirect
			automatically — no code copy-paste required.

			Use --no-browser if the browser cannot be opened automatically
			(e.g. SSH sessions, containers). The URL will be printed for you
			to open manually; the CLI still waits for the redirect.
		`),
		Example: heredoc.Doc(`
			# Sign in interactively (opens browser)
			$ algolia auth login

			# Auto-select an application by name
			$ algolia auth login --app-name "My App" --default

			# Print the URL instead of opening the browser
			$ algolia auth login --no-browser
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoginCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.AppName, "app-name", "", "Auto-select application by name")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Name for the CLI profile (defaults to application name)")
	cmd.Flags().BoolVar(&opts.Default, "default", true, "Set the profile as the default")
	cmd.Flags().BoolVar(&opts.NoBrowser, "no-browser", false, "Print the authorize URL instead of opening the browser")

	return cmd
}

func runLoginCmd(opts *LoginOptions) error {
	return RunOAuthFlow(opts, false)
}

// RunOAuthFlow runs the full browser-based OAuth + profile setup flow.
// If signup is true, the browser opens to the sign-up page instead of sign-in.
func RunOAuthFlow(opts *LoginOptions, signup bool) error {
	cs := opts.IO.ColorScheme()
	client := opts.NewDashboardClient(auth.OAuthClientID())

	openBrowser := !opts.NoBrowser
	accessToken, err := auth.RunOAuth(opts.IO, client, signup, openBrowser)
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
		fmt.Fprintf(opts.IO.Out, "Multiple applications found:\n")
		for i, app := range apps {
			fmt.Fprintf(opts.IO.Out, "  %d. %s (%s)\n", i+1, app.Name, app.ID)
		}
		fmt.Fprintf(opts.IO.Out, "Use --app-name to select one.\n")
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
