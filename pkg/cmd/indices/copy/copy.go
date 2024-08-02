package copy

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type CopyOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

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
		SearchClient: f.V4_SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "copy <source-index> <destination-index>",
		Args:              validators.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "settings,editSettings,browse,addObject",
		},
		Short: "Make a copy of an index",
		Long: heredoc.Doc(`
			Make a copy of an index, including its records, settings, synonyms, and rules except for the "enableReRanking" setting.
		`),
		Example: heredoc.Doc(`
			# Copy the records, settings, synonyms and rules from the "SERIES" index to the "MOVIES" index
			$ algolia indices copy SERIES MOVIES

			# Copy only the synonyms of the "SERIES" to the "MOVIES" index
			$ algolia indices copy SERIES MOVIES --scope synonyms

			# Copy the synonyms and rules of the index "SERIES" to the "MOVIES" index
			$ algolia indices copy SERIES MOVIES --scope synonyms,rules
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
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
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
	cmd.Flags().
		StringSliceVarP(&opts.Scope, "scope", "s", []string{}, "scope to copy (default: all)")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "settings",
			"synonyms": "synonyms",
			"rules":    "rules",
		}, "copy"))

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

	message := fmt.Sprintf(
		"Are you sure you want to copy %s from %s to %s?",
		scopesDesc,
		opts.SourceIndex,
		opts.DestinationIndex,
	)

	var scopes []search.ScopeType
	for _, scope := range opts.Scope {
		scopes = append(scopes, search.ScopeType(scope))
	}

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

	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf(
			"Copying %s from %s to %s",
			scopesDesc,
			opts.SourceIndex,
			opts.DestinationIndex,
		),
	)
	res, err := client.OperationIndex(
		client.NewApiOperationIndexRequest(
			opts.SourceIndex,
			search.NewEmptyOperationIndexParams().
				SetOperation(search.OPERATION_TYPE_COPY).
				SetDestination(opts.DestinationIndex).
				SetScope(scopes),
		),
	)
	if err != nil {
		return err
	}

	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
		_, err = client.WaitForTask(opts.DestinationIndex, res.TaskID)
		if err != nil {
			return err
		}
	}
	opts.IO.StopProgressIndicator()

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Copied %s from %s to %s\n",
			cs.SuccessIcon(),
			scopesDesc,
			opts.SourceIndex,
			opts.DestinationIndex,
		)
	}

	return nil
}
