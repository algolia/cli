package unpublish

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type UnpublishOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID string
}

func NewUnpublishCmd(f *cmdutil.Factory, runF func(*UnpublishOptions) error) *cobra.Command {
	opts := &UnpublishOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "unpublish <agent-id>",
		Short: "Unpublish a published Agent Studio agent",
		Long: heredoc.Doc(`
			Move an agent from "published" back to "draft". The completions
			endpoint will reject calls to a draft agent until it is republished.
		`),
		Example: heredoc.Doc(`
			$ algolia agents unpublish 11111111-1111-1111-1111-111111111111
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

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Unpublishing agent")
	agent, err := client.UnpublishAgent(ctx, opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}
