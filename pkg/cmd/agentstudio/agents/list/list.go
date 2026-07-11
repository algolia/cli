package list

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// ListOptions holds the dependencies and flags for the list command.
type ListOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)

	Page       int32
	Limit      int32
	ProviderID string

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd returns the `agents list` command.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all agents",
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# List all agents
			$ algolia agents list

			# List agents using a specific provider
			$ algolia agents list --provider-id my-provider
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	cmd.Flags().Int32Var(&opts.Page, "page", 0, "Page number")
	cmd.Flags().Int32Var(&opts.Limit, "limit", 0, "Items per page")
	cmd.Flags().StringVar(&opts.ProviderID, "provider-id", "", "Filter by provider ID")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	req := client.NewApiListAgentsRequest()
	if opts.Page > 0 {
		req = req.WithPage(opts.Page)
	}
	if opts.Limit > 0 {
		req = req.WithLimit(opts.Limit)
	}
	if opts.ProviderID != "" {
		req = req.WithProviderId(opts.ProviderID)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching agents")

	res, err := client.ListAgents(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
