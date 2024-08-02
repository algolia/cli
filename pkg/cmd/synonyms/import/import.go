package importSynonyms

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index                   string
	ForwardToReplicas       bool
	ReplaceExistingSynonyms bool
	Scanner                 *bufio.Scanner
	BatchSize               int
}

// NewImportCmd creates and returns an import command for indice synonyms
func NewImportCmd(f *cmdutil.Factory, runF func(*ImportOptions) error) *cobra.Command {
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
			"acls": "editSettings",
		},
		Short: "Import synonyms to the index",
		Long: heredoc.Doc(`
			Import synonyms to the provided index.
			The file must contains one single JSON synonym per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import synonyms from the "synonyms.ndjson" file to the "MOVIES" index
			$ algolia synonyms import MOVIES -F synonyms.ndjson

			# Import synonyms from the standard input to the "MOVIES" index
			$ cat synonyms.ndjson | algolia synonyms import MOVIES -F -

			# Browse the synonyms in the "SERIES" index and import them to the "MOVIES" index
			$ algolia synonyms browse SERIES | algolia synonyms import MOVIES -F -

			# Import synonyms from the "synonyms.ndjson" file to the "MOVIES" index and replace existing synonyms
			$ algolia synonyms import MOVIES -F synonyms.ndjson -r

			# Import synonyms from the "synonyms.ndjson" file to the "MOVIES" index and don't forward the synonyms to the index replicas
			$ algolia synonyms import MOVIES -F synonyms.ndjson -f=false
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			scanner, err := cmdutil.ScanFile(file, opts.IO.In)
			if err != nil {
				return err
			}
			opts.Scanner = scanner

			if runF != nil {
				return runF(opts)
			}

			return runImportCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&file, "file", "F", "", "Read synonyms to import from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().
		BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", true, "Forward the synonyms to the replicas of the index")
	cmd.Flags().
		BoolVarP(&opts.ReplaceExistingSynonyms, "replace-existing-synonyms", "r", false, "Replace existing synonyms in the index")
	cmd.Flags().IntVarP(&opts.BatchSize, "batch-size", "b", 1000, "Specify the upload batch size")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	// Move the following code to another module?
	var (
		batchSize  = 1000
		batch      = make([]search.SynonymHit, 0, batchSize)
		count      = 0
		totalCount = 0
	)

	opts.IO.StartProgressIndicatorWithLabel("Importing synonyms")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		var synonym search.SynonymHit
		err := json.Unmarshal([]byte(line), &synonym)
		if err != nil {
			return err
		}
		batch = append(batch, synonym)
		count++

		// If requested, only clear synonyms for the first batch (otherwise we'll keep deleting added synonyms)
		replaceSynonyms := false
		if count == 1 {
			replaceSynonyms = opts.ReplaceExistingSynonyms
		}

		if count == batchSize {
			if _, err := client.SaveSynonyms(
				client.NewApiSaveSynonymsRequest(opts.Index, batch).
					WithReplaceExistingSynonyms(replaceSynonyms).
					WithForwardToReplicas(opts.ForwardToReplicas),
				search.WithBatchSize(batchSize)); err != nil {
				return err
			}
			batch = make([]search.SynonymHit, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d synonyms", totalCount))
			count = 0
		}
	}

	if count > 0 {
		totalCount += count
		if _, err := client.SaveSynonyms(
			client.NewApiSaveSynonymsRequest(opts.Index, batch).
				WithForwardToReplicas(opts.ForwardToReplicas),
			search.WithBatchSize(batchSize)); err != nil {
			return err
		}
	}

	if totalCount == 0 && opts.ReplaceExistingSynonyms {
		if _, err := client.ClearSynonyms(
			client.
				NewApiClearSynonymsRequest(opts.Index).
				WithForwardToReplicas(opts.ForwardToReplicas)); err != nil {
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Successfully imported %s synonyms to %s\n",
			cs.SuccessIcon(),
			cs.Bold(fmt.Sprint(totalCount)),
			opts.Index,
		)
	}

	return nil
}
