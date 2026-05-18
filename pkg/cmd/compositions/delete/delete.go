package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// DeleteOptions holds the dependencies and flags for the delete command.
type DeleteOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	CompositionID     string
	DoConfirm         bool
	PrintFlags        *cmdutil.PrintFlags
}

// NewDeleteCmd returns the `compositions delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "delete <composition-id>",
		Short: "Delete a composition",
		Args:  validators.ExactArgsWithMsg(1, "compositions delete requires a <composition-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Delete a composition (with confirmation prompt)
			$ algolia compositions delete my-comp

			# Delete without confirmation
			$ algolia compositions delete my-comp --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete composition %q?", opts.CompositionID),
					&confirmed,
				)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.DoConfirm, "confirm", "y", false, "Skip confirmation prompt")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting composition")

	res, err := client.DeleteComposition(
		client.NewApiDeleteCompositionRequest(opts.CompositionID),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if err := compinternal.WaitForTask(opts.IO, client, opts.CompositionID, res.TaskID, compinternal.PollInterval, compinternal.Timeout); err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted composition %s\n", cs.SuccessIcon(), opts.CompositionID)
	}

	return p.Print(opts.IO, res)
}
