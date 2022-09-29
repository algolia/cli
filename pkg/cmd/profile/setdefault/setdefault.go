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

// SetDefaultOptions represents the options for the setdefault command
type SetDefaultOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	Profile string
}

// NewSetDefaultCmd returns a new instance of SetDefaultCmd
func NewSetDefaultCmd(f *cmdutil.Factory, runF func(*SetDefaultOptions) error) *cobra.Command {
	opts := &SetDefaultOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:               "setdefault <profile>",
		Args:              validators.ExactArgsWithDefaultRequiredMsg(1),
		ValidArgsFunction: cmdutil.ConfiguredProfilesCompletionFunc(f),
		Short:             "Set the default profile",
		Example: heredoc.Doc(`
			# Set the default profile to "my-app"
			$ algolia profile setdefault my-app
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Profile = args[0]

			if !opts.config.ProfileExists(opts.Profile) {
				return fmt.Errorf("the specified profile does not exist: '%s'", opts.Profile)
			}

			if runF != nil {
				return runF(opts)
			}

			return runSetDefaultCmd(opts)
		},
	}

	return cmd
}

// runSetDefaultCmd executes the setdefault command
func runSetDefaultCmd(opts *SetDefaultOptions) error {
	var defaultName string
	for _, profile := range opts.config.ConfiguredProfiles() {
		if profile.Default {
			defaultName = profile.Name
		}
	}

	err := opts.config.SetDefaultProfile(opts.Profile)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	opts.config.Profile().LoadDefault()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		if defaultName != "" {
			fmt.Fprintf(opts.IO.Out, "%s Default profile successfuly changed from '%s' to '%s'.\n", cs.SuccessIcon(), defaultName, opts.Profile)
		} else {
			fmt.Fprintf(opts.IO.Out, "%s Default profile successfuly set to '%s'.\n", cs.SuccessIcon(), opts.Profile)
		}
	}

	return nil
}
