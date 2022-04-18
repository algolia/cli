package search

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
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

	Query    string
	Settings *algoliaSearch.Settings

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
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Long:              `Search for objects in your index.`,
		Example: heredoc.Doc(`
			# Search for objects in your index
			algolia search PRODUCTS --query "foo"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			if cmd.Flags().Changed("settings") {
				if searchSettings, err := cmd.Flags().GetString("settings"); err == nil {
					var settings algoliaSearch.Settings
					if err := json.Unmarshal([]byte(searchSettings), &settings); err != nil {
						return fmt.Errorf("invalid settings: %v", err)
					}
					opts.Settings = &settings
				}
			}

			return runSearchCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "Query")
	cmd.Flags().StringP("settings", "s", "", "Settings")

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
	res, err := indice.Search(opts.Query, opts.Settings)
	if err != nil {
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
