package update

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type UpdateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectID string
	Name     string
	Source   string
	Filters  string
}

// NewUpdateCmd creates and returns an update command for GenAI data sources.
func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "update <id> [--name <name>] [--source <source>] [--filters <filters>]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a GenAI data source",
		Long: heredoc.Doc(`
			Update an existing GenAI data source.
		`),
		Example: heredoc.Doc(`
			# Update a data source name
			$ algolia genai datasource update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --name "New Products"

			# Update the source index of a data source
			$ algolia genai datasource update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --source new-products-index

			# Update the filters for a data source
			$ algolia genai datasource update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --filters "category:\"new-phones\" AND price>600"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectID = args[0]

			// At least one flag should be specified
			if opts.Name == "" && opts.Source == "" && opts.Filters == "" {
				return cmdutil.FlagErrorf("at least one of --name, --source, or --filters must be specified")
			}

			if runF != nil {
				return runF(opts)
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the data source")
	cmd.Flags().StringVar(&opts.Source, "source", "", "New source index for the data source")
	cmd.Flags().StringVar(&opts.Filters, "filters", "", "New filters for the data source")

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Updating data source")

	input := genai.UpdateDataSourceInput{
		ObjectID: opts.ObjectID,
	}

	if opts.Name != "" {
		input.Name = opts.Name
	}

	if opts.Source != "" {
		input.Source = opts.Source
	}

	if opts.Filters != "" {
		input.Filters = opts.Filters
	}

	_, err = client.UpdateDataSource(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Data source %s updated\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(opts.ObjectID))
	}

	return nil
}
