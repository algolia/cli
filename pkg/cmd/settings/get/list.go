package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index string

	PrintFlags *cmdutil.PrintFlags
}

// NewGetCmd creates and returns a get command for settings
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "get <index>",
		Args:  validators.ExactArgs(1),
		Short: "Get the settings of the specified index.",
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "settings",
		},
		Example: heredoc.Doc(`
			# Store the settings of an index in a file
			$ algolia settings get MOVIES > movies_settings.json
		`),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *GetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Fetching settings for index ", opts.Index))
	res, err := client.GetSettings(client.NewApiGetSettingsRequest(opts.Index))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
