package importRules

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
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index              string
	ForwardToReplicas  bool
	ClearExistingRules bool
	Scanner            *bufio.Scanner

	DoConfirm bool
}

// NewImportCmd creates and returns an import command for indice rules
func NewImportCmd(f *cmdutil.Factory, runF func(*ImportOptions) error) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}

	var confirm bool
	var file string

	cmd := &cobra.Command{
		Use:               "import <index> -F <file>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Short: "Import rules to the specified index",
		Long: heredoc.Doc(`
			Import rules to the specified index.
			The file must contains one JSON rule per line (newline delimited JSON objects - ndjson format: https://ndjson.org/).
		`),
		Example: heredoc.Doc(`
			# Import rules from the "rules.ndjson" file to the "MOVIES" index
			$ algolia rules import MOVIES -F rules.ndjson

			# Import rules from the standard input to the "MOVIES" index
			$ cat rules.ndjson | algolia rules import MOVIES -F -

			# Browse the rules in the "SERIES" index and import them to the "MOVIES" index
			$ algolia rules browse SERIES | algolia rules import MOVIES -F -

			# Import rules from the "rules.ndjson" file to the "MOVIES" index and don't forward them to the index replicas
			$ algolia rules import MOVIES -F rules.ndjson -f=false
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			if !confirm && opts.ClearExistingRules {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				opts.DoConfirm = true
			}

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

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	cmd.Flags().
		StringVarP(&file, "file", "F", "", "Read rules to import from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	cmd.Flags().
		BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", true, "Forward the rules to the index replicas")
	cmd.Flags().
		BoolVarP(&opts.ClearExistingRules, "clear-existing-rules", "c", false, "Clear existing rules before importing new ones")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf(
				"Are you sure you want to replace all the existing rules on %q?",
				opts.Index,
			),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	// Move the following code to another module?
	var (
		batchSize  = 1000
		rules      = make([]search.Rule, 0, batchSize)
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

		rules = append(rules, rule)
		count++

		// If requested, only clear existing rules the first time
		clearExistingRules := opts.ClearExistingRules
		if count == batchSize {
			_, err := client.SaveRules(
				client.NewApiSaveRulesRequest(opts.Index, rules).
					WithClearExistingRules(clearExistingRules).
					WithForwardToReplicas(opts.ForwardToReplicas),
			)
			if err != nil {
				return err
			}
			rules = make([]search.Rule, 0, batchSize)
			totalCount += count
			opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Imported %d rules", totalCount))
			count = 0
			clearExistingRules = false
		}
	}

	if count > 0 {
		totalCount += count
		if _, err := client.SaveRules(client.NewApiSaveRulesRequest(opts.Index, rules).WithForwardToReplicas(opts.ForwardToReplicas)); err != nil {
			return err
		}
	}
	// Clear rules if 0 rules are imported and the clear existing is set
	if totalCount == 0 && opts.ClearExistingRules {
		if _, err := client.ClearRules(client.NewApiClearRulesRequest(opts.Index).WithForwardToReplicas(opts.ForwardToReplicas)); err != nil {
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
			"%s Successfully imported %s rules to %s\n",
			cs.SuccessIcon(),
			cs.Bold(fmt.Sprint(totalCount)),
			opts.Index,
		)
	}

	return nil
}
