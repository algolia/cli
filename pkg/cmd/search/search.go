package search

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	algoliaSearch "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
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
		Use:               "search <index>",
		Short:             "Search the given index",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Long:              `Search for objects in your index.`,
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "search",
		},
		Example: heredoc.Doc(`
			# Search for objects in the "MOVIES" index matching the query "toy story"
			$ algolia search MOVIES --query "toy story"

			# Search for objects in the "MOVIES" index matching the query "toy story" with filters
			$ algolia search MOVIES --query "toy story" --filters "'(genres:Animation OR genres:Family) AND original_language:en'"

			# Search for objects in the "MOVIES" index matching the query "toy story" while setting the number of hits per page and specifying the page to retrieve
			$ algolia search MOVIES --query "toy story" --hitsPerPage 2 --page 4

			# Search for objects in the "MOVIES" index matching the query "toy story" and export the response to a .json file
			$ algolia search MOVIES --query "toy story" > movies.json

			# Search for objects in the "MOVIES" index matching the query "toy story" and only export the results to a .json file
			$ algolia search MOVIES --query "toy story" --output="jsonpath={$.Hits}" > movies.json
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

	cmd.SetUsageFunc(
		cmdutil.UsageFuncWithFilteredAndInheritedFlags(f.IOStreams, cmd, []string{"query"}),
	)

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
