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

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index             string
	Settings          search.IndexSettings
	ForwardToReplicas bool
	Wait              bool
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
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Short: "Import index settings from a file.",
		Example: heredoc.Doc(`
			# Import the settings from "settings.json" to the "MOVIES" index
			$ algolia settings import MOVIES -F settings.json
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
	cmd.Flags().
		StringVarP(&settingsFile, "file", "F", "", "Import settings from a `file` (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the settings to the replicas")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "wait for the operation to complete")

	return cmd
}

func runImportCmd(opts *ImportOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Importing settings to index ", opts.Index))
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
		fmt.Fprintf(opts.IO.Out, "%s Imported settings on %v\n", cs.SuccessIcon(), opts.Index)
	}

	return nil
}
