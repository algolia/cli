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

	Page  int32
	Limit int32

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd returns the `providers list` command.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all providers",
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# List all providers
			$ algolia providers list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	cmd.Flags().Int32Var(&opts.Page, "page", 0, "Page number")
	cmd.Flags().Int32Var(&opts.Limit, "limit", 0, "Items per page")

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

	req := client.NewApiListProvidersRequest()
	if opts.Page > 0 {
		req = req.WithPage(opts.Page)
	}
	if opts.Limit > 0 {
		req = req.WithLimit(opts.Limit)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching providers")

	res, err := client.ListProviders(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
