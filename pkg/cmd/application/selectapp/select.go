package selectapp

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/auth/login"
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
			return runSelectCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.AppName, "app-name", "", "Select application by name (non-interactive)")

	return cmd
}

func runSelectCmd(opts *SelectOptions) error {
	cs := opts.IO.ColorScheme()
	client := opts.NewDashboardClient(login.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
	apps, err := client.ListApplications(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		newToken, reAuthErr := auth.ReauthenticateIfExpired(opts.IO, client, err)
		if reAuthErr != nil {
			return reAuthErr
		}
		accessToken = newToken
		opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
		apps, err = client.ListApplications(accessToken)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
	}

	if len(apps) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No applications found.\n", cs.WarningIcon())
		fmt.Fprintf(opts.IO.Out, "  Use %s to create one.\n", cs.Bold("algolia application create"))
		return nil
	}

	chosen, err := pickApplication(opts, apps)
	if err != nil {
		return err
	}

	// If a profile already exists for this app, switch the default
	// and ensure it has an API key.
	if exists, profileName := opts.Config.ApplicationIDExists(chosen.ID); exists {
		// Read the profile BEFORE SetDefaultProfile, because viper.Set() calls
		// inside SetDefaultProfile pollute the override map and cause
		// UnmarshalKey to return empty fields (known viper issue).
		var existingProfile *config.Profile
		for _, p := range opts.Config.ConfiguredProfiles() {
			if p.Name == profileName {
				existingProfile = p
				break
			}
		}

		if err := opts.Config.SetDefaultProfile(profileName); err != nil {
			return fmt.Errorf("failed to set default profile: %w", err)
		}
		fmt.Fprintf(opts.IO.Out, "%s Switched to profile %q (application %s).\n",
			cs.SuccessIcon(), profileName, cs.Bold(chosen.ID))

		if existingProfile != nil && existingProfile.APIKey == "" {
			app := &dashboard.Application{ID: chosen.ID, Name: chosen.Name}
			if err := apputil.EnsureAPIKey(opts.IO, client, accessToken, app); err != nil {
				return err
			}
			existingProfile.ApplicationID = chosen.ID
			existingProfile.APIKey = app.APIKey
			if err := existingProfile.Add(); err != nil {
				return err
			}
			fmt.Fprintf(opts.IO.Out, "%s Profile %q updated with API key.\n",
				cs.SuccessIcon(), profileName)
		}
		return nil
	}

	if err := apputil.EnsureAPIKey(opts.IO, client, accessToken, chosen); err != nil {
		return err
	}

	return apputil.ConfigureProfile(opts.IO, opts.Config, chosen, "", true)
}

func pickApplication(opts *SelectOptions, apps []dashboard.Application) (*dashboard.Application, error) {
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
