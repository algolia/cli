package entries

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/entries/browse"
	"github.com/algolia/cli/pkg/cmd/dictionary/entries/clear"
	"github.com/algolia/cli/pkg/cmd/dictionary/entries/delete"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewEntriesCmd returns a new command for dictionnaries' entries.
func NewEntriesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entries",
		Short: "Manage your Algolia dictionaries entries",
	}

	cmd.AddCommand(clear.NewClearCmd(f, nil))
	cmd.AddCommand(browse.NewBrowseCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))

	return cmd
}
