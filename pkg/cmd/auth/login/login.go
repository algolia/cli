package login

import (
	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/cli/internal/authflow"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type LoginOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams
}

func NewLoginCmd(f *cmdutil.Factory, runF func(*LoginOptions) error) *cobra.Command {
	opts := &LoginOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Log in to an Algolia account",
		Long: heredoc.Docf(`
			Authenticate with Algolia.
		`, "`"),
		Example: heredoc.Doc(`
			# Start interactive setup
			$ algolia auth login
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoginCmd(opts)
		},
	}

	return cmd
}

func runLoginCmd(opts *LoginOptions) error {
	authflow.AuthFlow("https://www.algolia.com", opts.IO, "Logging in...", nil)

	return nil
}
