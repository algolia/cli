package save

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/cli/pkg/cmdutil"
	validator "github.com/algolia/cli/pkg/cmdutil/validators"
	"github.com/algolia/cli/pkg/cmdutil/wording"
	"github.com/algolia/cli/pkg/utils"
)

// NewSaveCmd creates and returns a save command for index synonyms
func NewSaveCmd(f *cmdutil.Factory, runF func(*validator.SaveOptions) error) *cobra.Command {
	opts := &validator.SaveOptions{
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

			synonym, err := validator.ValidateFlags(*opts)
			if err != nil {
				return err
			}

			if synonym != nil {
				opts.Synonym = synonym
			}

			if runF != nil {
				return runF(opts)
			}

			return runSaveCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.SynonymID, "id", "i", "", "Synonym ID to save")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the delete request to the replicas")
	cmd.Flags().VarP(&opts.SynonymType, "type", "t", "Synonym type to save (default to regular)")
	cmd.Flags().StringVarP(&opts.SynonymInput, "input", "n", "", "Word of phrases to appear in query strings (one way synonyms only)")
	cmd.Flags().StringVarP(&opts.SynonymWord, "word", "w", "", "A single word, used as the basis for the array of corrections (alt correction synonyms only)")
	cmd.Flags().StringVarP(&opts.SynonymPlaceholder, "placeholder", "l", "", "A single word, used as the basis for the below array of replacements (placeholder synonyms only)")
	cmd.Flags().StringSliceVarP(&opts.Synonyms, "synonyms", "s", nil, "Synonyms to save")
	cmd.Flags().StringSliceVarP(&opts.SynonymCorrections, "corrections", "c", nil, "A list of corrections of the word (alt correction synonyms only)")
	cmd.Flags().StringSliceVarP(&opts.SynonymReplacements, "replacements", "r", nil, "An list of replacements of the placeholder (placeholder synonyms only)")

	return cmd
}

func runSaveCmd(opts *validator.SaveOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	forwardToReplicas := opt.ForwardToReplicas(opts.ForwardToReplicas)

	synonymToUpdate, _ := indice.GetSynonym(opts.SynonymID)
	synonymExist := false

	if synonymToUpdate != nil {
		synonymExist = true
	}

	_, err = indice.SaveSynonym(opts.Synonym, forwardToReplicas)
	if err != nil {
		action := "create"
		if synonymExist {
			action = "update"
		}

		err = fmt.Errorf("failed to %s %s synonym '%s' with %s (%s): %w",
			action,
			opts.SynonymType,
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
		fmt.Fprintf(opts.IO.Out, "%s %s '%s' successfully %s with %s (%s) to %s\n",
			cs.SuccessIcon(),
			wording.GetSynonymWording(opts.SynonymType),
			opts.SynonymID,
			action,
			utils.Pluralize(len(opts.Synonyms), "synonym"),
			strings.Join(opts.Synonyms, ", "),
			opts.Indice)
	}

	return nil
}
