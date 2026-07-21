package current

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CurrentOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

type currentApplication struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
	Plan  string `json:"plan"`
}

func NewCurrentCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CurrentOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the currently selected application",
		Long: heredoc.Doc(`
			Show which Algolia application is currently selected, along with its
			name and plan.

			The application ID (and its alias, when set) is shown even if the name
			and plan can't be fetched.
		`),
		Example: heredoc.Doc(`
			# Show the current application
			$ algolia application current

			# Output as JSON
			$ algolia application current --output json
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCurrentCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCurrentCmd(opts *CurrentOptions) error {
	cs := opts.IO.ColorScheme()

	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return fmt.Errorf(
			"no current application configured; run \"algolia application select\" or \"algolia auth login\" first: %w",
			err,
		)
	}

	current := currentApplication{ID: appID}
	if alias, ok := opts.Config.ApplicationAlias(appID); ok {
		current.Alias = alias
	}

	app, signedOut := fetchApplication(opts, appID)
	if app != nil {
		current.Name = app.Name
		current.Plan = app.PlanLabel
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, current)
	}

	fmt.Fprintf(opts.IO.Out, "%s Current application: %s\n", cs.SuccessIcon(), cs.Bold(appID))
	if current.Alias != "" {
		fmt.Fprintf(opts.IO.Out, "  Alias: %s\n", current.Alias)
	}
	if current.Name != "" {
		fmt.Fprintf(opts.IO.Out, "  Name:  %s\n", current.Name)
	}
	if current.Plan != "" {
		fmt.Fprintf(opts.IO.Out, "  Plan:  %s\n", current.Plan)
	}
	if current.Name == "" && current.Plan == "" {
		if signedOut {
			fmt.Fprintf(
				opts.IO.Out,
				"%s Sign in with \"algolia auth login\" to see the application name and plan.\n",
				cs.WarningIcon(),
			)
		} else {
			fmt.Fprintf(
				opts.IO.Out,
				"%s Couldn't fetch the application name and plan; showing the selected application only.\n",
				cs.WarningIcon(),
			)
		}
	}

	return nil
}

func fetchApplication(opts *CurrentOptions, appID string) (*dashboard.Application, bool) {
	client := opts.NewDashboardClient(auth.OAuthClientID())

	token, err := auth.GetValidToken(client)
	if err != nil {
		return nil, true
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching application")
	app, err := client.GetApplication(token, appID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return nil, false
	}

	return app, false
}
