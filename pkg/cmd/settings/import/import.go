package set

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index    string
	Settings search.Settings
}

// NewImportCmd creates and returns an import command for settings
func NewImportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ImportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var settingsFile string

	cmd := &cobra.Command{
		Use:               "import <index> -F <file>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Import the index settings from the given file",
		Example: heredoc.Doc(`
			# Import the settings from "settings.json" to the "TEST_PRODUCTS_1" index
			$ algolia settings import TEST_PRODUCTS_1 -F settings.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]
			b, err := cmdutil.ReadFile(settingsFile, opts.IO.In)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &opts.Settings)
			if err != nil {
				return err
			}
			return runImportCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&settingsFile, "file", "F", "", "Read settings from `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Importing settings to index ", opts.Index))
	_, err = client.InitIndex(opts.Index).SetSettings(opts.Settings)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Imported settings on %v\n", cs.SuccessIcon(), opts.Index)
	}

	return nil
}
