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
		Short:             "Browse the index objects",
		Long: heredoc.Doc(`
			This command browse the objects of the specified index.
		`),
		Example: heredoc.Doc(`
			# Browse the objects from the "BOOKS" index
			$ algolia objects browse BOOKS

			# Browse the objects from the "BOOKS" index and select which attributes to retrieve
			$ algolia objects browse BOOKS --attributesToRetrieve author,title,description

			# Browse the objects from the "BOOKS" index with filters
			$ algolia objects browse BOOKS --filters "'(category:Book OR category:Ebook) AND _tags:published'"

			# Browse the objects from the "BOOKS" and export the results to a new line delimited JSON (ndjson) file
			$ algolia objects browse BOOKS > books.ndjson
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
		p.Print(opts.IO, iObject)
	}
}
