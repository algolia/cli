package browse

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice       string
	BrowseParams map[string]interface{}

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
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "browse",
		},
		Short: "Browse records in an index.",
		Long: heredoc.Doc(`
			This command browses records in the specified index.
		`),
		Example: heredoc.Doc(`
			# Browse records in the "MOVIES" index
			$ algolia objects browse MOVIES

			# Browse records in the "MOVIES" index and select which attributes to retrieve
			$ algolia objects browse MOVIES --attributesToRetrieve title,overview

			# Browse records in the "MOVIES" index with filters
			$ algolia objects browse MOVIES --filters "genres:Drama"

			# Browse records in the "MOVIES" and export the results to a new line delimited JSON (ndjson) file
			$ algolia objects browse MOVIES > movies.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			browseParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.BrowseParamsObject...)
			if err != nil {
				return err
			}
			opts.BrowseParams = browseParams

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

	indice := client.InitIndex(opts.Indice)

	// We use the `opt.ExtraOptions` to pass the `SearchParams` to the API.
	query, ok := opts.BrowseParams["query"].(string)
	if !ok {
		query = ""
	} else {
		delete(opts.BrowseParams, "query")
	}
	res, err := indice.BrowseObjects(opt.Query(query), opt.ExtraOptions(opts.BrowseParams))
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	for {
		iObject, err := res.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err = p.Print(opts.IO, iObject); err != nil {
			return err
		}

	}
}
