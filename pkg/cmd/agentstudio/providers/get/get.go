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
	ProviderID        string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `providers get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "get <provider-id>",
		Short:             "Get a provider by ID",
		Args:              validators.ExactArgsWithMsg(1, "providers get requires a <provider-id> argument."),
		ValidArgsFunction: cmdutil.ProviderIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# Get the provider with ID "my-provider"
			$ algolia providers get my-provider
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
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

	opts.IO.StartProgressIndicatorWithLabel("Fetching provider")

	res, err := client.GetProvider(client.NewApiGetProviderRequest(opts.ProviderID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
