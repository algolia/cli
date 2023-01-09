package browse

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Dictionnaries           []search.DictionaryName
	All                     bool
	IncludeDefaultStopwords bool

	PrintFlags *cmdutil.PrintFlags
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
		ValidArgs: shared.DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return shared.DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Browse dictionary entries",
		Long: heredoc.Docf(`
			This command retrieves all entries from the specified %s dictionnaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Retrieve all entries from the "stopwords" dictionary (doesn't include default stopwords)
			$ algolia dictionary entries browse stopwords

			# Retrieve all entries from the "stopwords" and "plurals" dictionnaries
			$ algolia dictionary entries browse stopwords plurals

			# Retrieve all entries from all dictionnaries
			$ algolia dictionary entries browse --all

			# Retrieve all entries from the "stopwords" dictionnaries (including default stopwords)
			$ algolia dictionary entries browse stopwords --include-defaults
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
	cmd.Flags().BoolVarP(&opts.IncludeDefaultStopwords, "include-defaults", "d", false, "browse dictionaries and include default stopwords")

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

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	for _, dictionnary := range dictionaries {
		pageCount := 0
		maxPages := 1

		// implement infinite pagination
		for pageCount < maxPages {
			res, err := client.SearchDictionaryEntries(dictionnary, "", opt.HitsPerPage(1000), opt.Page(pageCount))
			if err != nil {
				return err
			}

			maxPages = res.NbPages

			data, err := json.Marshal(res.Hits)
			if err != nil {
				return fmt.Errorf("cannot unmarshal dictionary entries: error while marshalling original dictionary entries: %v", err)
			}

			var entries []DictionaryEntry
			err = json.Unmarshal(data, &entries)
			if err != nil {
				return fmt.Errorf("cannot unmarshal dictionary entries: error while unmarshalling original dictionary entries: %v", err)
			}

			if len(entries) == 0 {
				fmt.Fprintf(opts.IO.Out, "%s No entries in %s dictionary.\n\n", cs.WarningIcon(), dictionnary)
				// go to the next dictionary
				break
			}

			for _, entry := range entries {
				if opts.IncludeDefaultStopwords {
					// print all entries (default stopwords included)
					p.Print(opts.IO, entry)
				} else if entry.Type == shared.CustomEntryType {
					// print only custom entries
					p.Print(opts.IO, entry)
				}
			}

			pageCount++
		}
	}

	return nil
}
