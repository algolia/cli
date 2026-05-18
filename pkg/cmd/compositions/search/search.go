package search

import (
	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// SearchOptions holds the dependencies and flags for the search command.
type SearchOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	CompositionID     string
	Query             string
	HitsPerPage       *int32
	Page              *int32
	Filters           string
	PrintFlags        *cmdutil.PrintFlags
}

// NewSearchCmd returns the `compositions search` command.
func NewSearchCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	var hitsPerPage int32
	var page int32

	cmd := &cobra.Command{
		Use:   "search <composition-id> <query>",
		Short: "Search a composition",
		Args:  validators.ExactArgsWithMsg(2, "compositions search requires a <composition-id> and a <query> argument."),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# Search a composition with a query
			$ algolia compositions search my-comp "running shoes"

			# Search with filters
			$ algolia compositions search my-comp "shirt" --filters "brand:Nike"

			# Search with pagination
			$ algolia compositions search my-comp "shirt" --hits-per-page 20 --page 2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
			opts.Query = args[1]
			if cmd.Flags().Changed("hits-per-page") {
				opts.HitsPerPage = &hitsPerPage
			}
			if cmd.Flags().Changed("page") {
				opts.Page = &page
			}
			return runSearchCmd(opts)
		},
	}

	cmd.Flags().Int32Var(&hitsPerPage, "hits-per-page", 20, "Number of hits per page")
	cmd.Flags().Int32Var(&page, "page", 0, "Page number")
	cmd.Flags().StringVar(&opts.Filters, "filters", "", "Filter expression")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runSearchCmd(opts *SearchOptions) error {
	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	params := algoliaComposition.NewParams(
		algoliaComposition.WithParamsQuery(opts.Query),
	)
	if opts.HitsPerPage != nil {
		params.HitsPerPage = opts.HitsPerPage
	}
	if opts.Page != nil {
		params.Page = opts.Page
	}
	if opts.Filters != "" {
		params.Filters = &opts.Filters
	}

	reqBody := algoliaComposition.NewRequestBody(
		algoliaComposition.WithRequestBodyParams(*params),
	)

	opts.IO.StartProgressIndicatorWithLabel("Searching")

	res, err := client.Search(client.NewApiSearchRequest(opts.CompositionID, reqBody))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
