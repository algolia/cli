package save

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/shared/handler"
	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type SaveOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice            string
	ForwardToReplicas bool
	Synonym           search.Synonym
	SuccessMessage    string
}

// NewSaveCmd creates and returns a save command for index synonyms
func NewSaveCmd(f *cmdutil.Factory, runF func(*SaveOptions) error) *cobra.Command {
	opts := &SaveOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	flags := &shared.SynonymFlags{}

	cmd := &cobra.Command{
		Use:               "save <index> --id <id> --synonyms <synonyms>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Short:   "Save a synonym to the given index",
		Aliases: []string{"create", "edit"},
		Long: heredoc.Doc(`
			This command save a synonym to the specified index.
			If the synonym doesn't exist yet, a new one is created.
		`),
		Example: heredoc.Doc(`
			# Save one standard synonym with ID "1" and "foo" and "bar" synonyms to the "MOVIES" index
			$ algolia synonyms save MOVIES --id 1 --synonyms foo,bar
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			flagsHandler := &handler.SynonymHandler{
				Flags: flags,
				Cmd:   cmd,
			}

			err := handler.HandleFlags(flagsHandler, opts.IO.CanPrompt())
			if err != nil {
				return err
			}

			synonym, err := shared.FlagsToSynonym(*flags)
			if err != nil {
				return err
			}
			// Correct flags are passed
			opts.Synonym = synonym

			err, successMessage := GetSuccessMessage(*flags, opts.Indice)
			if err != nil {
				return err
			}
			opts.SuccessMessage = fmt.Sprintf(
				"%s %s",
				f.IOStreams.ColorScheme().SuccessIcon(),
				successMessage,
			)

			if runF != nil {
				return runF(opts)
			}

			return runSaveCmd(opts)
		},
	}

	// Common
	cmd.Flags().StringVarP(&flags.SynonymID, "id", "i", "", "Synonym ID to save")
	cmd.Flags().
		StringVarP(&flags.SynonymType, "type", "t", "", "Synonym type to save (default to regular)")
	_ = cmd.RegisterFlagCompletionFunc("type",
		cmdutil.StringCompletionFunc(map[string]string{
			shared.Regular:        "(default) Used when you want a word or phrase to find its synonyms or the other way around.",
			shared.OneWay:         "Used when you want a word or phrase to find its synonyms, but not the reverse.",
			shared.AltCorrection1: "Used when you want records with an exact query match to rank higher than a synonym match. (will return matches with one typo)",
			shared.AltCorrection2: "Used when you want records with an exact query match to rank higher than a synonym match. (will return matches with two typos)",
			shared.Placeholder:    "Used to place not-yet-defined “tokens” (that can take any value from a list of defined words).",
		}))
	cmd.Flags().
		BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the save request to the replicas")
	// Regular synonym
	cmd.Flags().StringSliceVarP(&flags.Synonyms, "synonyms", "s", nil, "Synonyms to save")
	// One way synonym
	cmd.Flags().
		StringVarP(&flags.SynonymInput, "input", "n", "", "Word of phrases to appear in query strings (one way synonyms only)")
	// Placeholder synonym
	cmd.Flags().
		StringVarP(&flags.SynonymPlaceholder, "placeholder", "l", "", "A single word, used as the basis for the below array of replacements (placeholder synonyms only)")
	cmd.Flags().
		StringSliceVarP(&flags.SynonymReplacements, "replacements", "r", nil, "An list of replacements of the placeholder (placeholder synonyms only)")
	// Alt correction synonym
	cmd.Flags().
		StringVarP(&flags.SynonymWord, "word", "w", "", "A single word, used as the basis for the array of corrections (alt correction synonyms only)")
	cmd.Flags().
		StringSliceVarP(&flags.SynonymCorrections, "corrections", "c", nil, "A list of corrections of the word (alt correction synonyms only)")

	return cmd
}

func runSaveCmd(opts *SaveOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	forwardToReplicas := opt.ForwardToReplicas(opts.ForwardToReplicas)

	_, err = indice.SaveSynonym(opts.Synonym, forwardToReplicas)
	if err != nil {
		err = fmt.Errorf("failed to save synonym: %w", err)
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprint(opts.IO.Out, opts.SuccessMessage)
	}

	return nil
}
