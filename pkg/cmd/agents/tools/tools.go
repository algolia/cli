package tools

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// NewToolsCmd returns `algolia agents tools` — helpers that patch agent
// configuration without authoring full JSON.
func NewToolsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Patch agent tool configuration",
		Long: heredoc.Doc(`
			Convenience commands that fetch an agent, merge structured tool
			fragments (for example Algolia Search index wiring), and PATCH
			the agent. Prefer raw 'algolia agents update -F' when you need
			full control.
		`),
	}
	cmd.AddCommand(newAddSearchIndexCmd(f, nil))
	return cmd
}
