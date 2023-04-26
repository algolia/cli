package update

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
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

type UpdateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index             string
	CreateIfNotExists bool
	Wait              bool

	File    string
	Scanner *bufio.Scanner

	ContinueOnError bool
}

// NewUpdateCmd creates and returns an update command for index objects
func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "update <index> -F <file> [--create-if-not-exists] [--wait] [--continue-on-errors]",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Update objects from a file to the specified index",
		Long: heredoc.Doc(`
			Update objects from a file to the specified index.
			
			The file must contains one single JSON object per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Update objects from the "objects.ndjson" file to the "TEST_PRODUCTS" index
			$ algolia objects update TEST_PRODUCTS -F objects.ndjson

			# Update objects from the "objects.ndjson" file to the "TEST_PRODUCTS" index and create the objects if they don't exist
			$ algolia objects update TEST_PRODUCTS -F objects.ndjson --create-if-not-exists

			# Update objects from the "objects.ndjson" file to the "TEST_PRODUCTS" index and wait for the operation to complete
			$ algolia objects update TEST_PRODUCTS -F objects.ndjson --wait

			# Update objects from the "objects.ndjson" file to the "TEST_PRODUCTS" index and continue updating objects even if some objects are invalid
			$ algolia objects update TEST_PRODUCTS -F objects.ndjson --continue-on-errors
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			scanner, err := cmdutil.ScanFile(opts.File, opts.IO.In)
			if err != nil {
				return err
			}
			opts.Scanner = scanner

			if runF != nil {
				return runF(opts)
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "Read objects to update from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().BoolVarP(&opts.CreateIfNotExists, "create-if-not-exists", "c", false, "If provided, updating a nonexistent object will create a new object with the objectID and the attributes defined in the object")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the operation to complete before returning")

	cmd.Flags().BoolVarP(&opts.ContinueOnError, "continue-on-error", "C", false, "Continue updating objects even if some objects are invalid.")

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	index := client.InitIndex(opts.Index)

	var (
		objects      []interface{}
		currentLine  = 0
		totalObjects = 0
	)

	// Scan the file
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Reading objects from %s", opts.File))
	elapsed := time.Now()

	var errors []string
	for opts.Scanner.Scan() {
		currentLine++
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		totalObjects++
		opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Read %s from %s", utils.Pluralize(totalObjects, "object"), opts.File))

		var obj Object
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			err := fmt.Errorf("line %d: %s", currentLine, err)
			errors = append(errors, err.Error())
			continue
		}

		objects = append(objects, obj)
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {
		return err
	}

	errorMsg := heredoc.Docf(`
		%s Found %s (out of %d objects) while parsing the file:
		%s
	`, cs.FailureIcon(), utils.Pluralize(len(errors), "error"), totalObjects, text.Indent(strings.Join(errors, "\n"), "  "))

	// No objects found
	if len(objects) == 0 {
		if len(errors) > 0 {
			return fmt.Errorf(errorMsg)
		}
		return fmt.Errorf("%s No objects found in the file", cs.FailureIcon())
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

	// Update the objects
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Updating %s objects on %s", cs.Bold(fmt.Sprint(len(objects))), cs.Bold(opts.Index)))
	res, err := index.PartialUpdateObjects(objects, opt.CreateIfNotExists(opts.CreateIfNotExists))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	// Wait for the operation to complete if requested
	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for operation to complete")
		if err := res.Wait(); err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()
	_, err = fmt.Fprintf(opts.IO.Out, "%s Successfully updated %s objects on %s in %v\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(len(objects))), cs.Bold(opts.Index), time.Since(elapsed))
	return err
}
