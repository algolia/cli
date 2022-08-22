package copy

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

type CopyOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	SourceIndex      string
	DestinationIndex string
	Scope            []string

	Wait bool

	DoConfirm bool
}

// NewCopyCmd creates and returns a copy command for indices
func NewCopyCmd(f *cmdutil.Factory, runF func(*CopyOptions) error) *cobra.Command {
	opts := &CopyOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "copy <source-index> <destination-index>",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Make a copy of an index",
		Long: heredoc.Doc(`
			Make a copy of an index, including its records, settings, synonyms, and rules except for the "enableReRanking" setting.
		`),
		Example: heredoc.Doc(`
			# Copy the records, settings, synonyms and rules from the "TEST_PRODUCTS_1" index to the "TEST_PRODUCTS_2" index
			$ algolia indices copy TEST_PRODUCTS DEV_PRODUCTS

			# Copy only the synonyms of the "TEST_PRODUCTS_1" to the "TEST_PRODUCTS_2" index
			$ algolia indices copy TEST_PRODUCTS DEV_PRODUCTS --scope synonyms

			# Copy the synonyms and rules of the index "TEST_PRODUCTS_1" to the "TEST_PRODUCTS_2" index
			$ algolia indices copy TEST_PRODUCTS DEV_PRODUCTS --scope synonyms,rules
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SourceIndex = args[0]
			opts.DestinationIndex = args[1]

			scope, err := cmd.Flags().GetStringSlice("scope")
			if err != nil {
				return err
			}
			opts.Scope = scope

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runCopyCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "scope to copy (default: all)")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	_ = cmd.RegisterFlagCompletionFunc("scope", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		allowedScopesMap := map[string]string{
			"settings": "copy only the settings",
			"synonyms": "copy only the synonyms",
			"rules":    "copy only the rules",
		}
		allowedScopes := make([]string, 0, len(allowedScopesMap))
		for scope, description := range allowedScopesMap {
			allowedScopes = append(allowedScopes, fmt.Sprintf("%s\t%s", scope, description))
		}
		return allowedScopes, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func runCopyCmd(opts *CopyOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	var scopesDesc string
	if len(opts.Scope) > 0 {
		scopesDesc = strings.Join(opts.Scope, ",")
	} else {
		scopesDesc = "records, settings, synonyms, and rules"
	}

	message := fmt.Sprintf("Are you sure you want to copy %s from %s to %s?", scopesDesc, opts.SourceIndex, opts.DestinationIndex)

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: message,
			Help:    "Copied items fully replace the corresponding scopes in the destination index.",
			Default: false,
		}
		err = prompt.SurveyAskOne(p, &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Copying %s from %s to %s", scopesDesc, opts.SourceIndex, opts.DestinationIndex))
	_, err = client.CopyIndex(opts.SourceIndex, opts.DestinationIndex, opt.Scopes(opts.Scope...))
	if err != nil {
		return err
	}

	// Wait() is broken right now on copy index
	// if opts.Wait {
	// 	opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
	// 	err = client.WaitTask(res.TaskID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	opts.IO.StopProgressIndicator()

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Copied %s from %s to %s\n", cs.SuccessIcon(), scopesDesc, opts.SourceIndex, opts.DestinationIndex)
	}

	return nil
}
