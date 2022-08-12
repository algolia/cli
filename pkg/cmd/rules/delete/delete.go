package delete

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice            string
	RuleIDs           []string
	ForwardToReplicas bool

	DoConfirm bool
}

// NewDeleteCmd creates and returns a delete command for index rules
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <index> --rule-ids <rule-ids> --confirm",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Delete rules from an index",
		Long: heredoc.Doc(`
			This command deletes the rules from the specified index.
		`),
		Example: heredoc.Doc(`
			# Delete one single rule with the ID "1" from the "TEST_PRODUCTS_1" index
			$ algolia rules delete TEST_PRODUCTS_1 --rule-ids 1

			# Delete multiple rules with the IDs "1" and "2" from the "TEST_PRODUCTS_1" index
			$ algolia rules delete TEST_PRODUCTS_1 --rule-ids 1,2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.RuleIDs, "rule-ids", "", nil, "Rule IDs to delete")
	_ = cmd.MarkFlagRequired("rule-ids")
	cmd.Flags().BoolVar(&opts.ForwardToReplicas, "forward-to-replicas", false, "Forward the delete request to the replicas")

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	for _, ruleID := range opts.RuleIDs {
		if _, err := indice.GetRule(ruleID); err != nil {
			// The original error is not helpful, so we print a more helpful message
			if strings.Contains(err.Error(), "ObjectID does not exist") {
				return fmt.Errorf("rule %s does not exist", ruleID)
			}
			return err
		}
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Delete the %s from %s?", utils.Pluralize(len(opts.RuleIDs), "rule"), opts.Indice), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	for _, ruleID := range opts.RuleIDs {
		_, err = indice.DeleteRule(ruleID, opt.ForwardToReplicas(opts.ForwardToReplicas))
		if err != nil {
			err = fmt.Errorf("failed to delete rule %s: %w", ruleID, err)
			return err
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully deleted %s from %s\n", cs.SuccessIcon(), utils.Pluralize(len(opts.RuleIDs), "rule"), opts.Indice)
	}

	return nil
}
