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

// GetOptions holds the dependencies and flags for the get command.
type GetOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	Prompter          interactive.Prompter
	CompositionID     string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `compositions get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		Prompter:          f.Prompter,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "get [composition-id]",
		Short:             "Get a composition by ID",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmdutil.CompositionIDs(opts.CompositionClient),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# Get the composition with ID "my-comp"
			$ algolia compositions get my-comp

			# Pick a composition to get from a list
			$ algolia compositions get
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := compinternal.ResolveCompositionID(args, opts.IO.CanPrompt(), opts.Prompter, opts.CompositionClient)
			if err != nil {
				return err
			}
			opts.CompositionID = id
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

	opts.IO.StartProgressIndicatorWithLabel("Fetching composition")

	res, err := client.GetComposition(client.NewApiGetCompositionRequest(opts.CompositionID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
