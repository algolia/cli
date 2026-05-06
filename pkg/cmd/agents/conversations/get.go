package conversations

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID         string
	ConversationID  string
	IncludeFeedback bool
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <agent-id> <conversation-id>",
		Short: "Get a conversation including all its messages",
		Long: heredoc.Doc(`
			Fetch a single conversation by ID, including its full message
			history.

			The response is passed through as JSON because the message
			schema is a discriminated union over role
			(system/user/assistant/tool) with per-role nested content
			arrays — pinning a Go type would silently break on backend
			schema bumps.

			Use --include-feedback to also include the feedback votes
			associated with each assistant message.
		`),
		Example: heredoc.Doc(`
			$ algolia agents conversations get <agent-id> <conv-id>
			$ algolia agents conversations get <agent-id> <conv-id> --include-feedback
		`),
		Args: validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.ConversationID = args[1]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.ConversationID == "" {
				return cmdutil.FlagErrorf("conversation-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVar(&opts.IncludeFeedback, "include-feedback", false, "Include feedback votes attached to assistant messages")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Fetching conversation")
	raw, err := client.GetConversation(ctx, opts.AgentID, opts.ConversationID, opts.IncludeFeedback)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		// Round-trip through generic decode so --output json applies
		// the printer's formatting (consistent indentation regardless
		// of what the backend returned).
		var anyV any
		if err := json.Unmarshal(raw, &anyV); err != nil {
			// Backend returned malformed JSON — surface the bytes
			// rather than crashing the formatter.
			_, _ = opts.IO.Out.Write(raw)
			return nil
		}
		return opts.PrintFlags.Print(opts.IO, anyV)
	}

	// Default human path: pretty-print compactly with 2-space indent.
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		_, _ = opts.IO.Out.Write(raw)
		return nil
	}
	_, _ = opts.IO.Out.Write(pretty.Bytes())
	_, _ = opts.IO.Out.Write([]byte("\n"))
	return nil
}
