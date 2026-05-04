package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID        string
	ConversationID string
	All            bool
	Confirm        bool
}

func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}
	cmd := &cobra.Command{
		Use:     "delete <agent_id> [conversation_id]",
		Aliases: []string{"rm"},
		Args:    cobra.RangeArgs(1, 2),
		Short:   "Delete a conversation, or all conversations for an agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Delete a single conversation
			$ algolia agentstudio conversations delete a1b2 conv-1 --confirm

			# Delete every conversation for the agent
			$ algolia agentstudio conversations delete a1b2 --all --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			if len(args) == 2 {
				opts.ConversationID = args[1]
			}
			if !opts.All && opts.ConversationID == "" {
				return fmt.Errorf("either pass <conversation_id> or use --all")
			}
			if opts.All && opts.ConversationID != "" {
				return fmt.Errorf("--all conflicts with a positional <conversation_id>")
			}
			if !opts.Confirm {
				return fmt.Errorf("--confirm is required to delete")
			}
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.All, "all", false, "Delete every conversation for the agent")
	cmd.Flags().BoolVar(&opts.Confirm, "confirm", false, "Skip the confirmation prompt and delete immediately")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Deleting conversation(s)")
	if opts.All {
		err = client.DeleteAllConversations(opts.AgentID)
	} else {
		err = client.DeleteConversation(opts.AgentID, opts.ConversationID)
	}
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		if opts.All {
			fmt.Fprintf(opts.IO.Out, "%s Deleted all conversations for agent %s\n", cs.SuccessIcon(), opts.AgentID)
		} else {
			fmt.Fprintf(opts.IO.Out, "%s Deleted conversation %s\n", cs.SuccessIcon(), opts.ConversationID)
		}
	}
	return nil
}
