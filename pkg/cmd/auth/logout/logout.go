package logout

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type LogoutOptions struct {
	IO                 *iostreams.IOStreams
	NewDashboardClient func(clientID string) *dashboard.Client
}

func NewLogoutCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LogoutOptions{
		IO: f.IOStreams,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Sign out of your Algolia account",
		Long: heredoc.Doc(`
			Sign out by revoking the stored OAuth tokens on the server
			and removing them from the local keychain.
		`),
		Example: heredoc.Doc(`
			$ algolia auth logout
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogoutCmd(opts)
		},
	}

	return cmd
}

func runLogoutCmd(opts *LogoutOptions) error {
	cs := opts.IO.ColorScheme()
	stored := auth.LoadToken()

	if stored == nil {
		fmt.Fprintf(opts.IO.Out, "%s Already signed out.\n", cs.SuccessIcon())
		return nil
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	if stored.AccessToken != "" {
		if err := client.RevokeToken(stored.AccessToken); err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to revoke access token: %s\n", cs.WarningIcon(), err)
		}
	}

	if stored.RefreshToken != "" {
		if err := client.RevokeToken(stored.RefreshToken); err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to revoke refresh token: %s\n", cs.WarningIcon(), err)
		}
	}

	auth.ClearToken()

	fmt.Fprintf(opts.IO.Out, "%s Signed out successfully.\n", cs.SuccessIcon())
	return nil
}
