package importRecords

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/sirupsen/logrus"
)

type ImportOptions struct {
	Config                         *config.Config
	IO                             *iostreams.IOStreams
	SearchClient                   func() (*search.Client, error)
	Index                          string
	AutoGenerateObjectIDIfNotExist bool
	Scanner                        *bufio.Scanner
	Verbose                        bool
	BatchSize                      int
}

// NewImportCmd creates and returns an import command for indice object
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var file string
	var verbose bool
	var batchSize int

	cmd := &cobra.Command{
		Use:               "import <index-1> -F <file-1>",
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
			opts.Verbose = verbose
			opts.BatchSize = batchSize

			return runImportCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "F", "", "Read records to import from `file` (use \"-\" to read from standard input)")
	cmd.Flags().BoolVar(&opts.AutoGenerateObjectIDIfNotExist, "auto-generate-object-id-if-not-exist", false, "Automatically generate object ID if not exist")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Activate verbose logging")
	cmd.Flags().IntVarP(&batchSize, "batch-size", "b", 1000, "Specify the upload batch size")
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
		batchSize  = opts.BatchSize
		batch      = make([]interface{}, 0, batchSize)
		count      = 0
		totalCount = 0
		log        = logrus.New()
	)

	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).DisableColors = true
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true

	options := []interface{}{opt.AutoGenerateObjectIDIfNotExist(opts.AutoGenerateObjectIDIfNotExist)}

	opts.IO.StartProgressIndicatorWithLabel("Importing records...")
	elapsed := time.Now()
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
			start := time.Now()
			if _, err := indice.SaveObjects(batch, options...); err != nil {
				return err
			}
			if opts.Verbose {
				log.Infof("Added [%d] records in %.3f seconds", count, time.Since(start).Seconds())
			}
			batch = make([]interface{}, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d objects so far...", totalCount))
			count = 0
		}
	}

	if count > 0 {
		totalCount += count
		start := time.Now()
		if _, err := indice.SaveObjects(batch, options...); err != nil {
			return err
		}
		if opts.Verbose {
			log.Infof("Added [%d] records in %.3f seconds", count, time.Since(start).Seconds())
		}
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s objects to %s in %.3f seconds\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(totalCount)), opts.Index, time.Since(elapsed).Seconds())
	}

	return nil
}
