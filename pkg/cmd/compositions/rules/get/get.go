package get

import (
	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/pkg/iostreams"
)

// GetOptions holds dependencies and flags for the rules get command.
type GetOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	ObjectID          string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `compositions rules get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "get [composition-id] [rule-id]",
		Short:             "Get a composition rule by ID",
		Args:              cobra.MaximumNArgs(2),
		ValidArgsFunction: cmdutil.CompositionAndRuleIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# Get rule "rule-1" for composition "my-comp"
			$ algolia compositions rules get my-comp rule-1

			# Pick a composition and rule from lists
			$ algolia compositions rules get
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			compositionID, ruleID, err := compinternal.ResolveCompositionAndRuleID(args, opts.IO.CanPrompt(), opts.Prompter, opts.CompositionClient)
			if err != nil {
				return err
			}
			opts.CompositionID = compositionID
			opts.ObjectID = ruleID
			return runGetCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching rule")

	res, err := client.GetRule(client.NewApiGetRuleRequest(opts.CompositionID, opts.ObjectID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
