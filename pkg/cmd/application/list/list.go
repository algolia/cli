package list

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

type ListOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your Algolia applications and optionally select one to configure",
		Long: heredoc.Doc(`
			List all Algolia applications associated with your account.
			Requires an active session (run "algolia auth login" first).
			Applications that already have a local CLI profile are marked.
			You can select an unconfigured application to add it as a CLI profile.
		`),
		Example: heredoc.Doc(`
			# List applications
			$ algolia application list
		`),
		Aliases: []string{"ls"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
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

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, apps)
	}

	if len(apps) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No applications found.\n", cs.WarningIcon())
		fmt.Fprintf(opts.IO.Out, "  Use %s to create one.\n", cs.Bold("algolia application create"))
		return nil
	}

	configuredProfiles := opts.Config.ConfiguredProfiles()
	configuredAppIDs := make(map[string]string)
	for _, p := range configuredProfiles {
		configuredAppIDs[p.ApplicationID] = p.Name
	}

	fmt.Fprintf(opts.IO.Out, "\nYour applications:\n\n")
	unconfigured := make([]dashboard.Application, 0)

	for _, app := range apps {
		profileName, configured := configuredAppIDs[app.ID]
		label := fmt.Sprintf("  %s  %s", app.ID, app.Name)
		if configured {
			fmt.Fprintf(opts.IO.Out, "%s  %s\n", label, cs.Greenf("(profile: %s)", profileName))
		} else {
			fmt.Fprintf(opts.IO.Out, "%s  %s\n", label, cs.Gray("(not configured)"))
			unconfigured = append(unconfigured, app)
		}
	}

	fmt.Fprintln(opts.IO.Out)

	if len(unconfigured) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s All applications are already configured as CLI profiles.\n", cs.SuccessIcon())
		return nil
	}

	if !opts.IO.CanPrompt() {
		return nil
	}

	var wantConfigure bool
	err = prompt.SurveyAskOne(
		&survey.Confirm{
			Message: "Would you like to configure an unconfigured application as a CLI profile?",
			Default: true,
		},
		&wantConfigure,
	)
	if err != nil || !wantConfigure {
		return err
	}

	appOptions := make([]string, len(unconfigured))
	for i, app := range unconfigured {
		appOptions[i] = fmt.Sprintf("%s (%s)", app.ID, app.Name)
	}

	var selected int
	err = prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select an application to configure:",
			Options: appOptions,
		},
		&selected,
	)
	if err != nil {
		return err
	}

	appDetails := &unconfigured[selected]

	if err := apputil.EnsureAPIKey(opts.IO, client, accessToken, appDetails); err != nil {
		return err
	}

	var setDefault bool
	err = prompt.SurveyAskOne(
		&survey.Confirm{
			Message: "Set as the default profile?",
			Default: false,
		},
		&setDefault,
	)
	if err != nil {
		return err
	}

	return apputil.ConfigureProfile(opts.IO, opts.Config, appDetails, "", setDefault)
}
