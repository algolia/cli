package update

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

type UpdateOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Name string

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Rename the current Algolia application",
		Long: heredoc.Doc(`
			Rename the application associated with the current CLI profile.
			Requires an active application to be selected (run "algolia application select" first).
		`),
		Example: heredoc.Doc(`
			# Rename the current application
			$ algolia application update --name "New Name"
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the current application")
	_ = cmd.MarkFlagRequired("name")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	cs := opts.IO.ColorScheme()

	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return fmt.Errorf(
			"no current application configured; configure a profile with \"algolia profile add\" or \"algolia application select\" first: %w",
			err,
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating application")
	app, err := client.UpdateApplication(accessToken, appID, opts.Name)
	opts.IO.StopProgressIndicator()
	if err != nil {
		newToken, reAuthErr := auth.ReauthenticateIfExpired(opts.IO, client, err)
		if reAuthErr != nil {
			return reAuthErr
		}
		accessToken = newToken
		opts.IO.StartProgressIndicatorWithLabel("Updating application")
		app, err = client.UpdateApplication(accessToken, appID, opts.Name)
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
		return p.Print(opts.IO, app)
	}

	fmt.Fprintf(
		opts.IO.Out,
		"%s Application %s renamed to %q.\n",
		cs.SuccessIcon(),
		cs.Bold(app.ID),
		app.Name,
	)
	return nil
}
