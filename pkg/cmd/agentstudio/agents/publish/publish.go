package publish

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type PublishOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	ID         string
	PrintFlags *cmdutil.PrintFlags
}

func NewPublishCmd(f *cmdutil.Factory, runF func(*PublishOptions) error) *cobra.Command {
	opts := &PublishOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "publish <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Publish an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio agents publish a1b2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}
			return runPublishCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runPublishCmd(opts *PublishOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Publishing agent")
	agent, err := client.PublishAgent(opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Published agent %s\n", cs.SuccessIcon(), agent.ID)
		return nil
	}
	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, agent)
}
