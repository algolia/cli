package copy

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/app/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/validators"
)

// NewCopyCmd creates and returns a copy command for app
func NewCopyCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &shared.CopyOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "copy <sourceAppId> <targetAppId>",
		Args:              validators.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Copy an app to another including indices, synonyms, rules and settings",
		Example: heredoc.Doc(`
			# Copy all indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app copy APP_ID_1 APP_ID_2

			# Copy all indices from app "APP_ID_1" to destination app "APP_ID_2" with scope
			$ algolia app copy APP_ID_1 APP_ID_2 -s synonyms,rules

			# Copy selected indices from app "APP_ID_1" to destination app "APP_ID_2"
			$ algolia app copy APP_ID_1 APP_ID_2 -i indice1,indice2
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

			return shared.RunCopyCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().StringSliceVarP(&opts.Indices, "indices", "i", nil, "Indices to copy. All indices are copied by default")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "Scope to copy (default: all)")
	cmd.Flags().BoolVarP(&opts.ContinueOnError, "continueOnError", "c", false,
		"Continue indices copy on error. Default to false: if an error occured, all operations are reverted")

	// Autocompletion
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "copy only the settings",
			"synonyms": "copy only the synonyms",
			"rules":    "copy only the rules",
		}))

	return cmd
}
