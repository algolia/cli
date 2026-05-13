package get

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig
	Ctx    context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	PrintFlags *cmdutil.PrintFlags

	AgentID string
}

func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <agent-id>",
		Short: "Get an Agent Studio agent by ID",
		Long: heredoc.Doc(`
			Fetch a single agent (including its config and tools) by ID.

			Output defaults to JSON because agent payloads are nested and
			designed for piping into jq, files, or templating.
		`),
		Example: heredoc.Doc(`
			# Print the agent as JSON
			$ algolia agents get 11111111-1111-1111-1111-111111111111

			# Extract a single field with jsonpath
			$ algolia agents get 11111111-1111-1111-1111-111111111111 --output jsonpath --template '{$.name}'
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
			return runGetCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching agent")
	agent, err := client.GetAgent(ctx, opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}
