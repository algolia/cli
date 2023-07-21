package dictionary

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/entries"
	"github.com/algolia/cli/pkg/cmd/dictionary/settings"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewDictionaryCmd returns a new command for dictionaries.
func NewDictionaryCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dictionary",
		Aliases: []string{"dictionaries"},
		Short:   "Manage your Algolia dictionaries",
		Annotations: map[string]string{
			"help:see-also": heredoc.Doc(`
			The below command examples are not directly related to the dictionary command, but are relevant to the use of dictionaries in general.

			# Open Algolia supported languages page
			$ algolia open languages

			# Set the 'ignorePlurals' setting to true: https://www.algolia.com/doc/api-reference/api-parameters/ignorePlurals/
			$ algolia settings set <index> --ignorePlurals

			# Set the 'removeStopWords' setting to true: https://www.algolia.com/doc/api-reference/api-parameters/removeStopWords/
			$ algolia settings set <index> --removeStopWords

			# Set the 'decompoundQuery' setting to true: https://www.algolia.com/doc/api-reference/api-parameters/decompoundQuery/
			$ algolia settings set <index> --decompoundQuery
			`),
		},
	}

	cmd.AddCommand(settings.NewSettingsCmd(f))
	cmd.AddCommand(entries.NewEntriesCmd(f))

	return cmd
}
