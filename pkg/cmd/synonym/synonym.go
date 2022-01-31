package synonym

import (
	"github.com/spf13/cobra"

	importSynonyms "github.com/algolia/cli/pkg/cmd/synonym/import"
	"github.com/algolia/cli/pkg/cmd/synonym/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSynonymCmd returns a new command for synonyms.
func NewSynonymCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synonym",
		Short: "Manage your Algolia synonyms",
	}

	cmd.AddCommand(importSynonyms.NewImportCmd(f))
	cmd.AddCommand(list.NewListCmd(f))

	return cmd
}
