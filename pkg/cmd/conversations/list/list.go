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
	PrintFlags        *cmdutil.PrintFlags

	AgentID         string
	Page            int32
	Limit           int32
	StartDate       string
	EndDate         string
	IncludeFeedback bool
	SetFeedback     bool
}

// NewListCmd returns the `conversations list` command.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "list <agent-id>",
		Short: "List Agent Studio conversations for an agent",
		Args: validators.ExactArgsWithMsg(
			1,
			"agentstudio conversations list requires a <agent-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# List conversations for an agent
			$ algolia agentstudio conversations list 11111111-1111-4111-8111-111111111111

			# Filter by date range
			$ algolia agentstudio conversations list 11111111-1111-4111-8111-111111111111 --start-date 2026-01-01 --end-date 2026-01-31
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.SetFeedback = cmd.Flags().Changed("include-feedback")
			return runListCmd(opts)
		},
	}

	cmd.Flags().Int32Var(&opts.Page, "page", 0, "Page number")
	cmd.Flags().Int32Var(&opts.Limit, "limit", 0, "Items per page")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Created after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Created before this date (YYYY-MM-DD)")
	cmd.Flags().
		BoolVar(&opts.IncludeFeedback, "include-feedback", false, "Include feedback per conversation")

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

	req := client.NewApiListAgentConversationsRequest(opts.AgentID)
	if opts.Page > 0 {
		req = req.WithPage(opts.Page)
	}
	if opts.Limit > 0 {
		req = req.WithLimit(opts.Limit)
	}
	if opts.StartDate != "" {
		req = req.WithStartDate(opts.StartDate)
	}
	if opts.EndDate != "" {
		req = req.WithEndDate(opts.EndDate)
	}
	if opts.SetFeedback {
		req = req.WithIncludeFeedback(opts.IncludeFeedback)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching conversations")

	res, err := client.ListAgentConversations(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
