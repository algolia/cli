package browse

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type ExportOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index string
}

// NewBrowseCmd creates and returns a browse command for index objects
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "browse <index-1>",
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
			opts.Index = args[0]

			return runBrowseCmd(opts)
		},
	}

	return cmd
}

func runBrowseCmd(opts *ExportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Index)
	res, err := indice.BrowseObjects()
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
		obj, err := json.Marshal(iObject)
		if err != nil {
			return err
		}
		opts.IO.Out.Write([]byte(fmt.Sprintf("%s\n", obj)))
	}
}
