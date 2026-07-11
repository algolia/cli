package add

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// AddOptions holds the dependencies and flags for the add command.
type AddOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	Domain            string
	PrintFlags        *cmdutil.PrintFlags
}

// NewAddCmd returns the `agents domains add` command.
func NewAddCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &AddOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "add <agent-id> <domain>",
		Short: "Add an allowed domain to an agent",
		Args: validators.ExactArgsWithMsg(
			2,
			"agents domains add requires an <agent-id> and a <domain> argument.",
		),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Allow "https://app.example.com" to embed the agent "my-agent"
			$ algolia agents domains add my-agent https://app.example.com

			# Allow any subdomain of "example.com"
			$ algolia agents domains add my-agent "*.example.com"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Domain = args[1]
			return runAddCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runAddCmd(opts *AddOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Adding allowed domain")

	domainCreate := agentStudio.NewAllowedDomainCreate(opts.Domain)
	res, err := client.CreateAgentAllowedDomain(
		client.NewApiCreateAgentAllowedDomainRequest(opts.AgentID, domainCreate),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
