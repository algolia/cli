package conversations

import (
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

var nowFn = time.Now

// NewConversationsCmd is the parent for `algolia agents conversations <verb>`.
// All verbs take the agent ID as the first positional argument.
func NewConversationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversations",
		Short: "Inspect and manage Agent Studio conversations",
		Long: heredoc.Doc(`
			Inspect, delete, and export conversations persisted by Agent
			Studio for a given agent. All verbs take the agent ID as the
			first positional argument. ` + "`purge`" + ` requires a date
			range (see docs/agents.md).
		`),
	}

	cmd.AddCommand(newListCmd(f, nil))
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	cmd.AddCommand(newPurgeCmd(f, nil))
	cmd.AddCommand(newExportCmd(f, nil))
	return cmd
}
