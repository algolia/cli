package indiceimport

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/shared/handler"
	config "github.com/algolia/cli/pkg/cmd/shared/handler/indices"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

// NewImportCmd creates and returns an import command for indices config
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &config.ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "import <index>",
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Import the indice config from a config file into one or multiple index(es)",
		Long: heredoc.Doc(`
			Import the indice config from a config file into one or multiple index(es) including settings, synonyms and rules.
		`),
		Example: heredoc.Doc(`
			# Import the config from a .json file into 'PROD_TEST_PRODUCTS' index
			$ algolia indices config import PROD_TEST_PRODUCTS --file export-STAGING_TEST_PRODUCTS-APP_ID-1666792448.json

			# Import the config from a .json file into 'PROD_TEST_PRODUCTS' and 'STAGING_NEW_PRODUCTS' indices
			$ algolia indices config import PROD_TEST_PRODUCTS STAGING_NEW_PRODUCTS --file export-STAGING_TEST_PRODUCTS-APP_ID-1666792448.json

			# Import synonyms and settings from a .json file into 'PROD_TEST_PRODUCTS' index
			$ algolia indices config import PROD_TEST_PRODUCTS --file export-STAGING_TEST_PRODUCTS-APP_ID-1666792448.json --scope synonyms, settings
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indices = args
			cs := opts.IO.ColorScheme()

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			// JSON is parsed, read, validated (and options asked if interactive mode)
			err := handler.HandleFlags(&handler.IndexConfigImportHandler{Opts: opts}, opts.IO.CanPrompt())
			if err != nil {
				return err
			}

			if opts.DoConfirm {
				var confirmed bool
				err = prompt.Confirm(fmt.Sprintf("%s\nImport config?",
					GetConfirmMessage(cs, opts.Scope, opts.ClearExistingRules, opts.ClearExistingSynonyms)), &confirmed)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runImportCmd(opts)
		},
	}

	// Common
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	// Options
	cmd.Flags().StringVarP(&opts.FilePath, "file", "f", "", "Directory path of the JSON config file")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "Scope to import (default: none)")
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "settings",
			"synonyms": "synonyms",
			"rules":    "rules",
		}, "import only"))
	cmd.Flags().BoolVar(&opts.ClearExistingSynonyms, "clearExistingSynonyms", false, "Clear existing synonyms of the index before import")
	cmd.Flags().BoolVar(&opts.ClearExistingRules, "clearExistingRules", false, "Clear existing rules of the index before import")
	// Replicas
	cmd.Flags().BoolVar(&opts.ForwardSynonymsToReplicas, "forwardSynonymsToReplicas", false, "Forward imported synonyms to replicas")
	cmd.Flags().BoolVar(&opts.ForwardRulesToReplicas, "forwardRulesToReplicas", false, "Forward imported rules to replicas")
	cmd.Flags().BoolVar(&opts.ForwardSettingsToReplicas, "forwardSettingsToReplicas", false, "Forward imported settings to replicas")

	return cmd
}

func runImportCmd(opts *config.ImportOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	for _, indiceName := range opts.Indices {
		indice := client.InitIndex(indiceName)

		if opts.ImportConfig.Settings != nil && utils.Contains(opts.Scope, "settings") {
			_, err = indice.SetSettings(*opts.ImportConfig.Settings, opt.ForwardToReplicas(opts.ForwardSettingsToReplicas))
			if err != nil {
				return fmt.Errorf("%s An error occurred when saving settings: %w", cs.FailureIcon(), err)
			}
		}
		if len(opts.ImportConfig.Synonyms) > 0 && utils.Contains(opts.Scope, "synonyms") {
			synonyms, err := SynonymsToSearchSynonyms(opts.ImportConfig.Synonyms)
			if err != nil {
				return err
			}
			_, err = indice.SaveSynonyms(synonyms,
				[]interface{}{
					opt.ForwardToReplicas(opts.ForwardSynonymsToReplicas),
					opt.ReplaceExistingSynonyms(opts.ClearExistingSynonyms),
				},
			)
			if err != nil {
				return fmt.Errorf("%s An error occurred when saving synonyms: %w", cs.FailureIcon(), err)
			}
		}
		if len(opts.ImportConfig.Rules) > 0 && utils.Contains(opts.Scope, "rules") {
			_, err = indice.SaveRules(opts.ImportConfig.Rules,
				[]interface{}{
					opt.ForwardToReplicas(opts.ForwardRulesToReplicas),
					opt.ClearExistingRules(opts.ClearExistingRules),
				})
			if err != nil {
				return fmt.Errorf("%s An error occurred when saving rules: %w", cs.FailureIcon(), err)
			}
		}

		fmt.Printf("%s Config successfully saved to '%s'", cs.SuccessIcon(), indiceName)
	}

	return nil
}
