package upsert

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

// UpsertOptions holds dependencies and flags for the upsert command.
type UpsertOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	File              string
	Interactive       bool
	PrintFlags        *cmdutil.PrintFlags
}

// NewUpsertCmd returns the `compositions upsert` command.
func NewUpsertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpsertOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "upsert <composition-id>",
		Short: "Create or update a composition",
		Args:  validators.ExactArgsWithMsg(1, "compositions upsert requires a <composition-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Upsert a composition from a JSON file
			$ algolia compositions upsert my-comp --file body.json

			# Upsert from stdin
			$ cat body.json | algolia compositions upsert my-comp --file -

			# Build a composition interactively
			$ algolia compositions upsert my-comp --interactive
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]

			if opts.Interactive == (opts.File != "") {
				return cmdutil.FlagErrorf("exactly one of `--file` or `--interactive` is required")
			}
			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("`--interactive` requires a terminal; use `--file` instead")
			}

			return runUpsertCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Build the composition interactively")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

// buildComposition produces the composition body either interactively or by
// reading and parsing the JSON file.
func buildComposition(opts *UpsertOptions) (algoliaComposition.Composition, error) {
	var comp algoliaComposition.Composition

	if opts.Interactive {
		comp.ObjectID = opts.CompositionID
		builder := &interactive.Builder{Prompter: opts.Prompter}
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
