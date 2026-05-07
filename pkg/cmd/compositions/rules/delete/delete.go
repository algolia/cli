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

// DeleteOptions holds dependencies and flags for the rules delete command.
type DeleteOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	CompositionID     string
	ObjectID          string
	DoConfirm         bool
	PrintFlags        *cmdutil.PrintFlags
}

// NewDeleteCmd returns the `compositions rules delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "delete <composition-id> <rule-id>",
		Short: "Delete a composition rule",
		Args:  validators.ExactArgsWithMsg(2, "compositions rules delete requires a <composition-id> and a <rule-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Delete a rule (with confirmation prompt)
			$ algolia compositions rules delete my-comp rule-1

			# Delete without confirmation
			$ algolia compositions rules delete my-comp rule-1 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
			opts.ObjectID = args[1]

			if !opts.DoConfirm {
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete rule %q from composition %q?", opts.ObjectID, opts.CompositionID),
					&confirmed,
				)
				if err != nil {
					return err
				}
				if !confirmed {
					return cmdutil.ErrCancel
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

	opts.IO.StartProgressIndicatorWithLabel("Deleting rule")

	res, err := client.DeleteCompositionRule(client.NewApiDeleteCompositionRuleRequest(opts.CompositionID, opts.ObjectID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if err := compinternal.WaitForTask(opts.IO, client, opts.CompositionID, res.TaskID, compinternal.PollInterval, compinternal.Timeout); err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
