package delete

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

type DeleteOptions struct {
	Config *config.Config
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
		Use:               "delete <index_1> <index_2> ...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Delete indices",
		Long: heredoc.Doc(`
			Delete the given indices.
			This command permanently removes one or multiple indices from your application, and removes their metadata and configured settings.
		`),
		Example: heredoc.Doc(`
			$ algolia indices delete TEST_PRODUCTS_1
			$ algolia indices delete TEST_PRODUCTS_1 TEST_PRODUCTS_2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indices = args

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when passing a single argument")
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

	deleted := []string{}

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: fmt.Sprintf("Delete the indices %v?", opts.Indices),
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

	for _, index := range opts.Indices {
		if _, err := client.InitIndex(index).Delete(); err != nil {
			return err
		}
		deleted = append(deleted, index)
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted indices %s\n", cs.SuccessIcon(), strings.Join(deleted, ", "))
	}

	return nil
}
