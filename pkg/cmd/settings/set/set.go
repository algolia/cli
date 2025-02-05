package set

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type SetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Settings          search.IndexSettings
	ForwardToReplicas bool
	Wait              bool

	Index string
}

// NewSetCmd creates and returns a set command for settings
func NewSetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}
	cmd := &cobra.Command{
		Use:  "set <index>",
		Args: validators.ExactArgs(1),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Short: "Set the settings of the specified index.",
		Example: heredoc.Doc(`
			# Set the typo tolerance to false on the MOVIES index
			$ algolia settings set MOVIES --typoTolerance="false"
		`),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			settings, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.IndexSettings...)
			if err != nil {
				return err
			}

			// Serialize / Deseralize the settings
			tmp, err := json.Marshal(settings)
			if err != nil {
				return err
			}
			err = json.Unmarshal(tmp, &opts.Settings)
			if err != nil {
				return err
			}

			return runSetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the settings to the replicas")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	cmdutil.AddIndexSettingsFlags(cmd)

	return cmd
}

func runSetCmd(opts *SetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf("Setting settings for index %s", opts.Index),
	)

	res, err := client.SetSettings(
		client.NewApiSetSettingsRequest(opts.Index, &opts.Settings).
			WithForwardToReplicas(opts.ForwardToReplicas),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	if opts.Wait {
		opts.IO.UpdateProgressIndicatorLabel("Waiting for the task to complete")
		_, err := client.WaitForTask(opts.Index, res.TaskID)
		if err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
	}

	opts.IO.StopProgressIndicator()

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Set settings on %v\n", cs.SuccessIcon(), opts.Index)
	}

	return nil
}
