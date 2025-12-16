package importentries

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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

	SearchClient func() (*search.APIClient, error)

	DictionaryType search.DictionaryType
	Wait           bool

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
		ValidArgs: shared.DictionaryTypes(),
		Annotations: map[string]string{
			"acls": "settings,editSettings",
		},
		Short: "Import dictionary entries from a file to the specified index",
		Long: heredoc.Doc(`
			Import dictionary entries from a file to the specified index.
			The file must contain one JSON object per line - in newline-delimited JSON (NDJSON) format: https://ndjson.org/.
		`),
		Example: heredoc.Doc(`
			# Import entries from the "entries.ndjson" file to the "stopwords" dictionary
			$ algolia dictionary import stopwords -F entries.ndjson

			# Import entries from the "entries.ndjson" file to the "plurals" dictionary and continue importing entries even if some entries are invalid
			$ algolia dictionary import plurals -F entries.ndjson --continue-on-errors
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := search.NewDictionaryTypeFromValue(args[0])
			if err != nil {
				return err
			}
			opts.DictionaryType = *d

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

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "Read entries to import from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().
		BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the operation to complete before returning")
	cmd.Flags().
		BoolVarP(&opts.ContinueOnError, "continue-on-error", "C", false, "Continue importing entries even if some entries are invalid.")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	var (
		entries      []*search.DictionaryEntry
		currentLine  = 0
		totalEntries = 0
	)

	// Scan the file
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Reading entries from %s", opts.File))
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

		var entry search.DictionaryEntry

		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			errors = append(errors, fmt.Errorf("line %d: %s", currentLine, err).Error())
			continue
		}

		fmt.Printf("TYPE: %v\n", opts.DictionaryType)
		fmt.Printf("ENTRY: %v\n", entry)

		dictionaryEntry, err := createDictionaryEntry(opts.DictionaryType, entry)
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
			return fmt.Errorf("%s", errorMsg)
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
	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf(
			"Updating %s entries on %s",
			cs.Bold(fmt.Sprint(len(entries))),
			cs.Bold(string(opts.DictionaryType)),
		),
	)

	var requests []search.BatchDictionaryEntriesRequest
	for _, e := range entries {
		requests = append(
			requests,
			*search.NewBatchDictionaryEntriesRequest(search.DICTIONARY_ACTION_ADD_ENTRY, *e),
		)
	}
	res, err := client.BatchDictionaryEntries(
		client.NewApiBatchDictionaryEntriesRequest(
			opts.DictionaryType,
			search.NewBatchDictionaryEntriesParams(requests),
		),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	// Wait for the operation to complete if requested
	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for operation to complete")
		if _, err := client.WaitForAppTask(res.TaskID); err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()
	_, err = fmt.Fprintf(
		opts.IO.Out,
		"%s Successfully imported %s entries on %s in %v\n",
		cs.SuccessIcon(),
		cs.Bold(fmt.Sprint(len(entries))),
		cs.Bold(string(opts.DictionaryType)),
		time.Since(elapsed),
	)
	return err
}

func createDictionaryEntry(
	dictionaryType search.DictionaryType,
	entry search.DictionaryEntry,
) (*search.DictionaryEntry, error) {
	if entry.ObjectID == "" {
		return nil, fmt.Errorf("objectID is missing")
	}
	switch dictionaryType {
	case search.DICTIONARY_TYPE_PLURALS:
		if len(entry.Words) == 0 {
			return nil, fmt.Errorf("words is missing")
		}
		if entry.Language == nil {
			return nil, fmt.Errorf("language is missing")
		}
		return search.NewDictionaryEntry(
			entry.ObjectID,
			search.WithDictionaryEntryLanguage(*entry.Language),
			search.WithDictionaryEntryWords(entry.Words),
		), nil
	case search.DICTIONARY_TYPE_COMPOUNDS:
		if entry.Word == nil {
			return nil, fmt.Errorf("word is missing")
		}
		return search.NewDictionaryEntry(
			entry.ObjectID,
			search.WithDictionaryEntryLanguage(*entry.Language),
			search.WithDictionaryEntryWord(*entry.Word),
			search.WithDictionaryEntryDecomposition(entry.Decomposition),
		), nil
	case search.DICTIONARY_TYPE_STOPWORDS:
		if entry.Word == nil {
			return nil, fmt.Errorf("word is missing")
		}
		return search.NewDictionaryEntry(
			entry.ObjectID,
			search.WithDictionaryEntryLanguage(*entry.Language),
			search.WithDictionaryEntryWord(*entry.Word),
		), nil
	}

	return nil, fmt.Errorf("wrong dictionary name")
}
