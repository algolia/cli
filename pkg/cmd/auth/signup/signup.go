package signup

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/cmd/auth/login"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/validators"
)

// NewSignupCmd returns a new instance of the signup command.
func NewSignupCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &login.LoginOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Create a new Algolia account",
		Long: heredoc.Doc(`
			Create a new Algolia account via the browser.
			Opens the Algolia Dashboard sign-up page, then completes the OAuth
			authorization code flow to configure the CLI.
		`),
		Example: heredoc.Doc(`
			# Create a new account (opens browser to sign-up page)
			$ algolia auth signup
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return login.RunOAuthFlow(cmd.Context(), opts, true)
		},
	}

	cmd.Flags().StringVar(&opts.AppName, "app-name", "", "Name for the first application")
	cmd.Flags().
		StringVar(&opts.Region, "region", "", "Region for the first application (e.g. EU, UK, USC, USE, USW)")
	cmd.Flags().StringVar(&opts.ProfileName, "profile-name", "", "Alias for the application (defaults to the application name)")
	cmd.Flags().BoolVar(&opts.Default, "default", true, "Set the application as the current one")
	cmd.Flags().BoolVar(&opts.NoBrowser, "no-browser", false, "Print the authorize URL instead of opening the browser")

	return cmd
}
