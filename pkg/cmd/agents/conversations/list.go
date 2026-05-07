package conversations

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID         string
	Page            int
	PerPage         int
	StartDate       string
	EndDate         string
	IncludeFeedback bool
	FeedbackVote    int
	feedbackVoteSet bool
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:     "list <agent-id>",
		Aliases: []string{"ls"},
		Short:   "List persisted conversations for an agent",
		Long: heredoc.Doc(`
			List conversations stored against an agent.

			Filters:
			  --start-date / --end-date  YYYY-MM-DD passthrough; backend
			                             validates and 422s on bad input.
			  --include-feedback         attaches the feedback array per
			                             conversation (off by default).
			  --feedback-vote 0|1        filters to up- (1) or downvotes (0).
			                             Requires --include-feedback;
			                             enforced at the CLI to match
			                             backend behaviour.
		`),
		Example: heredoc.Doc(`
			$ algolia agents conversations list <agent-id>
			$ algolia agents conversations list <agent-id> --start-date 2026-01-01 --end-date 2026-01-31
			$ algolia agents conversations list <agent-id> --include-feedback --feedback-vote 0 --output json
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			opts.feedbackVoteSet = cmd.Flags().Changed("feedback-vote")
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.feedbackVoteSet {
				if opts.FeedbackVote != 0 && opts.FeedbackVote != 1 {
					return cmdutil.FlagErrorf("--feedback-vote must be 0 (down) or 1 (up)")
				}
				if !opts.IncludeFeedback {
					return cmdutil.FlagErrorf("--feedback-vote requires --include-feedback")
				}
			}
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-indexed; 0 = backend default)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Items per page (0 = backend default, currently 20)")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Filter conversations >= date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Filter conversations <= date (YYYY-MM-DD)")
	cmd.Flags().
		BoolVar(&opts.IncludeFeedback, "include-feedback", false, "Include the feedback array on each conversation")
	cmd.Flags().IntVar(&opts.FeedbackVote, "feedback-vote", 0, "Filter by feedback vote: 0 (down) or 1 (up)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	params := agentstudio.ListConversationsParams{
		Page:            opts.Page,
		Limit:           opts.PerPage,
		StartDate:       opts.StartDate,
		EndDate:         opts.EndDate,
		IncludeFeedback: opts.IncludeFeedback,
	}
	if opts.feedbackVoteSet {
		v := opts.FeedbackVote
		params.FeedbackVote = &v
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching conversations")
	res, err := client.ListConversations(ctx, opts.AgentID, params)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}

	now := nowFn()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("TITLE", nil, nil)
		table.AddField("MSGS", nil, nil)
		table.AddField("TOKENS", nil, nil)
		table.AddField("LAST ACTIVITY", nil, nil)
		table.EndRow()
	}
	for _, conv := range res.Data {
		title := "-"
		if conv.Title != nil && *conv.Title != "" {
			title = *conv.Title
		}
		last := "-"
		if conv.LastActivityAt != nil && !conv.LastActivityAt.IsZero() {
			last = humanize.RelTime(now, *conv.LastActivityAt, "from now", "ago")
		}
		table.AddField(conv.ID, nil, nil)
		table.AddField(title, nil, nil)
		table.AddField(fmt.Sprintf("%d", conv.MessageCount), nil, nil)
		table.AddField(fmt.Sprintf("%d", conv.TotalTokens), nil, nil)
		table.AddField(last, nil, nil)
		table.EndRow()
	}
	if err := table.Render(); err != nil {
		return err
	}
	if table.IsTTY() {
		fmt.Fprintf(opts.IO.Out,
			"\n%d conversation(s) — page %d of %d (total %d).\n",
			len(res.Data),
			res.Pagination.Page, res.Pagination.TotalPages, res.Pagination.TotalCount)
	}
	return nil
}
