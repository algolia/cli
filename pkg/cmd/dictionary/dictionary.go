package dictionary

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/entries"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewDictionaryCmd returns a new command for dictionnaries.
func NewDictionaryCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dictionary",
		Aliases: []string{"dictionaries"},
		Short:   "Manage your Algolia dictionaries",
	}

	cmd.AddCommand(entries.NewEntriesCmd(f))

	return cmd
}
