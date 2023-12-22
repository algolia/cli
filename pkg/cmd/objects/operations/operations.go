package operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
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

	SearchClient func() (*search.Client, error)

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
			Perform several indexing operations

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

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "The file to read the indexing operations from (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the indexing operation(s) to complete before returning.")
	cmd.Flags().BoolVarP(&opts.ContinueOnError, "continue-on-error", "C", false, "Continue processing operations even if some operations are invalid.")

	return cmd
}

func runOperationsCmd(opts *OperationsOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	var (
		operations      []search.BatchOperationIndexed
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
		opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Read %s from %s", utils.Pluralize(totalOperations, "operation"), opts.File))

		var batchOperation search.BatchOperationIndexed
		if err := json.Unmarshal([]byte(line), &batchOperation); err != nil {
			err := fmt.Errorf("line %d: %s", currentLine, err)
			errors = append(errors, err.Error())
			continue
		}
		err = ValidateBatchOperation(batchOperation)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		operations = append(operations, batchOperation)
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
	if len(operations) == 0 {
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
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Processing %s operations", cs.Bold(fmt.Sprint(len(operations)))))
	res, err := client.MultipleBatch(operations)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	// Wait for the operation to complete if requested
	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the operations to complete")
		if err := res.Wait(); err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()
	_, err = fmt.Fprintf(opts.IO.Out, "%s Successfully processed %s operations in %v\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(len(operations))), time.Since(elapsed))
	return err
}

// ValidateBatchOperation checks that the batch operation is valid
func ValidateBatchOperation(p search.BatchOperationIndexed) error {
	allowedActions := []string{
		string(search.AddObject), string(search.UpdateObject), string(search.PartialUpdateObject),
		string(search.PartialUpdateObjectNoCreate), string(search.DeleteObject),
	}
	extra := fmt.Sprintf("valid actions are %s", utils.SliceToReadableString(allowedActions))

	if p.Action == "" {
		return fmt.Errorf("missing action")
	}
	if !utils.Contains(allowedActions, string(p.Action)) {
		return fmt.Errorf("invalid action \"%s\" (%s)", p.Action, extra)
	}
	if p.IndexName == "" {
		return fmt.Errorf("missing index name for action \"%s\"", p.Action)
	}
	if p.Action == search.DeleteObject {
		switch body := p.Body.(type) {
		case map[string]interface{}:
			if body["objectID"] == nil || body["objectID"] == "" {
				return fmt.Errorf("missing objectID for action %s", search.DeleteObject)
			}
		default:
			return fmt.Errorf("missing objectID for action %s", search.DeleteObject)
		}
	}

	return nil
}
