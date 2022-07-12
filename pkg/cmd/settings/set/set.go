package set

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
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

	Settings          search.Settings
	ForwardToReplicas bool

	Index string
}

// NewSetCmd creates and returns a set command for settings
func NewSetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:   "set <index>",
		Args:  validators.ExactArgs(1),
		Short: "Set the settings of the specified index.",
		Example: heredoc.Doc(`
			# Set the typo tolerance to false on the PRODUCTS index
			$ settings set PRODUCTS --typoTolerance="false"
		`),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			settings, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.IndexSettings...)
			if err != nil {
				return err
			}

			// Serialize / Unseralize the settings
			b, err := json.Marshal(settings)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &opts.Settings)
			if err != nil {
				return err
			}

			return runSetCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the settings to the replicas")

	cmdutil.AddIndexSettingsFlags(cmd)

	return cmd
}

func runSetCmd(opts *SetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Fetching settings for index ", opts.Index))
	_, err = client.InitIndex(opts.Index).SetSettings(opts.Settings, opt.ForwardToReplicas(opts.ForwardToReplicas))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return nil
}
