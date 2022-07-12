package browse

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type BrowseOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice       string
	SearchParams map[string]interface{}

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
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Browse the index objects",
		Long: heredoc.Doc(`
			This command browse the objects of the specified index.
		`),
		Example: heredoc.Doc(`
			# Browse the objects from the "TEST_PRODUCTS_1" index
			$ algolia objects browse TEST_PRODUCTS_1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			searchParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.SearchParamsObject...)
			if err != nil {
				return err
			}
			opts.SearchParams = searchParams

			return runBrowseCmd(opts)
		},
	}

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
	res, err := indice.BrowseObjects(opts.SearchParams)
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
