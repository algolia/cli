package list

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// ListOptions holds the dependencies and flags for the list command.
type ListOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	PrintFlags        *cmdutil.PrintFlags
}

// NewListCmd returns the `agents domains list` command.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "list <agent-id>",
		Short:             "List the allowed domains of an agent",
		Args:              validators.ExactArgsWithMsg(1, "agents domains list requires an <agent-id> argument."),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# List the allowed domains of the agent "my-agent"
			$ algolia agents domains list my-agent
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runListCmd(opts)
		},
	}

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

	opts.IO.StartProgressIndicatorWithLabel("Fetching allowed domains")

	res, err := client.ListAgentAllowedDomains(client.NewApiListAgentAllowedDomainsRequest(opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
