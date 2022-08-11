package synonyms

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/synonyms/browse"
	importSynonyms "github.com/algolia/cli/pkg/cmd/synonyms/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSynonymsCmd returns a new command for synonyms.
func NewSynonymsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "synonyms",
		Aliases: []string{"synonym"},
		Short:   "Manage your Algolia synonyms",
	}

	cmd.AddCommand(importSynonyms.NewImportCmd(f, nil))
	cmd.AddCommand(browse.NewBrowseCmd(f))

	return cmd
}
