package importrecords

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ImportOptions struct {
	Config       config.IConfig
	IO           *iostreams.IOStreams
	SearchClient func() (*search.APIClient, error)
	Index        string

	Scanner   *bufio.Scanner
	BatchSize int
	Wait      bool
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
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "addObject",
		},
		Short: "Import objects to the specified index",
		Long: heredoc.Doc(`
			Import objects to the specified index from a file / the standard input.
			The file must contains one single JSON object per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import objects from the "data.ndjson" file to the "MOVIES" index
			$ algolia objects import MOVIES -F data.ndjson

			# Import objects from the standard input to the "MOVIES" index
			$ cat data.ndjson | algolia objects import MOVIES -F -

			# Browse the objects in the "SERIES" index and import them to the "MOVIES" index
			$ algolia objects browse SERIES | algolia objects import MOVIES -F -
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

	cmd.Flags().
		StringVarP(&file, "file", "F", "", "Read records to import from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().IntVarP(&opts.BatchSize, "batch-size", "b", 1000, "Specify the upload batch size")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")
	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	count := 0
	var records []map[string]any
	opts.IO.StartProgressIndicatorWithLabel("Importing records")
	elapsed := time.Now()
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			err := fmt.Errorf("failed to parse JSON object on line %d: %s", count, err)
			return err
		}
		records = append(records, record)
		count++
	}

	responses, err := client.SaveObjects(opts.Index, records, search.WithBatchSize(opts.BatchSize))
	if err != nil {
		return err
	}

	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
		for _, res := range responses {
			_, err := client.WaitForTask(opts.Index, res.TaskID)
			if err != nil {
				opts.IO.StopProgressIndicator()
				return err
			}
		}
	}

	opts.IO.UpdateProgressIndicatorLabel(
		fmt.Sprintf("Imported %d objects in %v", len(records), time.Since(elapsed)),
	)

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Successfully imported %s objects to %s in %v\n",
			cs.SuccessIcon(),
			cs.Bold(fmt.Sprint(len(records))),
			opts.Index,
			time.Since(elapsed),
		)
	}

	return nil
}
