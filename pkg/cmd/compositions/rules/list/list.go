package list

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

// ListOptions holds dependencies and flags for the rules list command.
type ListOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	PrintFlags        *cmdutil.PrintFlags
}

// NewListCmd returns the `compositions rules list` command.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "list [composition-id]",
		Short:             "List rules for a composition",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmdutil.CompositionIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# List rules for the composition "my-comp"
			$ algolia compositions rules list my-comp
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := compinternal.ResolveCompositionID(args, opts.IO.CanPrompt(), opts.Prompter, opts.CompositionClient)
			if err != nil {
				return err
			}
			opts.CompositionID = id
			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching rules")

	res, err := client.SearchCompositionRules(client.NewApiSearchCompositionRulesRequest(opts.CompositionID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
