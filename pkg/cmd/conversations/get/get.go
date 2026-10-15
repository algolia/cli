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
	PrintFlags        *cmdutil.PrintFlags

	AgentID                string
	ConversationID         string
	IncludeFeedback        bool
	IncludeMessageEvents   bool
	IncludeImpactAnalytics bool
	SetFeedback            bool
	SetMessageEvents       bool
	SetImpactAnalytics     bool
}

// NewGetCmd returns the `conversations get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <agent-id> <conversation-id>",
		Short: "Get an Agent Studio conversation including its messages",
		Args: validators.ExactArgsWithMsg(
			2,
			"agentstudio conversations get requires <agent-id> and <conversation-id> arguments.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Get a conversation
			$ algolia agentstudio conversations get 11111111-1111-4111-8111-111111111111 33333333-3333-4333-8333-333333333333
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.ConversationID = args[1]
			opts.SetFeedback = cmd.Flags().Changed("include-feedback")
			opts.SetMessageEvents = cmd.Flags().Changed("include-message-events")
			opts.SetImpactAnalytics = cmd.Flags().Changed("include-impact-analytics")
			return runGetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVar(&opts.IncludeFeedback, "include-feedback", false, "Include feedback for the conversation")
	cmd.Flags().
		BoolVar(&opts.IncludeMessageEvents, "include-message-events", false, "Include Insights events on assistant messages")
	cmd.Flags().
		BoolVar(&opts.IncludeImpactAnalytics, "include-impact-analytics", false, "Include hasView/hasClick/hasConversion")

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
	if opts.SetFeedback {
		req = req.WithIncludeFeedback(opts.IncludeFeedback)
	}
	if opts.SetMessageEvents {
		req = req.WithIncludeMessageEvents(opts.IncludeMessageEvents)
	}
	if opts.SetImpactAnalytics {
		req = req.WithIncludeImpactAnalytics(opts.IncludeImpactAnalytics)
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
