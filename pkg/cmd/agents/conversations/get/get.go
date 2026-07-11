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

	AgentID         string
	ConversationID  string
	IncludeFeedback bool

	PrintFlags *cmdutil.PrintFlags
}

// NewGetCmd returns the `agents conversations get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <agent-id> <conversation-id>",
		Short: "Get a conversation and its messages",
		Args: validators.ExactArgsWithMsg(
			2,
			"agents conversations get requires an <agent-id> and a <conversation-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Get the conversation "conv_123" of the agent "my-agent"
			$ algolia agents conversations get my-agent conv_123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.ConversationID = args[1]
			return runGetCmd(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.IncludeFeedback, "include-feedback", false, "Include feedback for the conversation")

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

	req := client.NewApiGetConversationRequest(opts.ConversationID, opts.AgentID)
	if opts.IncludeFeedback {
		req = req.WithIncludeFeedback(opts.IncludeFeedback)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching conversation")

	res, err := client.GetConversation(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
