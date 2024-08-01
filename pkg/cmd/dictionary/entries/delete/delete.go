package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Dictionary search.DictionaryType
	ObjectIDs  []string

	DoConfirm bool
}

// NewDeleteCmd deletes and returns a delete command for dictionaries entries.
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4_SearchClient,
	}
	cmd := &cobra.Command{
		Use:       "delete <dictionary> --object-ids <object-ids> [--confirm]",
		Args:      validators.ExactArgs(1),
		ValidArgs: shared.DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return shared.DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Annotations: map[string]string{
			"acls": "settings,editSettings",
		},
		Short: "Delete dictionary entries",
		Long: heredoc.Docf(`
			This command deletes entries from the specified dictionary.
		`),
		Example: heredoc.Doc(`
			# Delete one single entry with the ID "1" from the "plurals" dictionary
			$ algolia dictionary entries delete plurals --object-ids 1

			# Delete multiple entries with the IDs "1" and "2" from the "plurals" dictionary
			$ algolia dictionary entries delete plurals --object-ids 1,2
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Dictionary = search.DictionaryType(args[0])
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().
		StringSliceVarP(&opts.ObjectIDs, "object-ids", "", nil, "Object IDs of dictionary entries to delete")
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

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(
			fmt.Sprintf(
				"Delete %s from %s?",
				pluralizeEntry(len(opts.ObjectIDs)),
				opts.Dictionary,
			),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	// Construct batch request
	var requests []search.BatchDictionaryEntriesRequest

	for _, id := range opts.ObjectIDs {
		requests = append(
			requests,
			*search.NewBatchDictionaryEntriesRequest(search.DICTIONARY_ACTION_DELETE_ENTRY, *search.NewDictionaryEntry(id)),
		)
	}

	_, err = client.BatchDictionaryEntries(
		client.NewApiBatchDictionaryEntriesRequest(
			search.DictionaryType(opts.Dictionary),
			search.NewBatchDictionaryEntriesParams(requests),
		),
	)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Successfully deleted %s from %s\n",
			cs.SuccessIcon(),
			pluralizeEntry(len(opts.ObjectIDs)),
			opts.Dictionary,
		)
	}

	return nil
}

func pluralizeEntry(count int) string {
	if count <= 1 {
		return fmt.Sprintf("%d entry", count)
	}
	return fmt.Sprintf("%d entries", count)
}
