package set

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type SetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	DisableStandardEntries []string
	EnableStandardEntries  []string
	ResetStandardEntries   bool
}

// NewSetCmd creates and returns a set command for dictionaries' settings.
func NewSetCmd(f *cmdutil.Factory, runF func(*SetOptions) error) *cobra.Command {
	opts := &SetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:  "set --disable-standard-entries <languages...>  --enable-standard-entries <languages...> [--reset-standard-entries]",
		Args: cobra.NoArgs,
		Annotations: map[string]string{
			"acls": "settings,editSettings",
		},
		Short: "Set dictionary settings",
		Long: heredoc.Doc(`
			Set the dictionary settings.

			You can turn the standard stop words dictionary on or off.
		`),
		Example: heredoc.Doc(`
			# Tuen off standard entries for English and French
			$ algolia dictionary settings set --disable-standard-entries en,fr

			# Enable standard entries for English and French
			$ algolia dictionary settings set --enable-standard-entries en,fr

			# Turn off standard entries for English and French and enable standard entries for Spanish.
			$ algolia dictionary settings set --disable-standard-entries en,fr --enable-standard-entries es

			# Reset standard entries to their default values
			$ algolia dictionary settings set --reset-standard-entries
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check that either --disable-standard-entries and --enable-standard-entries or --reset-standard-entries is set
			if !opts.ResetStandardEntries && (len(opts.DisableStandardEntries) == 0 && len(opts.EnableStandardEntries) == 0) {
				return cmdutil.FlagErrorf("Either --disable-standard-entries and/or --enable-standard-entries or --reset-standard-entries must be set")
			}

			// Check that the user isn't resetting standard entries and trying to turn standard entries on or off at the same time
			if opts.ResetStandardEntries && (len(opts.DisableStandardEntries) > 0 || len(opts.EnableStandardEntries) > 0) {
				return cmdutil.FlagErrorf("You can't reset standard entries and turn them on or off at the same time")
			}

			// Check if the user is trying to turn standard entries on or off for the same languages at the same time
			for _, disableLanguage := range opts.DisableStandardEntries {
				for _, enableLanguage := range opts.EnableStandardEntries {
					if disableLanguage == enableLanguage {
						return cmdutil.FlagErrorf("You can't turn standard entries on or off for the same language: %s", disableLanguage)
					}
				}
			}

			if runF != nil {
				return runF(opts)
			}

			return runSetCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.DisableStandardEntries, "disable-standard-entries", "d", []string{}, "Turn off standard entries for these languages")
	cmd.Flags().StringSliceVarP(&opts.EnableStandardEntries, "enable-standard-entries", "e", []string{}, "Turn on standard entries for these languages")
	cmd.Flags().BoolVarP(&opts.ResetStandardEntries, "reset-standard-entries", "r", false, "Reset standard entries to their default values")

	SupportedLanguages := make(map[string]string, len(LanguagesWithStopwordsSupport))
	for _, languageCode := range LanguagesWithStopwordsSupport {
		SupportedLanguages[languageCode] = Languages[languageCode]
	}
	_ = cmd.RegisterFlagCompletionFunc("disable-standard-entries", cmdutil.StringCompletionFunc(SupportedLanguages))
	_ = cmd.RegisterFlagCompletionFunc("enable-standard-entries", cmdutil.StringCompletionFunc(SupportedLanguages))

	return cmd
}

// runSetCmd executes the set command
func runSetCmd(opts *SetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	var disableStandardEntriesOpt *opt.DisableStandardEntriesOption
	if opts.ResetStandardEntries {
		disableStandardEntriesOpt = opt.DisableStandardEntries(map[string]map[string]bool{"stopwords": nil})
	}

	if len(opts.DisableStandardEntries) > 0 || len(opts.EnableStandardEntries) > 0 {
		stopwords := map[string]map[string]bool{"stopwords": {}}
		for _, language := range opts.DisableStandardEntries {
			stopwords["stopwords"][language] = true
		}
		for _, language := range opts.EnableStandardEntries {
			stopwords["stopwords"][language] = false
		}
		disableStandardEntriesOpt = opt.DisableStandardEntries(stopwords)
	}

	dictionarySettings := search.DictionarySettings{
		DisableStandardEntries: disableStandardEntriesOpt,
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating dictionary settings")
	res, err := client.SetDictionarySettings(dictionarySettings)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	// Wait for the task to complete (so if the user runs `algolia dictionary settings get` right after, the settings will be updated)
	err = res.Wait()
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Dictionary settings successfully updated\n", cs.SuccessIcon())
	}

	return nil
}
