package conversations

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID        string
	ConversationID string
	DoConfirm      bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <agent-id> <conversation-id> [--confirm]",
		Short: "Delete a single conversation",
		Long: heredoc.Doc(`
			Delete one conversation by ID.

			This is the surgical sibling of "purge". Mistyping a
			conversation ID here at worst nukes one extra conversation;
			use this verb when you have a specific conversation in mind
			and "purge" when you want to clean a range.

			Like "agents delete", interactive use prompts to confirm and
			non-interactive use requires --confirm.
		`),
		Example: heredoc.Doc(`
			$ algolia agents conversations delete <agent-id> <conv-id>
			$ algolia agents conversations delete <agent-id> <conv-id> -y
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
			doConfirm, err := shared.ResolveConfirm(opts.IO, confirm)
			if err != nil {
				return err
			}
			opts.DoConfirm = doConfirm
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}

	shared.AddConfirmFlag(cmd, &confirm)
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	if opts.DoConfirm {
		ok, err := shared.Confirm(
			fmt.Sprintf("Delete conversation %s on agent %s?", opts.ConversationID, opts.AgentID),
		)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Deleting conversation")
	err = client.DeleteConversation(ctx, opts.AgentID, opts.ConversationID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted conversation %s\n", cs.SuccessIcon(), opts.ConversationID)
	}
	return nil
}
