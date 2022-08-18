package delete

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indices   []string
	DoConfirm bool
}

// NewDeleteCmd creates and returns a delete command for indices
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "delete <index>",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Delete one or multiple indices",
		Long: heredoc.Doc(`
			Delete one or multiples indices.
			This command permanently removes one or multiple indices from your application, and removes their metadata and configured settings.
		`),
		Example: heredoc.Doc(`
			# Delete the index named "TEST_PRODUCTS_1"
			$ algolia indices delete TEST_PRODUCTS_1

			# Delete the index named "TEST_PRODUCTS_1", skipping the confirmation prompt
			$ algolia indices delete TEST_PRODUCTS_1 -y

			# Delete multiple indices
			$ algolia indices delete TEST_PRODUCTS_1 TEST_PRODUCTS_2 TEST_PRODUCTS_3
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indices = args

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

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(fmt.Sprintf("Are you sure you want to delete the indices %q?", strings.Join(opts.Indices, ", ")), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	indices := make([]*search.Index, 0, len(opts.Indices))
	for _, indexName := range opts.Indices {
		index := client.InitIndex(indexName)
		exists, err := index.Exists()
		if err != nil || !exists {
			return fmt.Errorf("index %q does not exist", indexName)
		}
		indices = append(indices, index)
	}

	for _, index := range indices {
		if _, err := index.Delete(); err != nil {
			return fmt.Errorf("failed to delete index %q: %w", index.GetName(), err)
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted indices %s\n", cs.SuccessIcon(), strings.Join(opts.Indices, ", "))
	}

	return nil
}
