package upsert

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// UpsertOptions holds dependencies and flags for the rules upsert command.
type UpsertOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	SearchClient      func() (*search.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	ObjectID          string
	File              string
	Interactive       bool
	PrintFlags        *cmdutil.PrintFlags
}

func newOptions(f *cmdutil.Factory) *UpsertOptions {
	return &UpsertOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		SearchClient:      f.SearchClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
}

func addBodyFlags(cmd *cobra.Command, opts *UpsertOptions) {
	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Build the rule interactively")
	opts.PrintFlags.AddFlags(cmd)
}

// validateAndRun checks the file/interactive flags then writes the rule.
func validateAndRun(opts *UpsertOptions) error {
	if opts.Interactive == (opts.File != "") {
		return cmdutil.FlagErrorf("exactly one of `--file` or `--interactive` is required")
	}
	if opts.Interactive && !opts.IO.CanPrompt() {
		return cmdutil.FlagErrorf("`--interactive` requires a terminal; use `--file` instead")
	}
	return runUpsertCmd(opts)
}

// NewInsertCmd returns the `compositions rules insert` command: create a rule
// under a new id. Ids are typed, with no selector.
func NewInsertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := newOptions(f)

	cmd := &cobra.Command{
		Use:               "insert <composition-id> <rule-id>",
		Short:             "Create a composition rule",
		Args:              validators.ExactArgsWithMsg(2, "compositions rules insert requires a <composition-id> and a <rule-id> argument."),
		ValidArgsFunction: cmdutil.CompositionIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create a rule from a JSON file
			$ algolia compositions rules insert my-comp rule-1 --file rule.json

			# Build a new rule interactively
			$ algolia compositions rules insert my-comp rule-1 --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
			opts.ObjectID = args[1]
			return validateAndRun(opts)
		},
	}

	addBodyFlags(cmd, opts)
	return cmd
}

// NewUpdateCmd returns the `compositions rules update` command: edit an existing
// rule, selecting the composition and rule from lists when omitted on a terminal.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := newOptions(f)

	cmd := &cobra.Command{
		Use:               "update [composition-id] [rule-id]",
		Short:             "Update an existing composition rule",
		Args:              cobra.MaximumNArgs(2),
		ValidArgsFunction: cmdutil.CompositionAndRuleIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Update a rule from a JSON file
			$ algolia compositions rules update my-comp rule-1 --file rule.json

			# Pick a composition and rule to update, then build it interactively
			$ algolia compositions rules update --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			compositionID, ruleID, err := compinternal.ResolveCompositionAndRuleID(args, opts.IO.CanPrompt(), opts.Prompter, opts.CompositionClient)
			if err != nil {
				return err
			}
			opts.CompositionID = compositionID
			opts.ObjectID = ruleID
			return validateAndRun(opts)
		},
	}

	addBodyFlags(cmd, opts)
	return cmd
}

// buildRule produces the rule body either interactively or by reading and
// parsing the JSON file.
func buildRule(opts *UpsertOptions) (algoliaComposition.CompositionRule, error) {
	var rule algoliaComposition.CompositionRule

	if opts.Interactive {
		rule.ObjectID = opts.ObjectID
		builder := &interactive.Builder{
			Prompter: opts.Prompter,
			Suggest:  compinternal.CompositionSuggester(compinternal.SuggestDeps{Search: opts.SearchClient, Composition: opts.CompositionClient, CompositionID: opts.CompositionID}),
		}
		if err := builder.Build(&rule); err != nil {
			return rule, fmt.Errorf("building rule: %w", err)
		}
		return rule, nil
	}

	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return rule, fmt.Errorf("reading file: %w", err)
	}
	if err := json.Unmarshal(raw, &rule); err != nil {
		return rule, fmt.Errorf("parsing rule JSON: %w", err)
	}
	return rule, nil
}

func runUpsertCmd(opts *UpsertOptions) error {
	rule, err := buildRule(opts)
	if err != nil {
		return err
	}

	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Upserting rule")

	res, err := client.PutCompositionRule(client.NewApiPutCompositionRuleRequest(opts.CompositionID, opts.ObjectID, &rule))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if err := compinternal.WaitForTask(opts.IO, client, opts.CompositionID, res.TaskID, compinternal.PollInterval, compinternal.Timeout); err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Upserted rule %s in composition %s\n", cs.SuccessIcon(), opts.ObjectID, opts.CompositionID)
	}

	return p.Print(opts.IO, res)
}
