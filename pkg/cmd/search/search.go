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
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*algoliaSearch.Client, error)

	Indice string
	Query  string

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
		Use:               "search  <index-name>",
		Short:             "Search the given index",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Long:              `Search for objects in your index.`,
		Example: heredoc.Doc(`
			# Search for objects in your index
			algolia search PRODUCTS --query "foo"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			query, err := cmd.Flags().GetString("query")
			if err != nil {
				return err
			}
			opts.Query = query

			searchParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.SearchParams...)
			if err != nil {
				return err
			}
			opts.SearchParams = searchParams

			return runSearchCmd(opts)
		},
	}

	cmdutil.AddSearchFlags(cmd)

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
	res, err := indice.Search(opts.Query, opt.ExtraOptions(opts.SearchParams))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
