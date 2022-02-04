package settings

import (
	"github.com/spf13/cobra"

	importSettings "github.com/algolia/cli/pkg/cmd/settings/import"
	"github.com/algolia/cli/pkg/cmd/settings/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSettingsCmd returns a new command for managing settings.
func NewSettingsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage your Algolia settings",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(importSettings.NewImportCmd(f))

	return cmd
}
