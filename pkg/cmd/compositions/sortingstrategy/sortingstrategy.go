package sortingstrategy

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// Options holds dependencies and flags for the sorting-strategy command.
type Options struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	File              string
	Interactive       bool
	PrintFlags        *cmdutil.PrintFlags
}

// NewSortingStrategyCmd returns the `compositions sorting-strategy` command.
func NewSortingStrategyCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "sorting-strategy <composition-id>",
		Short: "Set the sorting strategy of a composition",
		Long: heredoc.Doc(`
			Replace a composition's sorting strategy: a mapping of sort labels to the
			indices (or replicas) that implement them. These labels are what the
			` + "`sortBy`" + ` field on composition rules and search params selects at runtime.
		`),
		Args: validators.ExactArgsWithMsg(1, "compositions sorting-strategy requires a <composition-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Set the sorting strategy from a JSON file
			$ algolia compositions sorting-strategy my-comp --file strategy.json

			# Set it from stdin
			$ echo '{"Price (asc)":"products_price_asc"}' | algolia compositions sorting-strategy my-comp --file -

			# Build the mapping interactively
			$ algolia compositions sorting-strategy my-comp --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]

			if opts.Interactive == (opts.File != "") {
				return cmdutil.FlagErrorf("exactly one of `--file` or `--interactive` is required")
			}
			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("`--interactive` requires a terminal; use `--file` instead")
			}

			return runCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Build the sorting strategy interactively")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

// buildStrategy produces the label->index map either interactively or by reading
// and parsing the JSON file.
func buildStrategy(opts *Options) (map[string]string, error) {
	if opts.Interactive {
		// The body is a bare map; the interactive Builder needs a struct, so wrap
		// it and reuse the engine's map handling (count, then key/value pairs).
		var doc struct {
			SortingStrategy map[string]string `json:"sortingStrategy"`
		}
		prompter := opts.Prompter
		if prompter == nil {
			prompter = interactive.NewSurveyPrompter(opts.IO)
		}
		if err := (&interactive.Builder{Prompter: prompter}).Build(&doc); err != nil {
			return nil, fmt.Errorf("building sorting strategy: %w", err)
		}
		return doc.SortingStrategy, nil
	}

	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	var strategy map[string]string
	if err := json.Unmarshal(raw, &strategy); err != nil {
		return nil, fmt.Errorf("parsing sorting strategy JSON: %w", err)
	}
	return strategy, nil
}

func runCmd(opts *Options) error {
	strategy, err := buildStrategy(opts)
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

	opts.IO.StartProgressIndicatorWithLabel("Updating sorting strategy")

	res, err := client.UpdateSortingStrategyComposition(
		client.NewApiUpdateSortingStrategyCompositionRequest(opts.CompositionID, strategy),
	)
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
		fmt.Fprintf(opts.IO.Out, "%s Updated sorting strategy for composition %s\n", cs.SuccessIcon(), opts.CompositionID)
	}

	return p.Print(opts.IO, res)
}
