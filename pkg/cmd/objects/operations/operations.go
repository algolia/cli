package operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/text"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type OperationsOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Wait bool

	File    string
	Scanner *bufio.Scanner

	ContinueOnError bool
}

// NewOperationsCmd creates and returns an operations command for object operations
func NewOperationsCmd(f *cmdutil.Factory, runF func(*OperationsOptions) error) *cobra.Command {
	opts := &OperationsOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:     "operations -F <file> [--wait] [--continue-on-errors]",
		Args:    validators.NoArgs(),
		Aliases: []string{"operation", "batch"},
		Annotations: map[string]string{
			"acls": "addObject,deleteObject,deleteIndex",
		},
		Short: "Perform several indexing operations",
		Long: heredoc.Doc(`
			Perform several indexing operations.

			The file must contains one single JSON object per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
			Each JSON object must be a valid indexing operation, as documented in the REST API documentation: https://www.algolia.com/doc/rest-api/search/#batch-write-operations-multiple-indices
		`),
		Example: heredoc.Doc(`
			# Batch operations from the "operations.ndjson" file
			$ algolia objects operations -F operations.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner, err := cmdutil.ScanFile(opts.File, opts.IO.In)
			if err != nil {
				return err
			}
			opts.Scanner = scanner

			if runF != nil {
				return runF(opts)
			}

			return runOperationsCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "The file to read the indexing operations from (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().
		BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the indexing operation(s) to complete before returning.")
	cmd.Flags().
		BoolVarP(&opts.ContinueOnError, "continue-on-error", "C", false, "Continue processing operations even if some operations are invalid.")

	return cmd
}

func runOperationsCmd(opts *OperationsOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	var (
		batchRequests   []search.MultipleBatchRequest
		currentLine     = 0
		totalOperations = 0
	)

	// Scan the file
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Reading operations from %s", opts.File))
	elapsed := time.Now()

	var errors []string
	for opts.Scanner.Scan() {
		currentLine++
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		totalOperations++
		opts.IO.UpdateProgressIndicatorLabel(
			fmt.Sprintf(
				"Read %s from %s",
				utils.Pluralize(totalOperations, "operation"),
				opts.File,
			),
		)

		var batchRequest search.MultipleBatchRequest
		if err := json.Unmarshal([]byte(line), &batchRequest); err != nil {
			err := fmt.Errorf("line %d: %s", currentLine, err)
			errors = append(errors, err.Error())
			continue
		}
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		batchRequests = append(batchRequests, batchRequest)
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	errorMsg := heredoc.Docf(`
		%s Found %s (out of %d operations) while parsing the file:
		%s
	`, cs.FailureIcon(), utils.Pluralize(len(errors), "error"), totalOperations, text.Indent(strings.Join(errors, "\n"), "  "))

	// No operations found
	if len(batchRequests) == 0 {
		if len(errors) > 0 {
			return fmt.Errorf(errorMsg)
		}
		return fmt.Errorf("%s No operations found in the file", cs.FailureIcon())
	}

	// Ask for confirmation if there are errors
	if len(errors) > 0 {
		if !opts.ContinueOnError {
			fmt.Print(errorMsg)

			var confirmed bool
			err = prompt.Confirm("Do you want to continue?", &confirmed)
			if err != nil {
				return fmt.Errorf("failed to prompt: %w", err)
			}
			if !confirmed {
				return nil
			}
		}
	}

	// Process operations
	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf("Processing %s operations", cs.Bold(fmt.Sprint(len(batchRequests)))),
	)
	res, err := client.MultipleBatch(
		client.NewApiMultipleBatchRequest(search.NewBatchParams(batchRequests)),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	// Wait for the operation to complete if requested
	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the operations to complete")
		for index, taskID := range res.TaskID {
			_, err := client.WaitForTask(index, taskID)
			if err != nil {
				opts.IO.StopProgressIndicator()
				return err
			}
		}
	}

	opts.IO.StopProgressIndicator()
	_, err = fmt.Fprintf(
		opts.IO.Out,
		"%s Successfully processed %s operations in %v\n",
		cs.SuccessIcon(),
		cs.Bold(fmt.Sprint(len(batchRequests))),
		time.Since(elapsed),
	)
	return err
}
