package conversations

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID        string
	ConversationID string
	DryRun         bool
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
			non-interactive use requires --confirm. --dry-run previews.
		`),
		Example: heredoc.Doc(`
			$ algolia agents conversations delete <agent-id> <conv-id>
			$ algolia agents conversations delete <agent-id> <conv-id> -y
			$ algolia agents conversations delete <agent-id> <conv-id> --dry-run
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
			if !confirm && !opts.DryRun {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				opts.DoConfirm = true
			}
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Print what would be deleted without calling the API")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out,
			"Dry run: would DELETE /1/agents/%s/conversations/%s\n",
			opts.AgentID, opts.ConversationID)
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf("Delete conversation %s on agent %s?", opts.ConversationID, opts.AgentID),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

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
