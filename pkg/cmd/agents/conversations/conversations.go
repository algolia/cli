package conversations

import (
	"context"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// nowFn is overridable for deterministic time-based output in tests.
var nowFn = time.Now

// NewConversationsCmd is the parent for `algolia agents conversations <verb>`.
//
// Conversations are persisted by the backend whenever an end-user
// interaction reaches a published agent (or `agents try`/`agents run`
// from the CLI). They live PER AGENT — every endpoint takes
// {agent_id} in the path — so every verb here takes the agent ID as a
// positional first argument, matching `agents publish/run/cache invalidate`.
//
// Five verbs, split per file (list/get/delete/purge/export). The
// single-record `delete` and the bulk `purge` are intentionally
// different command names: same HTTP method, vastly different blast
// radius. Keeping them apart makes the dangerous operation impossible
// to invoke by accident with a typo on the conversation ID.
func NewConversationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversations",
		Short: "Inspect and manage Agent Studio conversations",
		Long: heredoc.Doc(`
			Inspect, delete, and export conversations persisted by Agent
			Studio for a given agent. All verbs take the agent ID as the
			first positional argument.

			Common workflows:

			  - "agents conversations list <agent-id>" — page through
			    persisted conversations, optionally filter by date range
			    or feedback vote.
			  - "agents conversations get <agent-id> <conv-id>" —
			    retrieve a single conversation including all its messages.
			  - "agents conversations delete <agent-id> <conv-id>" —
			    delete one conversation.
			  - "agents conversations purge <agent-id>" — bulk delete
			    (requires --all OR a date range).
			  - "agents conversations export <agent-id>" — dump every
			    matching conversation to stdout or a file.
		`),
	}

	cmd.AddCommand(newListCmd(f, nil))
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	cmd.AddCommand(newPurgeCmd(f, nil))
	cmd.AddCommand(newExportCmd(f, nil))
	return cmd
}

// ctxOrBackground promotes a possibly-nil command Context to
// context.Background. Cobra always supplies one in production, but
// table-test invocations of run* helpers occasionally don't.
func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
