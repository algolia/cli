package create

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/auth/login"
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

	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:     f.IOStreams,
		Config: f.Config,
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
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "My First Application", "Name for the application")
	cmd.Flags().StringVar(&opts.Region, "region", "", "Region code (e.g. CA, US, EU)")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Name for the CLI profile (defaults to app name)")
	cmd.Flags().BoolVar(&opts.Default, "default", false, "Set the new profile as the default")

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client := opts.NewDashboardClient(login.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	appDetails, err := apputil.CreateAndFetchApplication(opts.IO, client, accessToken, opts.Region, opts.Name)
	if err != nil {
		return err
	}

	profileName := opts.ProfileName
	if profileName == "" {
		profileName = opts.Name
	}

	return apputil.ConfigureProfile(opts.IO, opts.Config, appDetails, profileName, opts.Default)
}
