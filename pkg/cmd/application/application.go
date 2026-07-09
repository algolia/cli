package application

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/application/create"
	"github.com/algolia/cli/pkg/cmd/application/downgrade"
	"github.com/algolia/cli/pkg/cmd/application/list"
	"github.com/algolia/cli/pkg/cmd/application/plans"
	"github.com/algolia/cli/pkg/cmd/application/selectapp"
	"github.com/algolia/cli/pkg/cmd/application/update"
	"github.com/algolia/cli/pkg/cmd/application/upgrade"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewApplicationCmd returns a new command for managing applications.
func NewApplicationCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application",
		Aliases: []string{"app"},
		Short:   "Create, select, and manage your Algolia applications",
	}

	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(selectapp.NewSelectCmd(f))
	cmd.AddCommand(update.NewUpdateCmd(f))
	cmd.AddCommand(plans.NewPlansCmd(f))
	cmd.AddCommand(upgrade.NewUpgradeCmd(f))
	cmd.AddCommand(downgrade.NewDowngradeCmd(f))

	return cmd
}
