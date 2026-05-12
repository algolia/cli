package create

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type CreateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	File string
	Body agentstudio.AgentConfigCreate

	PrintFlags *cmdutil.PrintFlags
}

func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "create",
		Args:  cobra.NoArgs,
		Short: "Create an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Create from a JSON file
			$ algolia agentstudio agents create -F agent.json

			# Create from individual flags
			$ algolia agentstudio agents create --name "Helper" --instructions "Be helpful"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.MergeFileAndFlagsInto(opts.IO, opts.File, cmd, cmdutil.AgentConfigCreate, &opts.Body); err != nil {
				return err
			}
			if opts.Body.Name == "" || opts.Body.Instructions == "" {
				return fmt.Errorf("--name and --instructions are required (or provide them in -F file)")
			}
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "Agent config JSON `file` (use \"-\" for stdin)")
	cmdutil.AddAgentConfigCreateFlags(cmd)
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Creating agent")
	agent, err := client.CreateAgent(opts.Body)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created agent %s (%s)\n", cs.SuccessIcon(), agent.Name, agent.ID)
		return nil
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, agent)
}
