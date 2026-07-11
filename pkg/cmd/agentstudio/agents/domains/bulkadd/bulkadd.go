package bulkadd

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// BulkAddOptions holds the dependencies and flags for the bulk-add command.
type BulkAddOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	Domains           []string
	PrintFlags        *cmdutil.PrintFlags
}

// NewBulkAddCmd returns the `agents domains bulk-add` command.
func NewBulkAddCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &BulkAddOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "bulk-add <agent-id> <domain>...",
		Short:             "Add multiple allowed domains to an agent",
		Args:              validators.AtLeastNArgs(2),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Allow multiple domains to embed the agent "my-agent"
			$ algolia agents domains bulk-add my-agent https://a.example.com https://b.example.com
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Domains = args[1:]
			return runBulkAddCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runBulkAddCmd(opts *BulkAddOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Adding allowed domains")

	bulkInsert := agentStudio.NewAllowedDomainBulkInsert(opts.Domains)
	res, err := client.BulkCreateAllowedDomains(
		client.NewApiBulkCreateAllowedDomainsRequest(opts.AgentID, bulkInsert),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
