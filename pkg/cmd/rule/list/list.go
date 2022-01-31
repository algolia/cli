package list

import (
	"fmt"
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ExportOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string

	Exporter cmdutil.Exporter
}

// NewListCmd creates and returns a list command for indice's rules
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "list <index_1>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Export the indice's rules",
		Long: heredoc.Doc(`
			Export the given indice's rules.
			This command export the rules of the specified indice.
		`),
		Example: heredoc.Doc(`
			$ algolia rule list TEST_PRODUCTS_1
			$ algolia rule list TEST_PRODUCTS_1 --json > rules.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter)

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

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, rules)
	}

	cs := opts.IO.ColorScheme()
	table := utils.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("DESCRIPTION", nil, nil)
		table.AddField("CONDITIONS", nil, nil)
		table.AddField("CONSEQUENCE", nil, nil)
		table.AddField("ENABLED", nil, nil)
		table.AddField("VALIDITY", nil, nil)
		table.EndRow()
	}

	for _, rule := range rules {
		table.AddField(rule.ObjectID, nil, nil)
		table.AddField(rule.Description, nil, nil)
		table.AddField(fmt.Sprintf("%v", rule.Conditions), nil, nil)
		table.AddField(fmt.Sprintf("%v", rule.Consequence), nil, nil)
		table.AddField(func() string {
			if rule.Enabled.Get() {
				return cs.SuccessIcon()
			} else {
				return cs.FailureIcon()
			}
		}(), nil, nil)
		table.AddField(fmt.Sprintf("%v", rule.Validity), nil, nil)
		table.EndRow()
	}
	return table.Render()
}
