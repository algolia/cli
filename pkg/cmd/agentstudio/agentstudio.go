package agentstudio

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers"
	"github.com/algolia/cli/pkg/cmd/agentstudio/secretkeys"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentStudioCmd returns the top-level agent-studio command group.
func NewAgentStudioCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent-studio",
		Short: "Manage your Algolia Agent Studio resources",
		Long:  "Manage Algolia Agent Studio agents, LLM providers, and secret keys.",
	}

	cmd.AddCommand(agents.NewAgentsCmd(f))
	cmd.AddCommand(providers.NewProvidersCmd(f))
	cmd.AddCommand(secretkeys.NewSecretKeysCmd(f))

	return cmd
}
