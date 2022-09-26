package search

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	algoliaSearch "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// SearchOptions represents the options for the search command
type SearchOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*algoliaSearch.Client, error)

	Indice string

	SearchParams map[string]interface{}

	PrintFlags *cmdutil.PrintFlags
}

// NewSearchCmd returns a new instance of the search command
func NewSearchCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "search  <index>",
		Short:             "Search the given index",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Long:              `Search for objects in your index.`,
		Example: heredoc.Doc(`
			# Search for objects in the "BOOKS" index matching the query "tolkien"
			$ algolia search BOOKS --query "tolkien"

			# Search for objects in the "BOOKS" index matching the query "tolkien" with filters
			$ algolia search BOOKS --query "tolkien" --filters "'(category:Book OR category:Ebook) AND _tags:published'"

			# Search for objects in the "BOOKS" index matching the query "tolkien" while setting the number of hits per page and specifying the page to retrieve
			$ algolia search BOOKS --query "tolkien" --hitsPerPage 2 --page 4

			# Search for objects in the "BOOKS" index matching the query "tolkien" and export the results to a new line delimited JSON (ndjson) file
			$ algolia search BOOKS --query "tolkien" > books.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]
			searchParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.SearchParamsObject...)
			if err != nil {
				return err
			}
			opts.SearchParams = searchParams

			return runSearchCmd(opts)
		},
	}

	cmdutil.AddSearchParamsObjectFlags(cmd)

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runSearchCmd(opts *SearchOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Searching")

	// We use the `opt.ExtraOptions` to pass the `SearchParams` to the API.
	query, ok := opts.SearchParams["query"].(string)
	if !ok {
		query = ""
	} else {
		delete(opts.SearchParams, "query")
	}
	res, err := indice.Search(query, opt.ExtraOptions(opts.SearchParams))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
