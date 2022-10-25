package shared

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

type CopyOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	SourceProfile *config.Profile
	TargetProfile *config.Profile
	Indices       []string
	Scope         []string

	ContinueOnError bool
	DoConfirm       bool
}

func RunCopyCmd(opts *CopyOptions) error {
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

	// Filter (or not) indices depending on --indices flag
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

	var hadError error
	var revertBatchOperations []search.BatchOperationIndexed
	for _, index := range sourceIndicesItems {
		// Replicas are copied automatically so we don't need to copy them
		if index.Primary != "" && len(index.Replicas) == 0 {
			continue
		}

		sourceIndexClient := sourceClient.InitIndex(index.Name)
		targetIndexClient := targetClient.InitIndex(index.Name)

		opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Copying index '%s'...", index.Name))
		_, err := CopyIndex(sourceIndexClient, targetIndexClient, opts.Scope...)
		opts.IO.StopProgressIndicator()

		if opts.ContinueOnError && err != nil {
			hadError = err
		}

		if !opts.ContinueOnError {
			// error happened during copy: revert everything
			if err != nil {
				revertBatchOperations = append(revertBatchOperations, CreateDeleteIndexBatchAction(index.Name))
				indiceError := fmt.Errorf("%s An error occured when copying index '%s' from app '%s': %w",
					cs.FailureIcon(), index.Name, opts.SourceProfile.ApplicationID, err)

				// An error occured in one of the indices: revert everything before throwing error
				fmt.Printf("%s One indice copy failed: reverting...\n", cs.WarningIcon())
				res, revertErr := targetClient.MultipleBatch(revertBatchOperations)
				err = res.Wait()
				if revertErr != nil || err != nil {
					return fmt.Errorf(
						"%s an error occured when reverting indices copy: check app state before trying again.\nAn error occured when copying index'%s' from app '%s': ¨%w",
						cs.FailureIcon(), index.Name, opts.SourceProfile.ApplicationID, err)
				}
				fmt.Printf("%s Copy operation reverted\n", cs.SuccessIcon())

				return indiceError
			}

			// Index sucessfully copied: we store the revert operation for the current index and its relicas
			// in case we have to revert
			revertBatchOperations = append(revertBatchOperations, CreateDeleteIndexBatchAction(index.Name))
			if len(index.Replicas) > 0 {
				for _, replicaIndex := range index.Replicas {
					revertBatchOperations = append(revertBatchOperations, CreateDeleteIndexBatchAction(replicaIndex))
				}
			}
		}
	}

	if hadError != nil {
		return fmt.Errorf("%s An error occured when copying app '%s' to '%s': %w",
			cs.FailureIcon(), opts.SourceProfile.ApplicationID, opts.TargetProfile.ApplicationID, hadError)
	}

	fmt.Printf("%s App '%s' (%s) successfuly copied to app '%s' (%s)",
		cs.SuccessIcon(), opts.SourceProfile.Name, opts.SourceProfile.ApplicationID, opts.TargetProfile.Name, opts.TargetProfile.ApplicationID)
	return nil
}

func ValidateCopy(sourceAppId string, targetAppId string, opts *CopyOptions) error {
	configuredProfile := opts.Config.ConfiguredProfiles()

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

	return nil
}

func CreateDeleteIndexBatchAction(indexName string) search.BatchOperationIndexed {
	return search.BatchOperationIndexed{
		IndexName:      indexName,
		BatchOperation: search.BatchOperation{Action: search.Delete},
	}
}

func findProfileByAppId(profiles []*config.Profile, appId string) *config.Profile {
	for _, profile := range profiles {
		if profile.ApplicationID == appId {
			return profile
		}
	}
	return nil
}
