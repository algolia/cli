package export

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

	Indice string
}

// NewExportCmd creates and returns an export command for indice records
func NewExportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "export <index_1>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Export the indice records",
		Long: heredoc.Doc(`
			Export the given indice records.
			This command export the records of the specified indice.
		`),
		Example: heredoc.Doc(`
			$ algolia indices export TEST_PRODUCTS_1
			$ algolia indices export TEST_PRODUCTS_1 > records.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runExportCmd(opts)
		},
	}

	return cmd
}

func runExportCmd(opts *ExportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
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
