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

// RemoveOptions represents the options for the remove command
type RemoveOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	Profile        string
	DefaultProfile string

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
		Use:               "remove <profile>",
		Args:              validators.ExactArgsWithDefaultRequiredMsg(1),
		ValidArgsFunction: cmdutil.ConfiguredProfilesCompletionFunc(f),
		Short:             "Remove the specified profile",
		Long:              `Remove the specified profile from the configuration.`,
		Example: heredoc.Doc(`
			# Remove the profile named "my-app" from the configuration
			$ algolia profile remove my-app
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Profile = args[0]

			if opts.config.Default() != nil {
				opts.DefaultProfile = opts.config.Default().Name
			} else {
				opts.DefaultProfile = ""
			}

			if !confirm && opts.Profile == opts.DefaultProfile {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if !opts.config.ProfileExists(opts.Profile) {
				return fmt.Errorf("the specified profile does not exist: '%s'", opts.Profile)
			}

			if runF != nil {
				return runF(opts)
			}

			return runRemoveCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

// runRemoveCmd executes the remove command
func runRemoveCmd(opts *RemoveOptions) error {
	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(fmt.Sprintf("Are you sure you want to remove '%s', the default profile?", opts.Profile), &confirmed)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}
	}

	err := opts.config.RemoveProfile(opts.Profile)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		extra := "."
		if opts.DefaultProfile == opts.Profile {
			extra = ". Set a new default profile with 'algolia profile setdefault'."
		}
		if len(opts.config.ConfiguredProfiles()) == 0 {
			extra = ". Add a profile with 'algolia profile add'."
		}
		fmt.Fprintf(opts.IO.Out, "%s '%s' removed successfully%s\n", cs.SuccessIcon(), opts.Profile, extra)
	}

	return nil
}
