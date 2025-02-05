package browse

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index string

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for Rules
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "browse <index>",
		Args:              validators.ExactArgs(1),
		Aliases:           []string{"list", "l"},
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "List all the rules of an index",
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "settings",
		},
		Example: heredoc.Doc(`
			# List all the rules of the "MOVIES" index
			$ algolia rules browse MOVIES

			# List all the rules of the "MOVIES" index and save them to a 'rules.ndjson' file
			$ algolia rules browse MOVIES -o json > rules.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ExportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	// Check if index exists because the API just returns an empty list if it doesn't
	exists, err := client.IndexExists(opts.Index)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("index %s doesn't exist", opts.Index)
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	err = client.BrowseRules(
		opts.Index,
		*search.NewEmptySearchRulesParams(),
		search.WithAggregator(func(res any, _ error) {
			for _, rule := range res.(*search.SearchRulesResponse).Hits {
				if err = p.Print(opts.IO, rule); err != nil {
					continue
				}
			}
		}),
	)
	if err != nil {
		return err
	}
	return nil
}
