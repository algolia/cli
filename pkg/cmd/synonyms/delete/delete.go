package delete

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index             string
	SynonymIDs        []string
	ForwardToReplicas bool

	DoConfirm bool
}

// NewDeleteCmd creates and returns a delete command for index synonyms
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <index> --synonyms <synonym-ids> --confirm",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Short: "Delete synonyms from an index",
		Long: heredoc.Doc(`
			This command deletes the synonyms from the specified index.
		`),
		Example: heredoc.Doc(`
			# Delete one single synonym with the ID "1" from the "MOVIES" index
			$ algolia synonyms delete MOVIES --synonym-ids 1

			# Delete multiple synonyms with the IDs "1" and "2" from the "MOVIES" index
			$ algolia synonyms delete MOVIES --synonym-ids 1,2
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

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.SynonymIDs, "synonym-ids", "", nil, "Synonym IDs to delete")
	_ = cmd.MarkFlagRequired("synonym-ids")
	cmd.Flags().
		BoolVar(&opts.ForwardToReplicas, "forward-to-replicas", false, "Forward the delete request to the replicas")

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	// Tests if the synonyms exists.
	for _, synonymID := range opts.SynonymIDs {
		if _, err := client.GetSynonym(client.NewApiGetSynonymRequest(opts.Index, synonymID)); err != nil {
			// The original error is not helpful, so we print a more helpful message
			extra := "Operation aborted, no deletion action taken"
			if strings.Contains(err.Error(), "Synonym set does not exist") {
				return fmt.Errorf("synonym %s does not exist. %s", synonymID, extra)
			}
			return fmt.Errorf("%s. %s", err, extra)
		}
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(
			fmt.Sprintf(
				"Delete the %s from %s?",
				utils.Pluralize(len(opts.SynonymIDs), "synonym"),
				opts.Index,
			),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	for _, synonymID := range opts.SynonymIDs {
		_, err = client.DeleteSynonym(
			client.NewApiDeleteSynonymRequest(opts.Index, synonymID).
				WithForwardToReplicas(opts.ForwardToReplicas),
		)
		if err != nil {
			err = fmt.Errorf("failed to delete synonym %s: %w", synonymID, err)
			return err
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Successfully deleted %s from %s\n",
			cs.SuccessIcon(),
			utils.Pluralize(len(opts.SynonymIDs), "synonym"),
			opts.Index,
		)
	}

	return nil
}
