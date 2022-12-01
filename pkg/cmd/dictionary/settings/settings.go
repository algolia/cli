package settings

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/settings/get"
	"github.com/algolia/cli/pkg/cmd/dictionary/settings/set"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSettingsCmd returns a new command for dictionnaries' entries.
func NewSettingsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage your Algolia dictionaries settings",
	}

	cmd.AddCommand(set.NewSetCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))

	return cmd
}
