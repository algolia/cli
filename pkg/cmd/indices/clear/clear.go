package clear

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

type ClearOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (search.ClientInterface, error)

	Indices   []string
	DoConfirm bool
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
		Use:  "clear <index_1> <index_2> ...",
		Args: cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			client, err := opts.SearchClient()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			indexNames, err := cmdutil.IndexNames(client)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return indexNames, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Clear indices",
		Long: heredoc.Doc(`
			Clear the objects of an index without affecting its settings.
		`),
		Example: heredoc.Doc(`
			$ algolia indices clear TEST_PRODUCTS_1
			$ algolia indices clear TEST_PRODUCTS_1 TEST_PRODUCTS_2
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

			return runClearCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

func runClearCmd(opts *ClearOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cleared := []string{}

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: fmt.Sprintf("Clear the indices %v?", opts.Indices),
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
		if _, err := client.InitIndex(index).ClearObjects(); err != nil {
			return err
		}
		cleared = append(cleared, index)
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Cleared indices %s\n", cs.SuccessIcon(), strings.Join(cleared, ", "))
	}

	return nil
}
