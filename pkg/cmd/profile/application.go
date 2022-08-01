package profile

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/profile/add"
	"github.com/algolia/cli/pkg/cmd/profile/list"
	"github.com/algolia/cli/pkg/cmd/profile/remove"
	"github.com/algolia/cli/pkg/cmd/profile/setdefault"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewProfileCmd returns a new command for managing profiles.
func NewProfileCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"profiles"},
		Short:   "Manage your profiles",
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(add.NewAddCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(remove.NewRemoveCmd(f, nil))
	cmd.AddCommand(setdefault.NewSetDefaultCmd(f, nil))

	return cmd
}
