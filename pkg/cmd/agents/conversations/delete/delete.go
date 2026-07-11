package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// DeleteOptions holds the dependencies and flags for the delete command.
type DeleteOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	ConversationID    string
	DoConfirm         bool
}

// NewDeleteCmd returns the `agents conversations delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <agent-id> <conversation-id>",
		Short: "Delete a conversation",
		Args: validators.ExactArgsWithMsg(
			2,
			"agents conversations delete requires an <agent-id> and a <conversation-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Delete the conversation "conv_123" of the agent "my-agent"
			$ algolia agents conversations delete my-agent conv_123 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.ConversationID = args[1]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete conversation %q?", opts.ConversationID),
					&confirmed,
				)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.DoConfirm, "confirm", "y", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting conversation")

	err = client.DeleteConversation(client.NewApiDeleteConversationRequest(opts.ConversationID, opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted conversation %s\n", cs.SuccessIcon(), opts.ConversationID)
	}

	return nil
}
