package indices

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/indices/clear"
	"github.com/algolia/cli/pkg/cmd/indices/copy"
	"github.com/algolia/cli/pkg/cmd/indices/delete"
	"github.com/algolia/cli/pkg/cmd/indices/list"
	"github.com/algolia/cli/pkg/cmd/indices/move"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewIndicesCmd returns a new command for indices.
func NewIndicesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "indices",
		Aliases: []string{"index"},
		Short:   "Manage your Algolia indices",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(clear.NewClearCmd(f, nil))
	cmd.AddCommand(copy.NewCopyCmd(f, nil))
	cmd.AddCommand(move.NewMoveCmd(f, nil))

	return cmd
}
