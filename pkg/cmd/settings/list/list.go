package list

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index string

	Exporter cmdutil.Exporter
}

// NewListCmd creates and returns a get command for settings
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:               "list <index-name>",
		Args:              validators.ExactArgs(1),
		Short:             "List the settings of the specified index.",
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter, true)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Fetching settings for index ", opts.Index))
	res, err := client.InitIndex(opts.Index).GetSettings()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.Exporter.Write(opts.IO, res)
}
