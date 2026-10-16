package agentstudio

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentStudioCmd returns the agentstudio command group.
func NewAgentStudioCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agentstudio",
		Short: "Manage Algolia Agent Studio",
		Long:  "Manage Algolia Agent Studio resources.",
	}

	cmd.AddCommand(agents.NewAgentsCmd(f))
	return cmd
}
