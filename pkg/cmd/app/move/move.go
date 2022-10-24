package move

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/app/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/validators"
)

// NewMoveCmd creates and returns a move command for app
func NewMoveCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &shared.CopyOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,

		ContinueOnError: false,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "move <sourceAppId> <targetAppId>",
		Args:              validators.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Move an app to another including indices, synonyms, rules and settings. If you just want to copy an app without altering the source, use `algolia app copy` instead.",
		Example: heredoc.Doc(`
			# Move all indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app move APP_ID_1 APP_ID_2

			# Move all indices from app "APP_ID_1" to destination app "APP_ID_2" with scope
			$ algolia app move APP_ID_1 APP_ID_2 -s synonyms,rules

			# Move selected indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app move APP_ID_1 APP_ID_2 -i indice1,indice2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceAppId := args[0]
			targetAppId := args[1]

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			err := shared.ValidateCopy(sourceAppId, targetAppId, opts)
			if err != nil {
				return err
			}

			return runMoveCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().StringSliceVarP(&opts.Indices, "indices", "i", nil, "Indices to copy. All indices are copied by default")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "Scope to copy (default: all)")

	// Autocompletion
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "move only the settings",
			"synonyms": "move only the synonyms",
			"rules":    "move only the rules",
		}))

	return cmd
}

func runMoveCmd(opts *shared.CopyOptions) error {
	cs := opts.IO.ColorScheme()

	err := shared.RunCopyCmd(opts)
	if err != nil {
		return err
	}

	sourceClient := search.NewClient(opts.SourceProfile.ApplicationID, opts.SourceProfile.AdminAPIKey)
	targetIndices, err := sourceClient.ListIndices()
	if err != nil {
		return err
	}

	var deleteBatchOperations []search.BatchOperationIndexed
	for _, index := range targetIndices.Items {
		deleteBatchOperations = append(deleteBatchOperations, shared.CreateDeleteIndexBatchAction(index.Name))
	}

	_, deleteErr := sourceClient.MultipleBatch(deleteBatchOperations)
	if deleteErr != nil {
		return deleteErr
	}

	fmt.Printf("\n%s App '%s' (%s) successfuly cleaned: all indices has been removed",
		cs.SuccessIcon(), opts.SourceProfile.Name, opts.SourceProfile.ApplicationID)
	return nil
}
