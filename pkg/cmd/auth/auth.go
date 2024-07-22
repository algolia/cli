package auth

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"

	"github.com/algolia/cli/pkg/cmd/auth/login"
)

func NewAuthCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Aliases: []string{"authenticate"},
		Short:   "Authenticate with Algolia",
	}

	cmd.AddCommand(login.NewLoginCmd(f, nil))

	return cmd
}
