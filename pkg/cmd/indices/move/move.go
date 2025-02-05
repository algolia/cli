package move

import (
	"fmt"

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

type MoveOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	SourceIndex      string
	DestinationIndex string

	Wait bool

	DoConfirm bool
}

// NewMoveCmd creates and returns a move command for indices
func NewMoveCmd(f *cmdutil.Factory, runF func(*MoveOptions) error) *cobra.Command {
	opts := &MoveOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "move <source-index> <destination-index>",
		Args:              validators.ExactArgs(2),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "addObject",
		},
		Short: "Move an index",
		Long: heredoc.Doc(`
			Move the full content (objects, synonyms, rules, settings) of the given source index into the destination one, effectively deleting the source index.
		`),
		Example: heredoc.Doc(`
			# Move the "TEST_MOVIES" index to "DEV_MOVIES"
			$ algolia indices move TEST_MOVIES DEV_MOVIES
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SourceIndex = args[0]
			opts.DestinationIndex = args[1]

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

			return runMoveCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	return cmd
}

func runMoveCmd(opts *MoveOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	message := fmt.Sprintf(
		"Are you sure you want to move %s to %s?",
		cs.Bold(opts.SourceIndex),
		cs.Bold(opts.DestinationIndex),
	)

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: message,
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
		fmt.Sprintf("Moving %s to %s", cs.Bold(opts.SourceIndex), cs.Bold(opts.DestinationIndex)),
	)
	res, err := client.OperationIndex(
		client.NewApiOperationIndexRequest(
			opts.SourceIndex,
			search.NewEmptyOperationIndexParams().
				SetDestination(opts.DestinationIndex).
				SetOperation(search.OPERATION_TYPE_MOVE),
		),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
		_, err := client.WaitForTask(opts.DestinationIndex, res.TaskID)
		if err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Moved %s to %s\n",
			cs.SuccessIcon(),
			cs.Bold(opts.SourceIndex),
			cs.Bold(opts.DestinationIndex),
		)
	}

	return nil
}
