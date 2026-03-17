package auth

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/auth/login"
	"github.com/algolia/cli/pkg/cmd/auth/logout"
	"github.com/algolia/cli/pkg/cmd/auth/signup"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAuthCmd returns a new command for authentication.
func NewAuthCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with your Algolia account",
	}

	auth.DisableAuthCheck(cmd)

	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(logout.NewLogoutCmd(f))
	cmd.AddCommand(signup.NewSignupCmd(f))

	return cmd
}
