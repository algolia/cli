package config

import (
	"github.com/spf13/cobra"

	configexport "github.com/algolia/cli/pkg/cmd/indices/config/export"
	configimport "github.com/algolia/cli/pkg/cmd/indices/config/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConfigCmd returns a new command for indice config management
func NewConfigCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage your Algolia index config (settings, synonyms, rules)",
	}

	cmd.AddCommand(configexport.NewExportCmd(f))
	cmd.AddCommand(configimport.NewImportCmd(f))

	return cmd
}
