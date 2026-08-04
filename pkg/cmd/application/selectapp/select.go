package selectapp

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

type SelectOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	AppID   string
	AppName string

	// NonInteractive disables every prompt and defaults the output to JSON, so
	// the command is usable from scripts.
	NonInteractive bool

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewSelectCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SelectOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select the current application",
		Long: heredoc.Doc(`
			Select an Algolia application to use as the current application for
			all CLI commands. Fetches your applications from the API and lets
			you pick one.
		`),
		Example: heredoc.Doc(`
			# Select interactively
			$ algolia application select

			# Select by name (non-interactive)
			$ algolia application select --app-name "My App"

			# Select by application ID (non-interactive)
			$ algolia application select --app-id "ABCDEF1234"

			# Select from a script: no prompts, JSON on stdout
			$ algolia application select --non-interactive --app-id "ABCDEF1234"
		`),
		Aliases: []string{"use"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.NonInteractive {
				cmdutil.ApplyNonInteractive(opts.IO, opts.PrintFlags)
			}

			// Fail before authenticating: with no selector there is nothing to pick.
			if opts.NonInteractive && opts.AppID == "" && opts.AppName == "" {
				return fmt.Errorf("--app-id or --app-name is required in non-interactive mode")
			}

			app, err := runSelectCmd(opts)
			if err != nil {
				return err
			}

			return printSelection(opts, app)
		},
	}

	cmd.Flags().
		StringVar(&opts.AppID, "app-id", "", "Select application by ID (non-interactive)")
	cmd.Flags().
		StringVar(&opts.AppName, "app-name", "", "Select application by name (non-interactive)")
	cmd.MarkFlagsMutuallyExclusive("app-id", "app-name")
	cmd.Flags().
		BoolVar(&opts.NonInteractive, "non-interactive", false, "Never prompt; output JSON unless --output is set (requires --app-id or --app-name)")
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// printSelection emits the structured document once the flow is done. The
// human-readable flow output has already been written to stderr by then.
func printSelection(opts *SelectOptions, app *dashboard.Application) error {
	if !opts.PrintFlags.HasStructuredOutput() {
		return nil
	}

	if app == nil {
		return fmt.Errorf(
			"no applications found; create one with \"algolia application create\"",
		)
	}

	return opts.PrintFlags.Print(opts.IO, apputil.NewApplicationOutput(opts.Config, app))
}

// Run executes the interactive application-selection flow and returns the
// chosen application. Other commands (e.g. open) use it to ensure an
// application is selected before proceeding. A nil application is returned
// when the account has no applications.
func Run(f *cmdutil.Factory) (*dashboard.Application, error) {
	opts := &SelectOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	return runSelectCmd(opts)
}

func runSelectCmd(opts *SelectOptions) (*dashboard.Application, error) {
	// Move the progress narration to stderr so stdout carries the JSON document
	// only.
	if opts.PrintFlags.HasStructuredOutput() {
		defer cmdutil.RedirectHumanOutput(opts.IO)()
	}

	cs := opts.IO.ColorScheme()
	client := opts.NewDashboardClient(auth.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return nil, err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
	apps, err := client.ListApplications(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		newToken, reAuthErr := auth.ReauthenticateIfExpired(opts.IO, client, err)
		if reAuthErr != nil {
			return nil, reAuthErr
		}
		accessToken = newToken
		opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
		apps, err = client.ListApplications(accessToken)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return nil, err
		}
	}

	if len(apps) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No applications found.\n", cs.WarningIcon())
		fmt.Fprintf(opts.IO.Out, "  Use %s to create one.\n", cs.Bold("algolia application create"))
		return nil, nil
	}

	chosen, err := pickApplication(opts, apps)
	if err != nil {
		return nil, err
	}

	// Reuse a key already stored for this application (keychain, then legacy
	// config.toml) before creating a new one on the dashboard. A migrated
	// application may have a usable key but no stored UUID (legacy config.toml
	// never recorded one); without a UUID, commands like `apikeys rotate` can't
	// target it, so regenerate a fresh CLI-managed key whenever none is on record.
	_, hasUUID := opts.Config.APIKeyUUID(chosen.ID)
	if !hasUUID || !apputil.ReuseExistingAPIKey(opts.Config, chosen) {
		if err := apputil.EnsureAPIKey(opts.IO, client, accessToken, chosen); err != nil {
			return nil, err
		}
	}

	if err := apputil.ConfigureProfile(opts.IO, opts.Config, chosen, "", true); err != nil {
		return nil, err
	}

	return chosen, nil
}

func pickApplication(
	opts *SelectOptions,
	apps []dashboard.Application,
) (*dashboard.Application, error) {
	if opts.AppID != "" {
		for i := range apps {
			if apps[i].ID == opts.AppID {
				return &apps[i], nil
			}
		}
		return nil, fmt.Errorf("application with ID %q not found", opts.AppID)
	}

	if opts.AppName != "" {
		for i := range apps {
			if apps[i].Name == opts.AppName {
				return &apps[i], nil
			}
		}
		return nil, fmt.Errorf("application %q not found", opts.AppName)
	}

	if !opts.IO.CanPrompt() {
		return nil, fmt.Errorf("--app-id or --app-name is required in non-interactive mode")
	}

	cs := opts.IO.ColorScheme()
	profileApps := apputil.ProfileApplicationIDs(opts.Config.ConfiguredProfiles())
	appOptions := make([]string, len(apps))
	for i, app := range apps {
		appOptions[i] = apputil.AppOptionLabel(opts.Config, profileApps, cs, app)
	}

	var selected int
	err := prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select an application:",
			Options: appOptions,
		},
		&selected,
	)
	if err != nil {
		return nil, err
	}

	return &apps[selected], nil
}
