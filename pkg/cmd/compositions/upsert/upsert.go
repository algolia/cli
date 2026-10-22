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

// UpsertOptions holds dependencies and flags for the upsert command.
type UpsertOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	SearchClient      func() (*search.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
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
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Build the composition interactively")
	opts.PrintFlags.AddFlags(cmd)
}

// validateAndRun checks the file/interactive flags then writes the composition.
func validateAndRun(opts *UpsertOptions) error {
	if opts.Interactive == (opts.File != "") {
		return cmdutil.FlagErrorf("exactly one of `--file` or `--interactive` is required")
	}
	if opts.Interactive && !opts.IO.CanPrompt() {
		return cmdutil.FlagErrorf("`--interactive` requires a terminal; use `--file` instead")
	}
	return runUpsertCmd(opts)
}

// NewInsertCmd returns the `compositions insert` command: create a composition
// under a new id. The id is typed, with no selector over existing ones.
func NewInsertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := newOptions(f)

	cmd := &cobra.Command{
		Use:   "insert <composition-id>",
		Short: "Create a composition",
		Args:  validators.ExactArgsWithMsg(1, "compositions insert requires a <composition-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create a composition from a JSON file
			$ algolia compositions insert my-comp --file body.json

			# Build a new composition interactively
			$ algolia compositions insert my-comp --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
			return validateAndRun(opts)
		},
	}

	addBodyFlags(cmd, opts)
	return cmd
}

// NewUpdateCmd returns the `compositions update` command: edit an existing
// composition, selectable from a list when the id is omitted on a terminal.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := newOptions(f)

	cmd := &cobra.Command{
		Use:               "update [composition-id]",
		Short:             "Update an existing composition",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmdutil.CompositionIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Update a composition from a JSON file
			$ algolia compositions update my-comp --file body.json

			# Pick a composition to update, then build it interactively
			$ algolia compositions update --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := compinternal.ResolveCompositionID(args, opts.IO.CanPrompt(), opts.Prompter, opts.CompositionClient)
			if err != nil {
				return err
			}
			opts.CompositionID = id
			return validateAndRun(opts)
		},
	}

	addBodyFlags(cmd, opts)
	return cmd
}

// buildComposition produces the composition body either interactively or by
// reading and parsing the JSON file.
func buildComposition(opts *UpsertOptions) (algoliaComposition.Composition, error) {
	var comp algoliaComposition.Composition

	if opts.Interactive {
		comp.ObjectID = opts.CompositionID
		builder := &interactive.Builder{
			Prompter: opts.Prompter,
			Suggest:  compinternal.CompositionSuggester(compinternal.SuggestDeps{Search: opts.SearchClient, Composition: opts.CompositionClient, CompositionID: opts.CompositionID}),
		}
		if err := builder.Build(&comp); err != nil {
			return comp, fmt.Errorf("building composition: %w", err)
		}
		return comp, nil
	}

	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return comp, fmt.Errorf("reading file: %w", err)
	}
	if err := json.Unmarshal(raw, &comp); err != nil {
		return comp, fmt.Errorf("parsing composition JSON: %w", err)
	}
	return comp, nil
}

func runUpsertCmd(opts *UpsertOptions) error {
	comp, err := buildComposition(opts)
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

	opts.IO.StartProgressIndicatorWithLabel("Upserting composition")

	res, err := client.PutComposition(client.NewApiPutCompositionRequest(opts.CompositionID, &comp))
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
		fmt.Fprintf(opts.IO.Out, "%s Upserted composition %s\n", cs.SuccessIcon(), opts.CompositionID)
	}

	return p.Print(opts.IO, res)
}
