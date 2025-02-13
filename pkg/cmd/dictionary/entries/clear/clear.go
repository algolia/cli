package clear

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
)

type ClearOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Dictionaries []search.DictionaryType
	All          bool

	DoConfirm bool
}

// NewClearCmd creates and returns a clear command for dictionaries' entries.
func NewClearCmd(f *cmdutil.Factory, runF func(*ClearOptions) error) *cobra.Command {
	var confirm bool
	cs := f.IOStreams.ColorScheme()

	opts := &ClearOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}
	cmd := &cobra.Command{
		Use:       "clear {<dictionary>... | --all} [--confirm]",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: shared.DictionaryTypes(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return shared.DictionaryTypes(), cobra.ShellCompDirectiveNoFileComp
		},
		Annotations: map[string]string{
			"acls": "settings,editSettings",
		},
		Short: "Clear dictionary entries",
		Long: heredoc.Docf(`
			This command deletes all entries from the specified %s dictionaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Delete all entries from the "stopwords" dictionary
			$ algolia dictionary entries clear stopwords

			# Delete all entries from the "stopwords" and "plurals" dictionaries
			$ algolia dictionary entries clear stopwords plurals

			# Delete all entries from all dictionaries
			$ algolia dictionary entries clear --all
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.All && len(args) > 0 || !opts.All && len(args) == 0 {
				return cmdutil.FlagErrorf(
					"Either specify dictionaries' names or use --all to clear all dictionaries",
				)
			}

			if opts.All {
				opts.Dictionaries = search.AllowedDictionaryTypeEnumValues
			} else {
				opts.Dictionaries = make([]search.DictionaryType, len(args))
				for i, dictionary := range args {
					opts.Dictionaries[i] = search.DictionaryType(dictionary)
				}
			}

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

			return runClearCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "clear all dictionaries")

	return cmd
}

// runClearCmd executes the clear command
func runClearCmd(opts *ClearOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	dictionaries := opts.Dictionaries
	dictionariesNames := make([]string, len(dictionaries))
	dictionariesCustomEntriesNb := make([]int, len(dictionaries))
	for i, dictionary := range dictionaries {
		nbCustomEntries, err := customEntriesNb(client, dictionary)
		if err != nil {
			return err
		}
		dictionariesCustomEntriesNb[i] = nbCustomEntries
		dictionariesNames[i] = string(dictionary)
	}

	totalEntries := 0
	for _, nb := range dictionariesCustomEntriesNb {
		totalEntries += nb
	}

	if totalEntries == 0 {
		if _, err = fmt.Fprintf(opts.IO.Out, "%s No entries to clear in %s dictionary.\n", cs.WarningIcon(), utils.SliceToReadableString(dictionariesNames)); err != nil {
			return err
		}
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(
			fmt.Sprintf(
				"Clear %d entries from %s dictionary?",
				totalEntries,
				utils.SliceToReadableString(dictionariesNames),
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

	for _, dict := range dictionaries {
		_, err = client.BatchDictionaryEntries(
			client.NewApiBatchDictionaryEntriesRequest(
				dict,
				search.NewEmptyBatchDictionaryEntriesParams().
					SetClearExistingDictionaryEntries(true),
			),
		)
		if err != nil {
			return err
		}
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Successfully cleared %d entries from %s dictionary\n",
			cs.SuccessIcon(),
			totalEntries,
			utils.SliceToReadableString(dictionariesNames),
		)
	}

	return nil
}

func customEntriesNb(client *search.APIClient, dictionary search.DictionaryType) (int, error) {
	res, err := client.SearchDictionaryEntries(
		client.NewApiSearchDictionaryEntriesRequest(
			dictionary,
			search.NewEmptySearchDictionaryEntriesParams().SetHitsPerPage(1000).SetQuery(""),
		),
	)
	if err != nil {
		return 0, err
	}
	data, err := json.Marshal(res.Hits)
	if err != nil {
		return 0, fmt.Errorf(
			"cannot unmarshal dictionary entries: error while marshalling original dictionary entries: %v",
			err,
		)
	}

	var entries []search.DictionaryEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return 0, fmt.Errorf(
			"cannot unmarshal dictionary entries: error while unmarshalling original dictionary entries: %v",
			err,
		)
	}

	var customEntriesNb int
	for _, entry := range entries {
		if entry.Type != nil && *entry.Type == search.DICTIONARY_ENTRY_TYPE_CUSTOM {
			customEntriesNb++
		}
	}

	return customEntriesNb, nil
}
