package login

import (
	"context"
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
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/validators"
)

// LoginOptions holds all options for the login command.
type LoginOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	AppName     string
	ProfileName string
	Region      string
	Default     bool

	// NoBrowser disables automatic browser opening; the authorize URL is
	// printed instead. The CLI still starts a local callback server and
	// waits for the redirect.
	NoBrowser bool

	// NonInteractive disables every prompt and defaults the output to JSON, so
	// the command is usable from scripts. It signs in only: unless --app-name
	// names one, choosing an application is left to `algolia application select`.
	// The browser step is unaffected - the authorize URL still has to be opened
	// for the flow to complete.
	NonInteractive bool

	// PrintFlags is nil for callers that don't expose output flags (e.g. signup).
	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

// loginResult is the machine-readable outcome of the flow. Application is nil
// when the run signed in without configuring one.
type loginResult struct {
	Success     bool                       `json:"success"`
	Email       string                     `json:"email,omitempty"`
	Application *apputil.ApplicationOutput `json:"application,omitempty"`
}

// NewLoginCmd returns a new instance of the login command.
func NewLoginCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
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
			automatically - no code copy-paste required.

			Use --no-browser if the browser cannot be opened automatically
			(e.g. SSH sessions, containers). The URL will be printed for you
			to open manually; the CLI still waits for the redirect.

			Use --non-interactive to sign in without any prompt and print the
			result as JSON. It signs in only: no application is configured unless
			--app-name names one, so pick one afterwards with
			"algolia application select".
		`),
		Example: heredoc.Doc(`
			# Sign in interactively (opens browser)
			$ algolia auth login

			# Auto-select an application by name
			$ algolia auth login --app-name "My App" --default

			# Print the URL instead of opening the browser
			$ algolia auth login --no-browser

			# Sign in from a script: no prompts, JSON on stdout, no application
			$ algolia auth login --non-interactive
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			applyNonInteractive(opts)
			return runLoginCmd(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.AppName, "app-name", "", "Auto-select application by name")
	cmd.Flags().
		StringVar(&opts.Region, "region", "", "Region for the first application when the account has none (e.g. EU, UK, USC, USE, USW)")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Alias for the application (defaults to the application name)")
	cmd.Flags().BoolVar(&opts.Default, "default", true, "Set the application as the current one")
	cmd.Flags().BoolVar(&opts.NoBrowser, "no-browser", false, "Print the authorize URL instead of opening the browser")
	cmd.Flags().
		BoolVar(&opts.NonInteractive, "non-interactive", false, "Sign in without prompting and print JSON; configures no application unless --app-name is set")
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// applyNonInteractive turns off every prompt and defaults the output to JSON,
// leaving an explicit --output untouched.
func applyNonInteractive(opts *LoginOptions) {
	if !opts.NonInteractive {
		return
	}

	opts.IO.SetNeverPrompt(true)

	if opts.PrintFlags != nil && opts.PrintFlags.OutputFormat != nil &&
		!opts.PrintFlags.HasStructuredOutput() {
		*opts.PrintFlags.OutputFormat = "json"
	}
}

func runLoginCmd(ctx context.Context, opts *LoginOptions) error {
	return RunOAuthFlow(ctx, opts, false)
}

// RunOAuthFlow runs the full browser-based OAuth + profile setup flow.
// If signup is true, the browser opens to the sign-up page instead of sign-in.
func RunOAuthFlow(ctx context.Context, opts *LoginOptions, signup bool) error {
	flow := telemetry.FlowLogin
	if signup {
		flow = telemetry.FlowSignup
	}
	tracker := telemetry.NewFlowTracker()
	telemetry.TrackEvent(ctx, telemetry.AuthStarted(flow, opts.NoBrowser))

	app, err := runOAuthFlowSteps(ctx, opts, signup, tracker)
	trackOAuthFlowOutcome(ctx, flow, tracker, err)
	if err != nil {
		return err
	}

	return printLoginResult(opts, app)
}

// printLoginResult emits the structured document once the flow is done, and
// nothing at all in the default human-readable mode.
func printLoginResult(opts *LoginOptions, app *dashboard.Application) error {
	if !opts.PrintFlags.HasStructuredOutput() {
		return nil
	}

	result := loginResult{Success: true}
	if token := auth.LoadToken(); token != nil {
		result.Email = token.Email
	}
	if app != nil {
		application := apputil.NewApplicationOutput(opts.Config, app)
		result.Application = &application
	}

	return opts.PrintFlags.Print(opts.IO, result)
}

