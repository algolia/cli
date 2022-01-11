package synonyms

import (
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmd/synonyms/export"
	importSynonyms "github.com/algolia/algolia-cli/pkg/cmd/synonyms/import"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
)

// NewSynonymsCmd returns a new command for synonyms.
func NewSynonymsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synonyms",
		Short: "Manage your Algolia synonyms",
	}

	cmd.AddCommand(importSynonyms.NewImportCmd(f))
	cmd.AddCommand(export.NewExportCmd(f))

	return cmd
}
