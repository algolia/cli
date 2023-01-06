package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Dictionnary search.DictionaryName
	ObjectIDs   []string

	DoConfirm bool
}

// DictionaryEntry can be plural, compound or stopword entry.
type DictionaryEntry struct {
	Type          shared.EntryType
	Word          string   `json:"word,omitempty"`
	Words         []string `json:"words,omitempty"`
	Decomposition string   `json:"decomposition,omitempty"`
	ObjectID      string
	Language      string
}

// NewDeleteCmd deletes and returns a delete command for dictionaries entries.
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:       "delete <dictionary> --object-ids <object-ids> --confirm",
		Args:      validators.ExactArgs(1),
		ValidArgs: shared.DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return shared.DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Delete dictionary entries",
		Long: heredoc.Docf(`
			This command deletes entries from the specified dictionary.
		`),
		Example: heredoc.Doc(`
			# Delete one single object with the ID "1" from the "plural" dictionary
			$ algolia dictionary entries delete plural --object-ids 1

			# Delete multiple objects with the IDs "1" and "2" from the "plural" index
			$ algolia dictionary entries delete plural --object-ids 1,2
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Dictionnary = search.DictionaryName(args[0])
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.ObjectIDs, "object-ids", "", nil, "Object IDs to delete")
	_ = cmd.MarkFlagRequired("object-ids")

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

// runDeleteCmd executes the delete command
func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	dictionary := opts.Dictionnary
	objectIDs := opts.ObjectIDs

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Delete the %s from %s?", utils.Pluralize(len(objectIDs), "object"), dictionary), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	_, err = client.DeleteDictionaryEntries(dictionary, objectIDs)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully deleted entries from %s\n", cs.SuccessIcon(), dictionary)
	}

	return nil
}
