package index

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/index/clear"
	"github.com/algolia/cli/pkg/cmd/index/copy"
	"github.com/algolia/cli/pkg/cmd/index/delete"
	"github.com/algolia/cli/pkg/cmd/index/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewIndexCmd returns a new command for indices.
func NewIndexCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Manage your Algolia indices",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(clear.NewClearCmd(f, nil))
	cmd.AddCommand(copy.NewCopyCmd(f, nil))

	return cmd
}
