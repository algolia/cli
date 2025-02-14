package browse

import (
	"encoding/json"

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

	Index        string
	BrowseParams search.BrowseParamsObject

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for index objects
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &BrowseOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "browse <index>",
		Aliases:           []string{"list", "l"},
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "browse",
		},
		Short: "Browse the index objects",
		Long: heredoc.Doc(`
			This command browse the objects of the specified index.
		`),
		Example: heredoc.Doc(`
			# Browse the objects from the "MOVIES" index
			$ algolia objects browse MOVIES

			# Browse the objects from the "MOVIES" index and select which attributes to retrieve
			$ algolia objects browse MOVIES --attributesToRetrieve title,overview

			# Browse the objects from the "MOVIES" index with filters
			$ algolia objects browse MOVIES --filters "genres:Drama"

			# Browse the objects from the "MOVIES" and export the results to a new line delimited JSON (ndjson) file
			$ algolia objects browse MOVIES > movies.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			browseParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.BrowseParamsObject...)
			if err != nil {
				return err
			}

			// Convert map to object
			tmp, err := json.Marshal(browseParams)
			if err != nil {
				return err
			}
			err = json.Unmarshal(tmp, &opts.BrowseParams)
			if err != nil {
				return err
			}

			return runBrowseCmd(opts)
		},
	}

	cmd.SetUsageFunc(cmdutil.UsageFuncWithInheritedFlagsOnly(f.IOStreams, cmd))

	cmdutil.AddSearchParamsObjectFlags(cmd)
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runBrowseCmd(opts *BrowseOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	err = client.BrowseObjects(
		opts.Index,
		opts.BrowseParams,
		search.WithAggregator(func(res any, err error) {
			if err != nil {
				return
			}
			for _, hit := range res.(*search.BrowseResponse).Hits {
				err := p.Print(opts.IO, hit)
				if err != nil {
					return
				}
			}
		}),
	)
	if err != nil {
		return err
	}
	return nil
}
