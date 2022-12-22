package browse

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
	"github.com/algolia/cli/pkg/utils"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Dictionnaries        []search.DictionaryName
	All                  bool
	ShowDefaultStopwords bool

	PrintFlags *cmdutil.PrintFlags
}

// EntryType represents the type of an entry in a dictionnary.
// It can be either a custom entry or a standard entry.
type EntryType string
type DictionaryType int

// DictionaryEntry can be plural, compound or stopword entry.
type DictionaryEntry struct {
	Type          EntryType
	Word          string   `json:"word,omitempty"`
	Words         []string `json:"words,omitempty"`
	Decomposition string   `json:"decomposition,omitempty"`
	ObjectID      string
	Language      string
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

// NewBrowseCmd creates and returns a browse command for dictionnaries' entries.
func NewBrowseCmd(f *cmdutil.Factory, runF func(*BrowseOptions) error) *cobra.Command {
	cs := f.IOStreams.ColorScheme()

	opts := &BrowseOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:       "browse {<dictionary>... | --all} [--include-defaults]",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Browse dictionary entries",
		Long: heredoc.Docf(`
			This command retrieves all entries from the specified %s dictionnaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Retrieve all entries from the "stopword" dictionary (doesn't include default stopwords)
			$ algolia dictionary entries browse stopword

			# Retrieve all entries from the "stopword" and "plural" dictionnaries
			$ algolia dictionary entries browse stopword plural

			# Retrieve all entries from all dictionnaries
			$ algolia dictionary entries browse --all

			# Retrieve all entries from the "stopword" dictionnaries (including default stopwords)
			$ algolia dictionary entries browse stopword --showDefaultStopwords
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.All && len(args) > 0 || !opts.All && len(args) == 0 {
				return cmdutil.FlagErrorf("Either specify dictionaries' names or use --all to browse all dictionaries")
			}

			if opts.All {
				opts.Dictionnaries = []search.DictionaryName{search.Stopwords, search.Plurals, search.Compounds}
			} else {
				opts.Dictionnaries = make([]search.DictionaryName, len(args))
				for i, dictionnary := range args {
					opts.Dictionnaries[i] = search.DictionaryName(dictionnary)
				}
			}

			return runBrowseCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "browse all dictionnaries")
	cmd.Flags().BoolVar(&opts.ShowDefaultStopwords, "showDefaultStopwords", false, "browse dictionaries and include default stopwords")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runBrowseCmd executes the browse command
func runBrowseCmd(opts *BrowseOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	dictionaries := opts.Dictionnaries
	dictionariesNames := make([]string, len(dictionaries))

	dictionariesCustomEntriesNb := make([]int, len(dictionaries))
	dictionariesAllEntriesNb := make([]int, len(dictionaries))

	dictionariesCustomEntries := make([][]DictionaryEntry, len(dictionaries))
	dictionariesAllEntries := make([][]DictionaryEntry, len(dictionaries))

	for i, dictionnary := range dictionaries {
		// get all entries, custom entries, and number of both entries
		allEntries, customEntries, nbAllEntries, nbCustomEntries, err := dictionaryEntriesData(client, dictionnary)
		if err != nil {
			return err
		}
		dictionariesCustomEntriesNb[i] = nbCustomEntries
		dictionariesAllEntriesNb[i] = nbAllEntries

		dictionariesNames[i] = string(dictionnary)

		dictionariesCustomEntries[i] = customEntries
		dictionariesAllEntries[i] = allEntries
	}

	totalCustomEntries := 0
	totalAllEntries := 0

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	// incase the `--showDefaultDictionary` flag is set
	if opts.ShowDefaultStopwords {
		for _, nb := range dictionariesAllEntriesNb {
			totalAllEntries += nb
		}
		if totalAllEntries == 0 {
			fmt.Fprintf(opts.IO.Out, "%s No entries in %s dictionary.\n", cs.WarningIcon(), utils.SliceToReadableString(dictionariesNames))
			return nil
		}
		for _, dictHits := range dictionariesAllEntries {
			// print dictionary entries
			p.Print(opts.IO, dictHits)
		}
	} else {
		for _, nb := range dictionariesCustomEntriesNb {
			totalCustomEntries += nb
		}
		if totalCustomEntries == 0 {
			fmt.Fprintf(opts.IO.Out, "%s No custom entries in %s dictionary.\n", cs.WarningIcon(), utils.SliceToReadableString(dictionariesNames))
			return nil
		}
		for _, dictHits := range dictionariesCustomEntries {
			// print dictionary entries
			p.Print(opts.IO, dictHits)
		}
	}

	return nil
}

// returns all entries, custom entries, the number of entries and the number of custom entries
func dictionaryEntriesData(client *search.Client, dictionnary search.DictionaryName) ([]DictionaryEntry, []DictionaryEntry, int, int, error) {

	pageCount := 0
	maxPages := 1

	allEntries := make([]DictionaryEntry, 0)
	customEntries := make([]DictionaryEntry, 0)
	var customEntriesNb int
	var allEntriesNb int

	// implement infinite pagination
	for pageCount < maxPages {
		res, err := client.SearchDictionaryEntries(dictionnary, "", opt.HitsPerPage(1000), opt.Page(pageCount))

		maxPages = res.NbPages
		if err != nil {
			return []DictionaryEntry{}, []DictionaryEntry{}, 0, 0, err
		}

		data, err := json.Marshal(res.Hits)
		if err != nil {
			return []DictionaryEntry{}, []DictionaryEntry{}, 0, 0, fmt.Errorf("cannot unmarshal dictionary entries: error while marshalling original dictionary entries: %v", err)
		}

		var entries []DictionaryEntry
		err = json.Unmarshal(data, &entries)
		if err != nil {
			return []DictionaryEntry{}, []DictionaryEntry{}, 0, 0, fmt.Errorf("cannot unmarshal dictionary entries: error while unmarshalling original dictionary entries: %v", err)
		}

		for _, entry := range entries {
			if entry.Type == CustomEntryType {
				customEntriesNb++
				customEntries = append(customEntries, entry)
			}

			allEntries = append(allEntries, entry)
			allEntriesNb++
		}
		pageCount++
	}

	return allEntries, customEntries, allEntriesNb, customEntriesNb, nil
}
