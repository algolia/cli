package get

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// GetOptions holds the dependencies and flags for the get command.
type GetOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	DomainID          string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `agents domains get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <agent-id> <domain-id>",
		Short: "Get an allowed domain of an agent",
		Args: validators.ExactArgsWithMsg(
			2,
			"agents domains get requires an <agent-id> and a <domain-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# Get the allowed domain "domain_123" of the agent "my-agent"
			$ algolia agents domains get my-agent domain_123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.DomainID = args[1]
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

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching allowed domain")

	res, err := client.GetAllowedDomain(client.NewApiGetAllowedDomainRequest(opts.DomainID, opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
