package importRecords

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type ImportOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index string

	AutoGenerateObjectIDIfNotExist bool

	Scanner *bufio.Scanner
}

// NewImportCmd creates and returns an import command for indice object
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var file string

	cmd := &cobra.Command{
		Use:               "import <index> -F <file>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Import objects to the specified indice",
		Long: heredoc.Doc(`
			Import objects to the specified indice from a file / the standard input.
			The file must contains one JSON object per line (newline delimited JSON objects - ndjson format).
		`),
		Example: heredoc.Doc(`
			# Import objects from the "data.ndjson" file to the "TEST_PRODUCTS_1" index
			$ algolia objects import TEST_PRODUCTS_1 -F data.ndjson

			# Import objects from the standard input to the "TEST_PRODUCTS_1" index
			$ cat data.ndjson | algolia objects import TEST_PRODUCTS_1 -F -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			scanner, err := cmdutil.ScanFile(file, opts.IO.In)
			if err != nil {
				return err
			}
			opts.Scanner = scanner

			return runImportCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "F", "", "Read records to import from `file` (use \"-\" to read from standard input)")

	cmd.Flags().BoolVar(&opts.AutoGenerateObjectIDIfNotExist, "auto-generate-object-id-if-not-exist", false, "Automatically generate object ID if not exist")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Index)

	// Move the following code to another module?
	var (
		batchSize  = 1000
		batch      = make([]interface{}, 0, batchSize)
		count      = 0
		totalCount = 0
	)

	options := []interface{}{opt.AutoGenerateObjectIDIfNotExist(opts.AutoGenerateObjectIDIfNotExist)}

	opts.IO.StartProgressIndicatorWithLabel("Importing records")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		var obj interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			return err
		}

		batch = append(batch, obj)
		count++

		if count == batchSize {
			if _, err := indice.SaveObjects(batch, options...); err != nil {
				return err
			}

			batch = make([]interface{}, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d objects", totalCount))
			count = 0
		}
	}

	if count > 0 {
		totalCount += count
		if _, err := indice.SaveObjects(batch, options...); err != nil {
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s objects to %s\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(totalCount)), opts.Index)
	}

	return nil
}
