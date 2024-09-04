package get

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	PrintFlags *cmdutil.PrintFlags
}

// NewGetCmd creates and returns a get command for dictionaries' settings.
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Annotations: map[string]string{
			"acls": "settings",
		},
		Short: "Get the dictionary settings",
		Long: heredoc.Doc(`
			Retrieve the dictionary override settings for plurals, stop words, and compound words.
		`),
		Example: heredoc.Doc(`
			# Get the dictionary settings
			$ algolia dictionary settings get
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runGetCmd executes the get command
func runGetCmd(opts *GetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	res, err := client.GetDictionarySettings()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
