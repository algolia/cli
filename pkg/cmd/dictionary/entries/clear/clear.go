package clear

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

type ClearOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Dictionnaries []search.DictionaryName
	All           bool

	DoConfirm bool
}

// EntryType represents the type of an entry in a dictionnary.
// It can be either a custom entry or a standard entry.
type EntryType string

// DictionaryEntry is a simple type alias for the search.DictionaryEntry type (which do not include the type of the entry).
type DictionaryEntry struct {
	Type EntryType
}

const (
	// CustomEntryType is the type of a custom entry in a dictionnary (i.e. added by the user).
	CustomEntryType EntryType = "custom"
	// StandardEntryType is the type of a standard entry in a dictionnary (i.e. added by Algolia).
	StandardEntryType EntryType = "standard"
)

var (
	// DictionaryNames returns the list of available dictionnaries.
	DictionaryNames = func() []string {
		return []string{
			string(search.Stopwords),
			string(search.Compounds),
			string(search.Plurals),
		}
	}
)

// NewClearCmd creates and returns a clear command for dictionnaries' entries.
func NewClearCmd(f *cmdutil.Factory, runF func(*ClearOptions) error) *cobra.Command {
	var confirm bool
	cs := f.IOStreams.ColorScheme()

	opts := &ClearOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:       "clear {<dictionnary>... | --all} [--confirm]",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Clear dictionary entries",
		Long: heredoc.Docf(`
			This command deletes all entries from the specified %s dictionnaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Delete all entries from the "stopword" dictionnary
			$ algolia dictionary entries clear stopword

			# Delete all entries from the "stopword" and "plural" dictionnaries
			$ algolia dictionary entries clear stopword plural

			# Delete all entries from all dictionnaries
			$ algolia dictionary entries clear --all
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.All && len(args) > 0 || !opts.All && len(args) == 0 {
				return cmdutil.FlagErrorf("Either specify dictionnaries' names or use --all to clear all dictionnaries")
			}

			if opts.All {
				opts.Dictionnaries = []search.DictionaryName{search.Stopwords, search.Plurals, search.Compounds}
			} else {
				opts.Dictionnaries = make([]search.DictionaryName, len(args))
				for i, dictionnary := range args {
					opts.Dictionnaries[i] = search.DictionaryName(dictionnary)
				}
			}

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
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
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "clear all dictionnaries")

	return cmd
}

// runClearCmd executes the clear command
func runClearCmd(opts *ClearOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	dictionaries := opts.Dictionnaries
	dictionariesNames := make([]string, len(dictionaries))
	dictionariesCustomEntriesNb := make([]int, len(dictionaries))
	for i, dictionnary := range dictionaries {
		nbCustomEntries, err := customEntriesNb(client, dictionnary)
		if err != nil {
			return err
		}
		dictionariesCustomEntriesNb[i] = nbCustomEntries
		dictionariesNames[i] = string(dictionnary)
	}

	totalEntries := 0
	for _, nb := range dictionariesCustomEntriesNb {
		totalEntries += nb
	}

	if totalEntries == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No entries to clear in %s dictionary.\n", cs.WarningIcon(), utils.SliceToReadableString(dictionariesNames))
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Clear %d entries from %s dictionary?", totalEntries, utils.SliceToReadableString(dictionariesNames)), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	for _, dictionary := range dictionaries {
		_, err = client.ClearDictionaryEntries(dictionary)
		if err != nil {
			return err
		}
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully cleared %d entries from %s dictionary\n", cs.SuccessIcon(), totalEntries, utils.SliceToReadableString(dictionariesNames))
	}

	return nil
}

func customEntriesNb(client *search.Client, dictionnary search.DictionaryName) (int, error) {
	res, err := client.SearchDictionaryEntries(dictionnary, "", opt.HitsPerPage(1000))
	if err != nil {
		return 0, err
	}
	data, err := json.Marshal(res.Hits)
	if err != nil {
		return 0, fmt.Errorf("cannot unmarshal dictionary entries: error while marshalling original dictionary entries: %v", err)
	}

	var entries []DictionaryEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return 0, fmt.Errorf("cannot unmarshal dictionary entries: error while unmarshalling original dictionary entries: %v", err)
	}

	var customEntriesNb int
	for _, entry := range entries {
		if entry.Type == CustomEntryType {
			customEntriesNb++
		}
	}

	return customEntriesNb, nil
}
