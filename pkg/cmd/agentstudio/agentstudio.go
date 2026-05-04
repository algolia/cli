package agentstudio

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents"
	"github.com/algolia/cli/pkg/cmd/agentstudio/conversations"
	"github.com/algolia/cli/pkg/cmdutil"
)

const authMethodHelpMsg = `In order to use the 'agentstudio' commands, you need an Algolia application ID and an admin API key. Set them via:
  - The ALGOLIA_APPLICATION_ID and ALGOLIA_API_KEY environment variables
  - Your profile (run 'algolia profile add' or edit ~/.config/algolia/config.tml)`

// NewAgentStudioCmd returns the top-level Agent Studio command.
func NewAgentStudioCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agentstudio",
		Aliases: []string{"agent-studio", "as"},
		Short:   "Manage Algolia Agent Studio agents and conversations",
		Long: heredoc.Docf(`
			Manage Algolia Agent Studio (RAG API) agents, conversations, and completions.

			%s
		`, authMethodHelpMsg),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := f.AgentStudioClient(); err != nil {
				fmt.Fprintf(f.IOStreams.ErrOut, "Agent Studio authentication error: %s\n\n", err)
				fmt.Fprintln(f.IOStreams.ErrOut, authMethodHelpMsg)
				return errors.New("authError")
			}
			return nil
		},
	}

	cmd.AddCommand(agents.NewAgentsCmd(f))
	cmd.AddCommand(conversations.NewConversationsCmd(f))

	return cmd
}
