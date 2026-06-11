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

	AppName string

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewSelectCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SelectOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select an application to use as the active profile",
		Long: heredoc.Doc(`
			Select an Algolia application to use as the default CLI profile.
			Fetches your applications from the API and lets you pick one.

			If the selected application already has a local profile, it is set
			as the default. Otherwise, a new profile is created and set as default.
		`),
		Example: heredoc.Doc(`
			# Select interactively
			$ algolia application select

			# Select by name (non-interactive)
			$ algolia application select --app-name "My App"
		`),
		Aliases: []string{"use"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := runSelectCmd(opts)
			return err
		},
	}

	cmd.Flags().
		StringVar(&opts.AppName, "app-name", "", "Select application by name (non-interactive)")

	return cmd
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
	// config.toml) before creating a new one on the dashboard.
	if !apputil.ReuseExistingAPIKey(opts.Config, chosen) {
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
	if opts.AppName != "" {
		for i := range apps {
			if apps[i].Name == opts.AppName {
				return &apps[i], nil
			}
		}
		return nil, fmt.Errorf("application %q not found", opts.AppName)
	}

	if !opts.IO.CanPrompt() {
		return nil, fmt.Errorf("--app-name is required in non-interactive mode")
	}

	configuredProfiles := opts.Config.ConfiguredProfiles()
	configuredAppIDs := make(map[string]string)
	for _, p := range configuredProfiles {
		configuredAppIDs[p.ApplicationID] = p.Name
	}

	cs := opts.IO.ColorScheme()
	appOptions := make([]string, len(apps))
	for i, app := range apps {
		label := fmt.Sprintf("%s (%s)", app.ID, app.Name)
		if profileName, ok := configuredAppIDs[app.ID]; ok {
			appOptions[i] = fmt.Sprintf("%s  %s", label, cs.Greenf("profile: %s", profileName))
		} else {
			appOptions[i] = label
		}
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
