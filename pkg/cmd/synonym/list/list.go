package list

import (
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ListOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string

	Exporter cmdutil.Exporter
}

// NewListCmd creates and returns a list command for synonyms
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "list <index_1>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Export the indice synonyms",
		Long: heredoc.Doc(`
			List the given indice synonyms.
			This command list the synonyms of the specified indice.
		`),
		Example: heredoc.Doc(`
			$ algolia synonym list TEST_PRODUCTS_1
			$ algolia synonym list TEST_PRODUCTS_1 > synonyms.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	res, err := indice.BrowseSynonyms()
	if err != nil {
		return err
	}

	synonyms := make([]search.Synonym, 0)
	for {
		iObject, err := res.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		synonyms = append(synonyms, iObject)
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, synonyms)
	}

	// cs := opts.IO.ColorScheme()
	table := utils.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("TYPE", nil, nil)
		table.AddField("INPUT / WORD / PLACEHOLDER", nil, nil)
		table.AddField("SYNONYMS / CORRECTIONS / REPLACEMENTS", nil, nil)
		table.EndRow()
	}

	out := fmt.Sprintf("Synonyms: %v", synonyms)
	fmt.Println(out)

	for _, synonym := range synonyms {
		table.AddField(synonym.ObjectID(), nil, nil)
		table.AddField(string(synonym.Type()), nil, nil)

		switch synonym.Type() {
		case search.RegularSynonymType:
			table.AddField("", nil, nil)
			table.AddField(fmt.Sprintf("%v", synonym.(search.RegularSynonym).Synonyms), nil, nil)
		case search.OneWaySynonymType:
			table.AddField(fmt.Sprintf("%v", synonym.(search.OneWaySynonym).Input), nil, nil)
			table.AddField(fmt.Sprintf("%v", synonym.(search.OneWaySynonym).Synonyms), nil, nil)
		case search.AltCorrection1Type, search.AltCorrection2Type:
			table.AddField(fmt.Sprintf("%v", synonym.(search.AltCorrection1).Word), nil, nil)
			table.AddField(fmt.Sprintf("%v", synonym.(search.AltCorrection1).Corrections), nil, nil)
		case search.PlaceholderType:
			table.AddField(fmt.Sprintf("%v", synonym.(search.Placeholder).Placeholder), nil, nil)
			table.AddField(fmt.Sprintf("%v", synonym.(search.Placeholder).Replacements), nil, nil)
		}

		table.EndRow()
	}
	return table.Render()
}
