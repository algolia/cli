package config

import (
	"github.com/spf13/cobra"

	indexexport "github.com/algolia/cli/pkg/cmd/indices/config/export"
	indeximport "github.com/algolia/cli/pkg/cmd/indices/config/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConfigCmd returns a new command for indice config management
func NewConfigCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage your Algolia index config (settings, synonyms, rules)",
	}

	cmd.AddCommand(indexexport.NewExportCmd(f))
	cmd.AddCommand(indeximport.NewImportCmd(f))

	return cmd
}
