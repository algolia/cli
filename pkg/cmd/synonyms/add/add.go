package add

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type AddOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice            string
	SynonymID         string
	ForwardToReplicas bool
	SynonymValues     []string
}

// NewAddCmd creates and returns an add command for index synonyms
func NewAddCmd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "add <index> --id <id> --values <values>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Add a synonym to the given index",
		Long: heredoc.Doc(`
			This command add a synonym to the specified index.
		`),
		Example: heredoc.Doc(`
			# Add one standard synonym with ID "1" and "foo" and "bar" values to the "TEST_PRODUCTS_1" index
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runAddCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.SynonymID, "id", "i", "", "Synonym ID to add")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().StringSliceVarP(&opts.SynonymValues, "values", "v", nil, "Synonym values to add")
	_ = cmd.MarkFlagRequired("values")
	cmd.Flags().BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the delete request to the replicas")

	return cmd
}

func runAddCmd(opts *AddOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	forwardToReplicas := opt.ForwardToReplicas(opts.ForwardToReplicas)
	synonym := search.NewRegularSynonym(
		opts.SynonymID,
		opts.SynonymValues...,
	)

	synonymToUpdate, _ := indice.GetSynonym(opts.SynonymID)
	synonymExist := false

	if synonymToUpdate != nil {
		synonymExist = true
	}

	_, err = indice.SaveSynonym(synonym, forwardToReplicas)
	if err != nil {
		action := "create"
		if synonymExist {
			action = "update"
		}
		err = fmt.Errorf("failed to %s synonym '%s' with %s (%s): %w",
			action,
			opts.SynonymID,
			utils.Pluralize(len(opts.SynonymValues), "value"),
			strings.Join(opts.SynonymValues, ", "),
			err)
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		action := "created"
		if synonymExist {
			action = "updated"
		}
		fmt.Fprintf(opts.IO.Out, "%s Synonym '%s' successfully %s with %s (%s) from %s\n",
			cs.SuccessIcon(),
			opts.SynonymID,
			action,
			utils.Pluralize(len(opts.SynonymValues), "value"),
			strings.Join(opts.SynonymValues, ", "),
			opts.Indice)
	}

	return nil
}
