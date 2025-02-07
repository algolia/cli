package importsynonyms

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
	Wait                    bool
	Scanner                 *bufio.Scanner
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
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	// Only clear existing synonyms on the first batch
	clearExistingSynonyms := opts.ReplaceExistingSynonyms

	// Move the following code to another module?
	var (
		batchSize  = 1000
		synonyms   = make([]search.SynonymHit, 0, batchSize)
		count      = 0
		totalCount = 0
		taskIDs    []int64
	)

	opts.IO.StartProgressIndicatorWithLabel("Importing synonyms")
	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		lineB := []byte(line)
		var synonym search.SynonymHit

		// Unmarshal as map[string]interface{} to get the type of the synonym
		if err := json.Unmarshal(lineB, &synonym); err != nil {
			err := fmt.Errorf("failed to parse JSON synonym on line %d: %s", count, err)
			return err
		}

		synonyms = append(synonyms, synonym)
		count++

		if count == batchSize {
			res, err := client.SaveSynonyms(
				client.NewApiSaveSynonymsRequest(opts.Index, synonyms).
					WithReplaceExistingSynonyms(clearExistingSynonyms).
					WithForwardToReplicas(opts.ForwardToReplicas),
			)
			if err != nil {
				return err
			}
			if opts.Wait {
				taskIDs = append(taskIDs, res.TaskID)
			}
			synonyms = make([]search.SynonymHit, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d synonyms", totalCount))
			count = 0
			clearExistingSynonyms = false
		}
	}

	if count > 0 {
		totalCount += count
		res, err := client.SaveSynonyms(
			client.NewApiSaveSynonymsRequest(opts.Index, synonyms).
				WithForwardToReplicas(opts.ForwardToReplicas),
		)
		if err != nil {
			return err
		}
		if opts.Wait {
			taskIDs = append(taskIDs, res.TaskID)
		}
	}

	if totalCount == 0 && opts.ReplaceExistingSynonyms {
		res, err := client.ClearSynonyms(
			client.NewApiClearSynonymsRequest(opts.Index).
				WithForwardToReplicas(opts.ForwardToReplicas),
		)
		if err != nil {
			return err
		}
		if opts.Wait {
			taskIDs = append(taskIDs, res.TaskID)
		}
	}

	if len(taskIDs) > 0 {
		for _, taskID := range taskIDs {
			_, err := client.WaitForTask(opts.Index, taskID)
			if err != nil {
				return err
			}
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
