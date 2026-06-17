package search

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/pkg/iostreams"
)

// SearchOptions holds the dependencies and flags for the search command.
type SearchOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	Query             string
	HitsPerPage       *int32
	Page              *int32
	Filters           string
	Interactive       bool
	PrintFlags        *cmdutil.PrintFlags
}

// NewSearchCmd returns the `compositions search` command.
func NewSearchCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	var hitsPerPage int32
	var page int32

	cmd := &cobra.Command{
		Use:   "search <composition-id> [query]",
		Short: "Search a composition",
		Args:  cobra.RangeArgs(1, 2),
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

			# Build the search request interactively
			$ algolia compositions search my-comp --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]

			if opts.Interactive {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("`--interactive` requires a terminal")
				}
				return runSearchCmd(opts)
			}

			if len(args) < 2 {
				return cmdutil.FlagErrorf("a <query> argument is required (or use `--interactive`)")
			}
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
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Build the search request interactively")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

// buildRequestBody assembles the search request body, interactively when
// requested or from the query argument and flags otherwise.
func buildRequestBody(opts *SearchOptions) (*algoliaComposition.RequestBody, error) {
	if opts.Interactive {
		var body algoliaComposition.RequestBody
		prompter := opts.Prompter
		if prompter == nil {
			prompter = interactive.NewSurveyPrompter(opts.IO)
		}
		if err := (&interactive.Builder{Prompter: prompter}).Build(&body); err != nil {
			return nil, fmt.Errorf("building search request: %w", err)
		}
		return &body, nil
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
	return algoliaComposition.NewRequestBody(
		algoliaComposition.WithRequestBodyParams(*params),
	), nil
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

	reqBody, err := buildRequestBody(opts)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Searching")

	res, err := client.Search(client.NewApiSearchRequest(opts.CompositionID, reqBody))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
