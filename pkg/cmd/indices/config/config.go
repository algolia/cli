package config

import (
	"github.com/spf13/cobra"

	indiceexport "github.com/algolia/cli/pkg/cmd/indices/config/export"
	indiceimport "github.com/algolia/cli/pkg/cmd/indices/config/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConfigCmd returns a new command for indice config management
func NewConfigCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage your Algolia indice config",
	}

	cmd.AddCommand(indiceexport.NewExportCmd(f))
	cmd.AddCommand(indiceimport.NewImportCmd(f))

	return cmd
}
