package browse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index string

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for synonyms
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &BrowseOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "browse <index>",
		Aliases:           []string{"list"},
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		Short:             "List all the the synonyms of the given index",
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "settings",
		},
		Example: heredoc.Doc(`
			# List all the synonyms of the 'MOVIES' index
			$ algolia synonyms browse MOVIES

			# List all the synonyms of the 'MOVIES' and save them to the 'synonyms.json' file
			$ algolia synonyms browse MOVIES > synonyms.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			return runBrowseCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runBrowseCmd(opts *BrowseOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}
	// Check if index exists, because the API just returns an empty list if it doesn't
	exists, err := client.IndexExists(opts.Index)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Index %s doesn't exist.", opts.Index)
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	err = client.BrowseSynonyms(
		opts.Index,
		*search.NewEmptySearchSynonymsParams(),
		search.WithAggregator(func(res any, _ error) {
			for _, synonym := range res.(*search.SearchSynonymsResponse).Hits {
				p.Print(opts.IO, synonym)
			}
		}),
	)
	if err != nil {
		return err
	}
	return nil
}
