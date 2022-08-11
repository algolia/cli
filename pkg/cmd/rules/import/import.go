package importRules

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice  string
	Scanner *bufio.Scanner
}

// NewImportCmd creates and returns an import command for indice rules
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
		Short:             "Import rules to the specified index",
		Long: heredoc.Doc(`
			Import rules to the specified index.
			The file must contains one JSON rule per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import rules from the "rules.ndjson" file to the "TEST_PRODUCTS_1" index
			$ algolia import TEST_PRODUCTS_1 -F rules.ndjson

			# Import rules from the standard input to the "TEST_PRODUCTS_1" index
			$ cat rules.ndjson | algolia rules import TEST_PRODUCTS_1 -F -

			# Browse the rules in the "TEST_PRODUCTS_1" index and import them to the "TEST_PRODUCTS_2" index
			$ algolia rules browse TEST_PRODUCTS_2 | algolia rules import TEST_PRODUCTS_2 -F -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			scanner, err := cmdutil.ScanFile(file, opts.IO.In)
			if err != nil {
				return err
			}
			opts.Scanner = scanner

			return runImportCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "F", "", "Read rules to import from `file` (use \"-\" to read from standard input)")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)

	// Move the following code to another module?
	var (
		batchSize  = 1000
		batch      = make([]search.Rule, 0, batchSize)
		count      = 0
		totalCount = 0
	)

	opts.IO.StartProgressIndicatorWithLabel("Importing rules")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		var rule search.Rule
		if err := json.Unmarshal([]byte(line), &rule); err != nil {
			err := fmt.Errorf("failed to parse JSON rule on line %d: %s", count, err)
			return err
		}

		batch = append(batch, rule)
		count++

		if count == batchSize {
			if _, err := indice.SaveRules(batch); err != nil {
				return err
			}

			batch = make([]search.Rule, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d rules", totalCount))
			count = 0
		}
	}

	if count > 0 {
		totalCount += count
		if _, err := indice.SaveRules(batch); err != nil {
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s rules to %s\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(totalCount)), opts.Indice)
	}

	return nil
}
