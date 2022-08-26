package save

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

type SaveOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice            string
	SynonymID         string
	ForwardToReplicas bool
	Synonyms          []string
}

// NewSaveCmd creates and returns a save command for index synonyms
func NewSaveCmd(f *cmdutil.Factory, runF func(*SaveOptions) error) *cobra.Command {
	opts := &SaveOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "save <index> --id <id> --synonyms <synonyms>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Save a synonym to the given index",
		Aliases:           []string{"create", "edit"},
		Long: heredoc.Doc(`
			This command save a synonym to the specified index.
			If the synonym doesn't exist yet, a new one is created.
		`),
		Example: heredoc.Doc(`
			# Save one standard synonym with ID "1" and "foo" and "bar" synonyms to the "TEST_PRODUCTS_1" index
			$ algolia save TEST_PRODUCTS_1 --id 1 --synonyms foo,bar
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runSaveCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.SynonymID, "id", "i", "", "Synonym ID to save")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().StringSliceVarP(&opts.Synonyms, "synonyms", "s", nil, "Synonyms to save")
	_ = cmd.MarkFlagRequired("synonyms")
	cmd.Flags().BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the delete request to the replicas")

	return cmd
}

func runSaveCmd(opts *SaveOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	forwardToReplicas := opt.ForwardToReplicas(opts.ForwardToReplicas)
	synonym := search.NewRegularSynonym(
		opts.SynonymID,
		opts.Synonyms...,
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
			utils.Pluralize(len(opts.Synonyms), "synonym"),
			strings.Join(opts.Synonyms, ", "),
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
			utils.Pluralize(len(opts.Synonyms), "synonym"),
			strings.Join(opts.Synonyms, ", "),
			opts.Indice)
	}

	return nil
}
