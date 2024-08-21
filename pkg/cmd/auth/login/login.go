package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/cli/internal/authflow"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// LoginOptions holds the options for the login command
type LoginOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams
}

// NewLoginCmd creates and returns a login command
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

// runLoginCmd executes the login command
func runLoginCmd(opts *LoginOptions) error {
	cs := opts.IO.ColorScheme()
	token, refreshToken, err := authflow.AuthFlow(opts.IO, "Logging in...")
	if err != nil {
		return err
	}

	err = opts.Config.Auth().Login(token, refreshToken)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "%s Successfully logged in\n", cs.SuccessIcon())

	return nil
}
