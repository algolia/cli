package login

import (
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
	"github.com/algolia/algolia-cli/pkg/login"
	"github.com/algolia/algolia-cli/pkg/validators"
)

type LoginOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams
}

func NewLoginCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		Config: f.Config,
		IO:     f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "login",
		Args:  validators.NoArgs,
		Short: "Login to your Algolia account",
		Long:  `Login to your Algolia account to setup the CLI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoginCmd(opts)
		},
	}

	return cmd
}

func runLoginCmd(opts *LoginOptions) error {
	return login.InteractiveLogin(opts.Config, opts.IO)
}
