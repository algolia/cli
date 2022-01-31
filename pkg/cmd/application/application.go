package application

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/application/add"
	"github.com/algolia/cli/pkg/cmd/application/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewApplicationCmd returns a new command for managing applications.
func NewApplicationCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "application",
		Short: "Manage your configured Algolia applications",
	}

	cmd.AddCommand(add.NewAddCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f, nil))

	return cmd
}
