package settings

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/settings/get"
	"github.com/algolia/cli/pkg/cmd/settings/set"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSettingsCmd returns a new command for managing settings.
func NewSettingsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage your Algolia settings",
	}

	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(set.NewSetCmd(f))

	return cmd
}
