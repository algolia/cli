package browse

import (
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type ExportOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string

	Exporter cmdutil.Exporter
}

// NewBrowseCmd creates and returns a browse command for indice's rules
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "browse <index-name>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "List all the rules of an index",
		Example: heredoc.Doc(`
			# List all the rules of the "TEST_PRODUCTS_1" index
			$ algolia rule browse TEST_PRODUCTS_1

			# List all the rules of the "TEST_PRODUCTS_1" index and save them to a 'rules.ndjson' file
			$ algolia rule browse TEST_PRODUCTS_1 --json > rules.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter, true)

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

	rules := make([]*search.Rule, 0)
	for {
		iObject, err := res.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		rules = append(rules, iObject)
	}

	return opts.Exporter.Write(opts.IO, rules)
}
