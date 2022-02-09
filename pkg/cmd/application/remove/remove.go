package remove

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// RemoveOptions represents the options for the list command
type RemoveOptions struct {
	config *config.Config
	IO     *iostreams.IOStreams

	ApplicationName string

	DoConfirm bool
}

// NewRemoveCmd returns a new instance of RemoveCmd
func NewRemoveCmd(f *cmdutil.Factory, runF func(*RemoveOptions) error) *cobra.Command {
	opts := &RemoveOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "remove <app-name>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.ConfiguredApplicationsCompletionFunc(f),
		Short:             "Remove the specified application",
		Long:              `Remove the specified application from the configuration.`,
		Example: heredoc.Doc(`
			# Remove the application named "my-app" from the configuration
			$ algolia application remove my-app
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ApplicationName = args[0]
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runListCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

// runRemoveCmd executes the remove command
func runListCmd(opts *RemoveOptions) error {
	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(fmt.Sprintf("Are you sure you want to remove the application %q?", opts.ApplicationName), &confirmed)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Application successfully removed: %s\n", cs.SuccessIcon(), opts.ApplicationName)
	}

	return nil
}