// trackOAuthFlowOutcome reports how the auth flow ended: completed, aborted by
// the user, or failed.
func trackOAuthFlowOutcome(
	ctx context.Context,
	flow telemetry.Flow,
	tracker *telemetry.FlowTracker,
	err error,
) {
	switch {
	case err == nil:
		telemetry.TrackEvent(ctx, telemetry.AuthCompleted(flow, tracker))
	case cmdutil.IsUserCancellation(err):
		telemetry.TrackEvent(ctx, telemetry.AuthAborted(flow, tracker))
	default:
		telemetry.TrackEvent(ctx, telemetry.AuthFailed(flow, tracker, err))
	}
}

func runOAuthFlowSteps(
	ctx context.Context,
	opts *LoginOptions,
	signup bool,
	tracker *telemetry.FlowTracker,
) (*dashboard.Application, error) {
	// Drop progress narration so the command emits the JSON document only. This
	// includes the authorize URL, so --non-interactive relies on the browser
	// opening: pair it with --no-browser only if the URL isn't needed.
	if opts.PrintFlags.HasStructuredOutput() {
		defer cmdutil.DiscardHumanOutput(opts.IO)()
	}

	cs := opts.IO.ColorScheme()
	client := opts.NewDashboardClient(auth.OAuthClientID())

	openBrowser := !opts.NoBrowser
	accessToken, err := auth.RunOAuth(opts.IO, client, signup, openBrowser, tracker)
	if err != nil {
		return nil, err
	}

	applyStoredIdentity(ctx)

	// Non-interactive login authenticates only: with no name to go on there is
	// nothing to pick, and creating an application (or a key) behind a script's
	// back would be a side effect it never asked for. `algolia application
	// select` configures one afterwards.
	if opts.NonInteractive && opts.AppName == "" {
		return nil, nil
	}

	tracker.SetStep(telemetry.StepAppsFetch)
	opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
	apps, err := client.ListApplications(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return nil, err
	}

	var appDetails *dashboard.Application

	if len(apps) == 0 {
		fmt.Fprintf(opts.IO.Out, "\n%s No applications found. Let's create one.\n", cs.WarningIcon())

		tracker.SetStep(telemetry.StepAppCreate)

		appName := opts.AppName
		if appName == "" && opts.IO.CanPrompt() {
			appName, err = apputil.PromptName()
			if err != nil {
				return nil, err
			}
		}

		// No create tracker: this creation belongs to the auth funnel, which
		// stays on the app_create step.
		appDetails, _, err = apputil.CreateAndFetchApplication(opts.IO, client, accessToken, opts.Region, appName, nil)
		if err != nil {
			return nil, err
		}
	} else {
		tracker.SetStep(telemetry.StepAppSelect)
		interactive := opts.IO.CanPrompt()
		app, err := selectApplication(opts, apps, interactive)
		if err != nil {
			return nil, err
		}

		appDetails = app
		// A migrated application may have a usable key but no stored UUID, which
		// commands like `apikeys rotate` need; regenerate a fresh CLI-managed key
		// whenever none is on record (mirrors `application select`).
		_, hasUUID := opts.Config.APIKeyUUID(appDetails.ID)
		if !hasUUID || !apputil.ReuseExistingAPIKey(opts.Config, appDetails) {
			if err := apputil.EnsureAPIKey(opts.IO, client, accessToken, appDetails); err != nil {
				return nil, err
			}
		}
	}

	profileName := opts.ProfileName
	if profileName == "" {
		profileName = appDetails.Name
	}

	tracker.SetStep(telemetry.StepProfileConfigure)
	if err := apputil.ConfigureProfile(opts.IO, opts.Config, appDetails, profileName, opts.Default); err != nil {
		return nil, err
	}

	return appDetails, nil
}

// applyStoredIdentity copies the persisted user identity from the stored token
// onto the request's telemetry metadata so the flow's own events carry the user.
// The Identify itself is sent once at command completion. It reports whether an
// identity was applied.
func applyStoredIdentity(ctx context.Context) bool {
	token := auth.LoadToken()
	if token == nil || token.UserID == "" {
		return false
	}

	metadata := telemetry.GetEventMetadata(ctx)
	if metadata == nil {
		return false
	}

	metadata.SetUser(token.UserID, token.Email, token.Name)
	return true
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
		// stderr: this listing explains the error below, so it must not be
		// dropped or mixed into structured output.
		fmt.Fprintf(opts.IO.ErrOut, "Multiple applications found:\n")
		for i, app := range apps {
			fmt.Fprintf(opts.IO.ErrOut, "  %d. %s (%s)\n", i+1, app.Name, app.ID)
		}
		fmt.Fprintf(opts.IO.ErrOut, "Use --app-name to select one.\n")
		return nil, fmt.Errorf("multiple applications found - use --app-name to select one")
	}

	cs := opts.IO.ColorScheme()
	profileApps := apputil.ProfileApplicationIDs(opts.Config.ConfiguredProfiles())
	appNames := make([]string, len(apps))
	for i, app := range apps {
		appNames[i] = apputil.AppOptionLabel(opts.Config, profileApps, cs, app)
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
