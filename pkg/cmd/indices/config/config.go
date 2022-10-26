package config

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/indices/config/export"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConfigCmd returns a new command for indice config management
func NewConfigCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage your Algolia indice config",
	}

	cmd.AddCommand(export.NewExportCmd(f, nil))

	return cmd
}
