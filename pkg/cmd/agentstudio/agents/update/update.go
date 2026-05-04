package update

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type UpdateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	ID   string
	File string
	Body agentstudio.AgentConfigUpdate

	PrintFlags *cmdutil.PrintFlags
}

func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "update <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Update an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Update from a JSON file
			$ algolia agentstudio agents update a1b2 -F changes.json

			# Update one field via flag
			$ algolia agentstudio agents update a1b2 --instructions "be more helpful"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			// Update reuses the create-flag set (same property names; all
			// optional in the update body).
			if err := cmdutil.MergeFileAndFlagsInto(opts.IO, opts.File, cmd, cmdutil.AgentConfigCreate, &opts.Body); err != nil {
				return err
			}
			if runF != nil {
				return runF(opts)
			}
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "Agent config JSON `file` (use \"-\" for stdin)")
	cmdutil.AddAgentConfigCreateFlags(cmd)
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating agent")
	agent, err := client.UpdateAgent(opts.ID, opts.Body)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated agent %s\n", cs.SuccessIcon(), agent.ID)
		return nil
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, agent)
}
