package export

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"
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
		Use:  "export <index_1>",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			client, err := opts.SearchClient()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			indexNames, err := cmdutil.IndexNames(client)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return indexNames, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Export the indice records",
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

			client, err := opts.SearchClient()
			if err != nil {
				return err
			}

			// Test that the provided indices exists
			exists, err := client.InitIndex(opts.Indice).Exists()
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("the index %s does not exist", opts.Indice)
			}

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
