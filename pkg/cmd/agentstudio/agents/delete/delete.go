package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	ID      string
	Confirm bool
}

func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}
	cmd := &cobra.Command{
		Use:     "delete <agent_id>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Delete an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio agents delete a1b2 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if !opts.Confirm {
				return fmt.Errorf("--confirm is required to delete agent %s", opts.ID)
			}
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Confirm, "confirm", false, "Skip the confirmation prompt and delete immediately")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting agent")
	err = client.DeleteAgent(opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted agent %s\n", cs.SuccessIcon(), opts.ID)
	}
	return nil
}
