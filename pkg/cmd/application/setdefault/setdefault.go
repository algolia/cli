package setdefault

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// SetDefault represents the options for the setdefault command
type AddOptions struct {
	config *config.Config
	IO     *iostreams.IOStreams

	ApplicationName string
}

// NewSetDefaultCmd returns a new instance of SetDefaultCmd
func NewSetDefaultCmd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:               "setdefault [NAME]",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.ConfiguredApplicationsCompletionFunc(f),
		Short:             "Set the default application",
		Example: heredoc.Doc(`
			# Set the default application
			$ algolia application setdefault my-app
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			opts.ApplicationName = args[0]

			return runSetDefaultCmd(opts)
		},
	}

	return cmd
}

// runSetDefaultCmd executes the setdefault command
func runSetDefaultCmd(opts *AddOptions) error {
	opts.config.Application.Name = opts.ApplicationName
	err := opts.config.Application.SetDefault()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	opts.config.Application.LoadDefault()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Application '%s' successfuly set as default application\n", cs.SuccessIcon(), opts.ApplicationName)
	}

	return nil
}
