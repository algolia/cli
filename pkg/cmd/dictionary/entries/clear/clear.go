package clear

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

type ClearOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Dictionnaries []search.DictionaryName
	All           bool

	DoConfirm bool
}

var (
	// DictionaryNames returns the list of available dictionnaries.
	DictionaryNames = func() []string {
		return []string{
			string(search.Stopwords),
			string(search.Compounds),
			string(search.Plurals),
		}
	}
)

// NewClearCmd creates and returns a clear command for dictionnaries' entries.
func NewClearCmd(f *cmdutil.Factory, runF func(*ClearOptions) error) *cobra.Command {
	var confirm bool
	cs := f.IOStreams.ColorScheme()

	opts := &ClearOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:       "clear {<dictionnary>... | --all} [--confirm]",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: DictionaryNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return DictionaryNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Clear dictionary entries",
		Long: heredoc.Docf(`
			This command deletes all entries from the specified %s dictionnaries.
		`, cs.Bold("custom")),
		Example: heredoc.Doc(`
			# Delete all entries from the "stopword" dictionnary
			$ algolia dictionary entries clear stopword

			# Delete all entries from the "stopword" and "plural" dictionnaries
			$ algolia dictionary entries clear stopword plural

			# Delete all entries from all dictionnaries
			$ algolia dictionary entries clear --all
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.All && len(args) > 0 || !opts.All && len(args) == 0 {
				return cmdutil.FlagErrorf("Either specify dictionnaries' names or use --all to clear all dictionnaries")
			}

			if opts.All {
				opts.Dictionnaries = []search.DictionaryName{search.Stopwords, search.Plurals, search.Compounds}
			} else {
				opts.Dictionnaries = make([]search.DictionaryName, len(args))
				for i, dictionnary := range args {
					opts.Dictionnaries[i] = search.DictionaryName(dictionnary)
				}
			}

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runClearCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "clear all dictionnaries")

	return cmd
}

// runClearCmd executes the clear command
func runClearCmd(opts *ClearOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	dictionnaries := opts.Dictionnaries
	dictionnariesNames := make([]string, len(dictionnaries))
	for i, dictionnary := range dictionnaries {
		dictionnariesNames[i] = string(dictionnary)
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Clear all entries from %s dictionary?", utils.SliceToReadableString(dictionnariesNames)), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	for _, dictionnary := range dictionnaries {
		_, err = client.ClearDictionaryEntries(dictionnary)
		if err != nil {
			return err
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully cleared all entries from %s dictionary\n", cs.SuccessIcon(), utils.SliceToReadableString(dictionnariesNames))
	}

	return nil
}
