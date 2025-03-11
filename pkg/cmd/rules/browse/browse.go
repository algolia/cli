package browse

import (
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for indice's rules
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "browse <index>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "List all the rules of an index",
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "settings",
		},
		Example: heredoc.Doc(`
			# List all the rules of the "MOVIES" index
			$ algolia rules browse MOVIES

			# List all the rules of the "MOVIES" index and save them to a 'rules.ndjson' file
			$ algolia rules browse MOVIES -o json > rules.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ExportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	res, err := indice.BrowseRules()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	for {
		iObject, err := res.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err = p.Print(opts.IO, iObject); err != nil {
			return err
		}
	}
}
