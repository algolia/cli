package set

import (
	"encoding/json"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type SetOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice   string
	Settings search.Settings
}

// NewSetCmd creates and returns a set command for settings
func NewSetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var settingsFile string

	cmd := &cobra.Command{
		Use:               "set",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Set settings",
		Long:              `Set the settings for the specified index.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]
			b, err := cmdutil.ReadFile(settingsFile, opts.IO.In)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &opts.Settings)
			if err != nil {
				return err
			}
			return runListCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&settingsFile, "settings-file", "F", "", "Read settings from `file` (use \"-\" to read from standard input)")

	return cmd
}

func runListCmd(opts *SetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Setting settings for index ", opts.Indice))
	_, err = client.InitIndex(opts.Indice).SetSettings(opts.Settings)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Updated settings on %v\n", cs.SuccessIcon(), opts.Indice)
	}

	return nil
}
