package entries

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/entries/clear"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewEntriesCmd returns a new command for dictionnaries' entries.
func NewEntriesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entries",
		Short: "Manage your Algolia dictionaries entries",
	}

	cmd.AddCommand(clear.NewClearCmd(f, nil))

	return cmd
}
