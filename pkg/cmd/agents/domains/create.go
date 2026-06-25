package domains

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags
	AgentID, Domain      string
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "create <agent-id> --domain <pattern>",
		Short: "Add a single allowed domain to an agent",
		Example: heredoc.Doc(`
			$ algolia agents domains create <agent-id> --domain https://app.example.com
			$ algolia agents domains create <agent-id> --domain "*.example.com"
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.Domain == "" {
				return cmdutil.FlagErrorf("--domain is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}
	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Domain or pattern (required)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Adding allowed domain")
	d, err := client.CreateAgentAllowedDomain(
		client.NewApiCreateAgentAllowedDomainRequest(opts.AgentID, agentStudio.NewAllowedDomainCreate(opts.Domain)),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, d)
}
