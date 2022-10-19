package migrate

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type MigrateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	SourceProfile *config.Profile
	TargetProfile *config.Profile
	Indices       []string
	Scope         []string

	DoConfirm bool
}

func findProfileByAppId(profiles []*config.Profile, appId string) *config.Profile {
	for _, profile := range profiles {
		if profile.ApplicationID == appId {
			return profile
		}
	}
	return nil
}

func createDeleteIndexBatchAction(indexName string) search.BatchOperationIndexed {
	return search.BatchOperationIndexed{
		IndexName:      indexName,
		BatchOperation: search.BatchOperation{Action: search.Delete},
	}
}

// NewMigrateCmd creates and returns a migrate command for app
func NewMigrateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MigrateOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "migrate <sourceAppId> <targetAppId>",
		Args:              validators.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Migrate an app to another including indices, synonyms, rules and settings",
		Example: heredoc.Doc(`
			# Copy all indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app migrate APP_ID_1 APP_ID_2

			# Copy all indices from app "APP_ID_1" to destination app "APP_ID_2" with scope
			$ algolia app migrate APP_ID_1 APP_ID_2 -s synonyms,rules

			# Copy selected indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app migrate APP_ID_1 APP_ID_2 -i indice1,indice2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			configuredProfile := opts.Config.ConfiguredProfiles()
			sourceAppId := args[0]
			targetAppId := args[1]

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if sourceAppId == targetAppId {
				return fmt.Errorf("source and target apps must be different")
			}
			opts.SourceProfile = findProfileByAppId(configuredProfile, sourceAppId)
			if opts.SourceProfile == nil {
				return fmt.Errorf("no profile configured for source app ID: %s", sourceAppId)
			}
			opts.TargetProfile = findProfileByAppId(configuredProfile, targetAppId)
			if opts.TargetProfile == nil {
				return fmt.Errorf("no profile configured for destination app ID: %s", targetAppId)
			}

			return runMigrateCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().StringSliceVarP(&opts.Indices, "indices", "i", nil, "Indices to migrate. All indices are migrated by default")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "Scope to copy (default: all)")

	// Autocompletion
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "copy only the settings",
			"synonyms": "copy only the synonyms",
			"rules":    "copy only the rules",
		}))

	return cmd
}

func runMigrateCmd(opts *MigrateOptions) error {
	cs := opts.IO.ColorScheme()

	sourceClient := search.NewClient(opts.SourceProfile.ApplicationID, opts.SourceProfile.AdminAPIKey)
	targetClient := search.NewClient(opts.TargetProfile.ApplicationID, opts.TargetProfile.AdminAPIKey)

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to copy indices from '%s' to '%s'?",
				opts.SourceProfile.Name, opts.TargetProfile.Name),
			Help:    "Copied indices fully replace the corresponding scopes in the destination index.",
			Default: false,
		}

		err := prompt.SurveyAskOne(p, &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	sourceIndices, err := sourceClient.ListIndices()
	if err != nil {
		return err
	}

	var sourceIndicesItems []search.IndexRes
	if len(opts.Indices) > 0 {
		for _, indice := range sourceIndices.Items {
			if utils.Contains(opts.Indices, indice.Name) {
				sourceIndicesItems = append(sourceIndicesItems, indice)
			}
		}
	} else {
		sourceIndicesItems = sourceIndices.Items
	}

	var revertBatchOperations []search.BatchOperationIndexed
	for _, index := range sourceIndicesItems {
		// Replicas are copied automatically so we don't need to copy them
		if index.Primary != "" && len(index.Replicas) == 0 {
			continue
		}

		sourceIndexClient := sourceClient.InitIndex(index.Name)
		targetIndexClient := targetClient.InitIndex(index.Name)

		opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Copying index '%s'...", index.Name))
		_, err := search.NewAccount().CopyIndex(sourceIndexClient, targetIndexClient, opt.Scopes(opts.Scope...))
		opts.IO.StopProgressIndicator()

		if err != nil {
			revertBatchOperations = append(revertBatchOperations, createDeleteIndexBatchAction(index.Name))
			indiceError := fmt.Errorf("%s An error occured when copying index '%s' from app '%s': %w",
				cs.FailureIcon(), index.Name, opts.SourceProfile.ApplicationID, err)

			// An error occured in one of the indices: revert everything before throwing error
			fmt.Printf("%s One indice migration failed: reverting...\n", cs.WarningIcon())
			res, revertErr := targetClient.MultipleBatch(revertBatchOperations)
			err = res.Wait()
			if revertErr != nil || err != nil {
				return fmt.Errorf(
					"%s an error occured when reverting indices migration: check app state before trying again.\nAn error occured when copying index'%s' from app '%s': ¨%w",
					cs.FailureIcon(), index.Name, opts.SourceProfile.ApplicationID, err)
			}
			fmt.Printf("%s Migration operation reverted\n", cs.SuccessIcon())

			return indiceError
		}

		// Index sucessfully copied: we store the revert operation in case we have to revert
		revertBatchOperations = append(revertBatchOperations, createDeleteIndexBatchAction(index.Name))
		// Store revert operation for replicas of the current index
		if len(index.Replicas) > 0 {
			for _, replicaIndex := range index.Replicas {
				revertBatchOperations = append(revertBatchOperations, createDeleteIndexBatchAction(replicaIndex))
			}
		}
	}

	fmt.Printf("%s App '%s' (%s) successfuly copied to app '%s' (%s)",
		cs.SuccessIcon(), opts.SourceProfile.Name, opts.SourceProfile.ApplicationID, opts.TargetProfile.Name, opts.TargetProfile.ApplicationID)
	return nil
}
