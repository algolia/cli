package importentries

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/dictionary/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"

	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/text"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	DictionaryName    string
	CreateIfNotExists bool
	Wait              bool

	File    string
	Scanner *bufio.Scanner

	ContinueOnError bool
}

// NewImportCmd creates and returns an import command for dictionary
func NewImportCmd(f *cmdutil.Factory, runF func(*ImportOptions) error) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:       "import <dictionary> -F <file> [--wait] [--continue-on-errors]",
		Args:      validators.ExactArgs(1),
		ValidArgs: shared.DictionaryNames(),
		Short:     "Import dictionary entries from a file to the specified index",
		Long: heredoc.Doc(`
			Import dictionary entries from a file to the specified index.
			
			The file must contains one single JSON object per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import entries from the "entries.ndjson" file to "stopwords" dictionary
			$ algolia dictionary import stopwords -F entries.ndjson

			# Import entries from the "entries.ndjson" file to "plurals" dictionary and continue importing entries even if some entries are invalid
			$ algolia dictionary import plurals -F entries.ndjson --continue-on-errors
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.DictionaryName = args[0]

			scanner, err := cmdutil.ScanFile(opts.File, opts.IO.In)
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

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "Read entries to import from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the operation to complete before returning")
	cmd.Flags().BoolVarP(&opts.ContinueOnError, "continue-on-error", "C", false, "Continue importing entries even if some entries are invalid.")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	var (
		entries      []search.DictionaryEntry
		currentLine  = 0
		totalEntries = 0
	)

	// Scan the file
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Reading words from %s", opts.File))
	elapsed := time.Now()

	var errors []string
	for opts.Scanner.Scan() {
		currentLine++
		line := opts.Scanner.Text()
		if line == "" {
			continue
		}

		totalEntries++
		opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Read entries from %s", opts.File))

		var entry shared.DictionaryEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			err := fmt.Errorf("line %d: %s", currentLine, err)
			errors = append(errors, err.Error())
			continue
		}
		err := ValidateDictionaryEntry(entry, currentLine)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		dictionaryEntry, err := createDictionaryEntry(opts.DictionaryName, entry)
		if err != nil {
			errors = append(errors, fmt.Errorf("line %d: %s", currentLine, err.Error()).Error())
			continue
		}
		entries = append(entries, dictionaryEntry)
	}

	opts.IO.StopProgressIndicator()

	if err := opts.Scanner.Err(); err != nil {

		return err
	}

	errorMsg := heredoc.Docf(`
		%s Found %s (out of %d entries) while parsing the file:
		%s
	`, cs.FailureIcon(), utils.Pluralize(len(errors), "error"), totalEntries, text.Indent(strings.Join(errors, "\n"), "  "))

	// No entries found
	if len(entries) == 0 {
		if len(errors) > 0 {
			return fmt.Errorf(errorMsg)
		}
		return fmt.Errorf("%s No entries found in the file", cs.FailureIcon())
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

	// Import entries
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Updating %s entries on %s", cs.Bold(fmt.Sprint(len(entries))), cs.Bold(opts.DictionaryName)))

	res, err := client.SaveDictionaryEntries(search.DictionaryName(opts.DictionaryName), entries)
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
	fmt.Fprintf(opts.IO.Out, "%s Successfully imported %s entries on %s in %v\n", cs.SuccessIcon(), cs.Bold(fmt.Sprint(len(entries))), cs.Bold(opts.DictionaryName), time.Since(elapsed))
	return nil
}

func ValidateDictionaryEntry(entry shared.DictionaryEntry, currentLine int) error {
	if entry.ObjectID == "" {
		return fmt.Errorf("line %d: objectID is missing", currentLine)
	}
	if entry.Word == "" {
		return fmt.Errorf("line %d: word is missing", currentLine)
	}
	if entry.Language == "" {
		return fmt.Errorf("line %d: language is missing", currentLine)
	}

	return nil
}

func createDictionaryEntry(dictionaryName string, entry shared.DictionaryEntry) (search.DictionaryEntry, error) {
	switch dictionaryName {
	case string(search.Plurals):
		return search.NewPlural(entry.ObjectID, entry.Language, entry.Words), nil
	case string(search.Compounds):
		return search.NewCompound(entry.ObjectID, entry.Language, entry.Word, entry.Decomposition), nil
	case string(search.Stopwords):
		return search.NewStopword(entry.ObjectID, entry.Language, entry.Word, entry.State), nil
	}

	return nil, fmt.Errorf("Wrong dictionary name")
}
