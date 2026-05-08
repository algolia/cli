package feedback

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	MessageID         string
	AgentID           string
	Vote              int
	VoteSet           bool
	Tags              []string
	Notes             string
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:     "create --agent-id <id> --message-id <id> --vote 0|1 [--tags x,y] [--notes ...]",
		Aliases: []string{"submit"},
		Short:   "Submit feedback on a single agent message (vote 0=down, 1=up)",
		Example: heredoc.Doc(`
			$ algolia agents feedback create --agent-id <id> --message-id <id> --vote 1
			$ algolia agents feedback create --agent-id <id> --message-id <id> --vote 0 \
			    --tags hallucination,wrong-tone --notes "answered with stale data"
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			opts.VoteSet = cmd.Flags().Changed("vote")
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("--agent-id is required")
			}
			if opts.MessageID == "" {
				return cmdutil.FlagErrorf("--message-id is required")
			}
			if !opts.VoteSet {
				return cmdutil.FlagErrorf("--vote is required (0=down, 1=up)")
			}
			if opts.Vote != 0 && opts.Vote != 1 {
				return cmdutil.FlagErrorf("--vote must be 0 or 1")
			}
			if len(opts.Tags) > 10 {
				return cmdutil.FlagErrorf("--tags may have at most 10 entries (got %d)", len(opts.Tags))
			}
			for _, t := range opts.Tags {
				if len(t) > 50 {
					return cmdutil.FlagErrorf("tag %q exceeds 50-character limit", t)
				}
			}
			if len(opts.Notes) > 1000 {
				return cmdutil.FlagErrorf("--notes exceeds 1000-character limit")
			}
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}
	cmd.Flags().StringVar(&opts.AgentID, "agent-id", "", "Agent that produced the message (required)")
	cmd.Flags().StringVar(&opts.MessageID, "message-id", "", "ID of the assistant message being rated (required)")
	cmd.Flags().IntVar(&opts.Vote, "vote", 0, "Vote: 0=downvote, 1=upvote (required)")
	cmd.Flags().StringSliceVar(&opts.Tags, "tags", nil, "Optional tags (max 10, each <=50 chars)")
	cmd.Flags().StringVar(&opts.Notes, "notes", "", "Optional free-form notes (max 1000 chars)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body := agentstudio.FeedbackCreate{
		MessageID: opts.MessageID,
		AgentID:   opts.AgentID,
		Vote:      opts.Vote,
		Tags:      opts.Tags,
		Notes:     opts.Notes,
	}
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Submitting feedback")
	fb, err := client.CreateFeedback(shared.OrBackground(opts.Ctx), body)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, fb)
}
