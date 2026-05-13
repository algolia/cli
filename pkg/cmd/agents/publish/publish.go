package publish

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type PublishOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID string
}

func NewPublishCmd(f *cmdutil.Factory, runF func(*PublishOptions) error) *cobra.Command {
	opts := &PublishOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "publish <agent-id>",
		Short: "Publish a draft Agent Studio agent",
		Long: heredoc.Doc(`
			Transition an agent from "draft" to "published" so it can be
			invoked by the completions endpoint.
		`),
		Example: heredoc.Doc(`
			$ algolia agents publish 11111111-1111-1111-1111-111111111111
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
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

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Publishing agent")
	agent, err := client.PublishAgent(ctx, opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}
