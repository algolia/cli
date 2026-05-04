package unpublish

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type UnpublishOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	ID         string
	PrintFlags *cmdutil.PrintFlags
}

func NewUnpublishCmd(f *cmdutil.Factory, runF func(*UnpublishOptions) error) *cobra.Command {
	opts := &UnpublishOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "unpublish <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Unpublish an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio agents unpublish a1b2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}
			return runUnpublishCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUnpublishCmd(opts *UnpublishOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Unpublishing agent")
	agent, err := client.UnpublishAgent(opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Unpublished agent %s\n", cs.SuccessIcon(), agent.ID)
		return nil
	}
	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, agent)
}
