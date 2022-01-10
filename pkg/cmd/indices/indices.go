package indices

import (
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmd/indices/clear"
	"github.com/algolia/algolia-cli/pkg/cmd/indices/delete"
	"github.com/algolia/algolia-cli/pkg/cmd/indices/dump"
	"github.com/algolia/algolia-cli/pkg/cmd/indices/list"
	"github.com/algolia/algolia-cli/pkg/cmd/indices/load"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
)

// NewIndicesCmd returns a new command for indices.
func NewIndicesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indices",
		Short: "Manage your Algolia indices",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(clear.NewClearCmd(f))
	cmd.AddCommand(dump.NewDumpCmd(f))
	cmd.AddCommand(load.NewLoadCmd(f))

	return cmd
}
