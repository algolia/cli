package duplicate

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type DuplicateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	ID         string
	PrintFlags *cmdutil.PrintFlags
}

func NewDuplicateCmd(f *cmdutil.Factory, runF func(*DuplicateOptions) error) *cobra.Command {
	opts := &DuplicateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "duplicate <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Duplicate an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio agents duplicate a1b2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}
			return runDuplicateCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runDuplicateCmd(opts *DuplicateOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Duplicating agent")
	agent, err := client.DuplicateAgent(opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Duplicated %s -> %s\n", cs.SuccessIcon(), opts.ID, agent.ID)
		return nil
	}
	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, agent)
}
