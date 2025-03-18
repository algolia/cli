package clear

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type ClearOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index     string
	DoConfirm bool
	Wait      bool
}

// NewClearCmd creates and returns a clear command for indices
func NewClearCmd(f *cmdutil.Factory, runF func(*ClearOptions) error) *cobra.Command {
	opts := &ClearOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "clear <index>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "deleteIndex",
		},
		Short: "Remove all records from the specified index but don't delete the index.",
		Long: heredoc.Doc(`
			Remove an indices record without affecting its settings.
		`),
		Example: heredoc.Doc(`
			# Clear the index named "MOVIES"
			$ algolia index clear MOVIES
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

			if runF != nil {
				return runF(opts)
			}

			return runClearCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for the operation to complete")

	return cmd
}

func runClearCmd(opts *ClearOptions) error {
	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf("Are you sure you want to clear the index %q?", opts.Index),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf("Deleting all records from index %s", opts.Index),
	)
	res, err := client.ClearObjects(client.NewApiClearObjectsRequest(opts.Index))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
		_, err := client.WaitForTask(opts.Index, res.TaskID)
		if err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Cleared index %s\n", cs.SuccessIcon(), opts.Index)
	}

	return nil
}
