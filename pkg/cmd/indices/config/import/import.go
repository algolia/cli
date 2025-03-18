package configimport

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/shared/handler"
	config "github.com/algolia/cli/pkg/cmd/shared/handler/indices"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

// NewImportCmd creates and returns an import command for index config
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &config.ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool
	cs := opts.IO.ColorScheme()

	cmd := &cobra.Command{
		Use:               "import <index> -F <file> --scope <scope>...",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "settings,editSettings",
		},
		Short: "Import an index configuration (settings, synonyms, rules) from a file",
		Long: heredoc.Doc(`
			Import an index configuration (settings, synonyms, rules) from a file.
		`),
		Example: heredoc.Doc(`
			# Import the config from a .json file into 'PROD_MOVIES' index
			$ algolia index config import PROD_MOVIES -F export-STAGING_MOVIES-APP_ID-1666792448.json

			# Import only the synonyms and settings from a .json file to the 'PROD_MOVIES' index
			$ algolia index config import PROD_MOVIES -F export-STAGING_MOVIES-APP_ID-1666792448.json --scope synonyms, settings

			# Import only the synonyms from a .json file to the 'PROD_MOVIES' index and clear all existing ones
			$ algolia index config import PROD_MOVIES -F export-STAGING_MOVIES-APP_ID-1666792448.json --scope synonyms --clear-existing-synonyms

			# Import only the rules from a .json file to the 'PROD_MOVIES' index and clear all existing ones
			$ algolia index config import PROD_MOVIES -F export-STAGING_MOVIES-APP_ID-1666792448.json --scope rules --clear-existing-rules
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				opts.DoConfirm = true
			}

			// JSON is parsed, read, validated (and options asked if interactive mode)
			err := handler.HandleFlags(
				&handler.IndexConfigImportHandler{Opts: opts},
				opts.IO.CanPrompt(),
			)
			if err != nil {
				return err
			}

			if opts.DoConfirm {
				var confirmed bool
				fmt.Printf(
					"\n%s",
					GetConfirmMessage(
						cs,
						opts.Scope,
						opts.ClearExistingRules,
						opts.ClearExistingSynonyms,
					),
				)
				err = prompt.Confirm("Import config?", &confirmed)
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
	cmd.Flags().
		StringVarP(&opts.FilePath, "file", "F", "", "Directory path of the JSON config file")
	cmd.Flags().
		StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "Scope to import (default: none)")
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "settings",
			"synonyms": "synonyms",
			"rules":    "rules",
		}, "import only"))
	cmd.Flags().
		BoolVarP(&opts.ClearExistingSynonyms, "clear-existing-synonyms", "o", false, fmt.Sprintf("Clear %s existing synonyms of the index before import", cs.Bold("ALL")))
	cmd.Flags().
		BoolVarP(&opts.ClearExistingRules, "clear-existing-rules", "r", false, fmt.Sprintf("Clear %s existing rules of the index before import", cs.Bold("ALL")))
	// Replicas
	cmd.Flags().
		BoolVarP(&opts.ForwardSynonymsToReplicas, "forward-synonyms-to-replicas", "m", false, "Forward imported synonyms to replicas")
	cmd.Flags().
		BoolVarP(&opts.ForwardRulesToReplicas, "forward-rules-to-replicas", "l", false, "Forward imported rules to replicas")
	cmd.Flags().
		BoolVarP(&opts.ForwardSettingsToReplicas, "forward-settings-to-replicas", "t", false, "Forward imported settings to replicas")

	return cmd
}

func runImportCmd(opts *config.ImportOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	if opts.ImportConfig.Settings != nil && utils.Contains(opts.Scope, "settings") {
		_, err = client.SetSettings(
			client.NewApiSetSettingsRequest(opts.Index, opts.ImportConfig.Settings).
				WithForwardToReplicas(opts.ForwardSettingsToReplicas),
		)
		if err != nil {
			return fmt.Errorf(
				"%s An error occurred when saving settings: %w",
				cs.FailureIcon(),
				err,
			)
		}
	}
	if len(opts.ImportConfig.Synonyms) > 0 && utils.Contains(opts.Scope, "synonyms") {
		_, err = client.SaveSynonyms(
			client.NewApiSaveSynonymsRequest(opts.Index, opts.ImportConfig.Synonyms).
				WithForwardToReplicas(opts.ForwardSynonymsToReplicas).
				WithReplaceExistingSynonyms(opts.ClearExistingSynonyms),
		)
		if err != nil {
			return fmt.Errorf(
				"%s An error occurred when saving synonyms: %w",
				cs.FailureIcon(),
				err,
			)
		}
	}
	if len(opts.ImportConfig.Rules) > 0 && utils.Contains(opts.Scope, "rules") {
		_, err = client.SaveRules(
			client.NewApiSaveRulesRequest(opts.Index, opts.ImportConfig.Rules).
				WithForwardToReplicas(opts.ForwardRulesToReplicas).
				WithClearExistingRules(opts.ClearExistingRules),
		)
		if err != nil {
			return fmt.Errorf("%s An error occurred when saving rules: %w", cs.FailureIcon(), err)
		}
	}

	fmt.Printf("%s Config successfully saved to '%s'\n", cs.SuccessIcon(), opts.Index)

	return nil
}
