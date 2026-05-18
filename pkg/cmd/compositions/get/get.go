package get

import (
	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// GetOptions holds the dependencies and flags for the get command.
type GetOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	CompositionID     string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `compositions get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <composition-id>",
		Short: "Get a composition by ID",
		Args:  validators.ExactArgsWithMsg(1, "compositions get requires a <composition-id> argument."),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# Get the composition with ID "my-comp"
			$ algolia compositions get my-comp
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
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
