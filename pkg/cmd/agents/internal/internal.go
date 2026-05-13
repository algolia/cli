package internal

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// NewInternalCmd is the parent for `algolia agents internal <verb>` —
// hidden, unstable, x-hidden endpoints. See docs/agents.md.
func NewInternalCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "internal",
		Short:  "Hidden Agent Studio diagnostics + memory internals (unstable)",
		Hidden: true,
		Long: heredoc.Doc(`
			Wraps x-hidden Agent Studio endpoints. Not part of the
			documented public surface; may change without notice.
			memorize/ponder/consolidate hit the doubled path
			/1/agents/agents/{id}/<verb>.
		`),
	}
	cmd.AddCommand(newStatusCmd(f, nil))
	cmd.AddCommand(newMemorizeCmd(f, nil))
	cmd.AddCommand(newPonderCmd(f, nil))
	cmd.AddCommand(newConsolidateCmd(f, nil))
	return cmd
}
