package importSynonyms

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
)

type ImportOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice  string
	Scanner *bufio.Scanner
}

// NewImportCmd creates and returns an import command for indice synonyms
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var file string

	cmd := &cobra.Command{
		Use:  "import <index_1> -F <file_1>",
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
		Short: "Import synonyms for the indice",
		Long: heredoc.Doc(`
			Import the synonyms for the provided indice.
		`),
		Example: heredoc.Doc(`
			$ algolia import TEST_PRODUCTS_1 -F synonyms.json
			$ cat rules.json | algolia synonyms import TEST_PRODUCTS_1 -F -
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

	cmd.Flags().StringVarP(&file, "file", "F", "", "Read synonyms to import from `file` (use \"-\" to read from standard input)")

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

	opts.IO.StartProgressIndicatorWithLabel("Importing synonyms")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		var rule search.Rule
		if err := json.Unmarshal([]byte(line), &rule); err != nil {
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
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d synonyms", totalCount))
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
		fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s synonyms to %s\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(totalCount)), opts.Indice)
	}

	return nil
}
