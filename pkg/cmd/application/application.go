package application

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/application/add"
	"github.com/algolia/cli/pkg/cmd/application/list"
	"github.com/algolia/cli/pkg/cmd/application/remove"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewApplicationCmd returns a new command for managing applications.
func NewApplicationCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application",
		Aliases: []string{"app", "apps"},
		Short:   "Manage your configured Algolia applications",
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(add.NewAddCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(remove.NewRemoveCmd(f, nil))

	return cmd
}
