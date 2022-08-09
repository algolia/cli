package importSynonyms

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

	Index   string
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
		Use:               "import <index> -F <file>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Import synonyms to the indice",
		Long: heredoc.Doc(`
			Import synonyms to the provided indice.
			The file must contains one single JSON synonym per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import synonyms from the "synonyms.ndjson" file to the "TEST_PRODUCTS_1" index
			$ algolia import TEST_PRODUCTS_1 -F synonyms.ndjson

			# Import objects from the standard input to the "TEST_PRODUCTS_1" index
			$ cat synonyms.ndjson | algolia synonyms import TEST_PRODUCTS_1 -F -

			# Browse the synonyms in the "TEST_PRODUCTS_1" index and import them to the "TEST_PRODUCTS_2" index
			$ algolia synonyms browse TEST_PRODUCTS_1 | algolia synonyms import TEST_PRODUCTS_2 -F -
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

	cmd.Flags().StringVarP(&file, "file", "F", "", "Read synonyms to import from `file` (use \"-\" to read from standard input)")

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
		batch      = make([]search.Synonym, 0, batchSize)
		count      = 0
		totalCount = 0
	)

	opts.IO.StartProgressIndicatorWithLabel("Importing synonyms")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		lineB := []byte(line)
		var rawSynonym map[string]interface{}

		// Unmarshal as map[string]interface{} to get the type of the synonym
		if err := json.Unmarshal(lineB, &rawSynonym); err != nil {
			err := fmt.Errorf("failed to parse JSON synonym on line %d: %s", count, err)
			return err
		}
		typeString := rawSynonym["type"].(string)

		// This is really ugly, but algoliasearch package doesn't provide a way to
		// unmarshal a synonym from a JSON string.
		switch search.SynonymType(typeString) {
		case search.RegularSynonymType:
			var syn search.RegularSynonym
			err = json.Unmarshal(lineB, &syn)
			if err != nil {
				return err
			}
			batch = append(batch, syn)

		case search.OneWaySynonymType:
			var syn search.OneWaySynonym
			err = json.Unmarshal(lineB, &syn)
			if err != nil {
				return err
			}
			batch = append(batch, syn)

		case search.AltCorrection1Type:
			var syn search.AltCorrection1
			err = json.Unmarshal(lineB, &syn)
			if err != nil {
				return err
			}
			batch = append(batch, syn)

		case search.AltCorrection2Type:
			var syn search.AltCorrection2
			err = json.Unmarshal(lineB, &syn)
			if err != nil {
				return err
			}
			batch = append(batch, syn)

		case search.PlaceholderType:
			var syn search.Placeholder
			err = json.Unmarshal(lineB, &syn)
			if err != nil {
				return err
			}
			batch = append(batch, syn)

		default:
			return fmt.Errorf("cannot unmarshal synonym: unknown type %s", typeString)
		}

		count++

		if count == batchSize {
			if _, err := indice.SaveSynonyms(batch); err != nil {
				return err
			}

			batch = make([]search.Synonym, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d synonyms", totalCount))
			count = 0
		}
	}

	if count > 0 {
		totalCount += count
		if _, err := indice.SaveSynonyms(batch); err != nil {
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s synonyms to %s\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(totalCount)), opts.Index)
	}

	return nil
}
