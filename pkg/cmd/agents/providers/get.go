package providers

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID string
	Show       bool
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <provider-id>",
		Short: "Get an LLM provider authentication by ID",
		Long: heredoc.Doc(`
			Fetch a provider authentication by ID. By default, secret
			fields ("apiKey") are masked. Pass --show-secret to reveal.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers get 11111111-1111-1111-1111-111111111111
			$ algolia agents providers get <id> --show-secret --output json
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields (apiKey, ...) verbatim instead of masking")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Fetching provider")
	p, err := client.GetProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if !opts.Show {
		p.Input = shared.MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}
