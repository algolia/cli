package browse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Dictionaries            []search.DictionaryType
	All                     bool
	IncludeDefaultStopwords bool

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for dictionaries' entries.
func NewBrowseCmd(f *cmdutil.Factory, runF func(*BrowseOptions) error) *cobra.Command {
	cs := f.IOStreams.ColorScheme()

	opts := &BrowseOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4_SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:       "browse {<dictionary>... | --all} [--include-defaults]",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: shared.DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return shared.DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Annotations: map[string]string{
			"acls": "settings",
		},
		Short: "Browse dictionary entries",
		Long: heredoc.Docf(`
			This command retrieves all entries from the specified %s dictionaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Retrieve all entries from the "stopwords" dictionary (doesn't include default stopwords)
			$ algolia dictionary entries browse stopwords

			# Retrieve all entries from the "stopwords" and "plurals" dictionaries
			$ algolia dictionary entries browse stopwords plurals

			# Retrieve all entries from all dictionaries
			$ algolia dictionary entries browse --all

			# Retrieve all entries from the "stopwords" dictionaries (including default stopwords)
			$ algolia dictionary entries browse stopwords --include-defaults
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.All && len(args) > 0 || !opts.All && len(args) == 0 {
				return cmdutil.FlagErrorf(
					"Either specify dictionaries' names or use --all to browse all dictionaries",
				)
			}

			if opts.All {
				opts.Dictionaries = []search.DictionaryType{
					search.DICTIONARY_TYPE_STOPWORDS,
					search.DICTIONARY_TYPE_PLURALS,
					search.DICTIONARY_TYPE_COMPOUNDS,
				}
			} else {
				opts.Dictionaries = make([]search.DictionaryType, len(args))
				for i, dictionary := range args {
					opts.Dictionaries[i] = search.DictionaryType(dictionary)
				}
			}

			return runBrowseCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "browse all dictionaries")
	cmd.Flags().
		BoolVarP(&opts.IncludeDefaultStopwords, "include-defaults", "d", false, "include default stopwords")

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

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	for _, dictionary := range opts.Dictionaries {
		var pageCount int32 = 0
		var maxPages int32 = 1

		// implement infinite pagination
		for pageCount < maxPages {
			res, err := client.SearchDictionaryEntries(
				client.NewApiSearchDictionaryEntriesRequest(
					dictionary,
					search.NewEmptySearchDictionaryEntriesParams().
						SetHitsPerPage(1000).
						SetPage(pageCount),
				),
			)
			if err != nil {
				return err
			}

			maxPages = res.NbPages

			if res.NbHits == 0 {
				if _, err = fmt.Fprintf(opts.IO.Out, "%s No entries found.\n\n", cs.WarningIcon()); err != nil {
					return err
				}
				// go to the next dictionary
				break
			}

			for _, entry := range res.Hits {
				if opts.IncludeDefaultStopwords {
					// print all entries (default stopwords included)
					if err = p.Print(opts.IO, entry); err != nil {
						return err
					}
					// TODO: This will break when `type` is added as a type to the response
				} else if entry.AdditionalProperties["type"] == "custom" {
					// print only custom entries
					if err = p.Print(opts.IO, entry); err != nil {
						return err
					}
				}
			}

			pageCount++
		}
	}

	return nil
}
