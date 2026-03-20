package create

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Name        string
	Region      string
	ProfileName string
	Default     bool
	DryRun      bool

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Algolia application",
		Long: heredoc.Doc(`
			Create a new Algolia application and optionally configure it as a CLI profile.
			Requires an active session (run "algolia auth login" first).
		`),
		Example: heredoc.Doc(`
			# Create an application interactively
			$ algolia application create

			# Create with specific options
			$ algolia application create --name "My App" --region CA

			# Create and set as default profile
			$ algolia application create --name "My App" --region CA --default

			# Preview what would be created without actually creating it
			$ algolia application create --name "My App" --region CA --dry-run
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "My First Application", "Name for the application")
	cmd.Flags().StringVar(&opts.Region, "region", "", "Region code (e.g. CA, US, EU)")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Name for the CLI profile (defaults to app name)")
	cmd.Flags().BoolVar(&opts.Default, "default", false, "Set the new profile as the default")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the create request without sending it")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	if opts.DryRun {
		summary := map[string]any{
			"action":  "create_application",
			"name":    opts.Name,
			"region":  opts.Region,
			"default": opts.Default,
			"dryRun":  true,
		}
		return cmdutil.PrintRunSummary(
			opts.IO,
			opts.PrintFlags,
			summary,
			fmt.Sprintf("Dry run: would create application %q in region %q", opts.Name, opts.Region),
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	appDetails, err := apputil.CreateAndFetchApplication(opts.IO, client, accessToken, opts.Region, opts.Name)
	if err != nil {
		return err
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, appDetails)
	}

	return apputil.ConfigureProfile(opts.IO, opts.Config, appDetails, opts.ProfileName, opts.Default)
}
