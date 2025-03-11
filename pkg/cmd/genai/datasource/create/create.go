package create

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type CreateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	Name     string
	Source   string
	Filters  string
	ObjectID string
}

// NewCreateCmd creates and returns a create command for GenAI data sources.
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "create <name> --source <source> [--filters <filters>] [--id <id>]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a GenAI data source",
		Long: heredoc.Doc(`
			Create a new GenAI data source.

			A data source provides the contexts that the toolkit uses to generate relevant responses to your users' queries.
		`),
		Example: heredoc.Doc(`
			# Create a data source named "Products" using the "products" index
			$ algolia genai datasource create Products --source products

			# Create a data source named "Phones" using the "products" index with a filter
			$ algolia genai datasource create Phones --source products --filters "category:\"phones\" AND price>500"

			# Create a data source with a specific ID
			$ algolia genai datasource create Products --source products --id b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if opts.Source == "" {
				return cmdutil.FlagErrorf("--source is required")
			}

			if runF != nil {
				return runF(opts)
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Source, "source", "", "The Algolia index to use as the data source")
	cmd.Flags().StringVar(&opts.Filters, "filters", "", "Optional filters to apply to the data source")
	cmd.Flags().StringVar(&opts.ObjectID, "id", "", "Optional object ID for the data source")

	_ = cmd.MarkFlagRequired("source")

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Creating data source")

	input := genai.CreateDataSourceInput{
		Name:   opts.Name,
		Source: opts.Source,
	}

	if opts.Filters != "" {
		input.Filters = opts.Filters
	}

	if opts.ObjectID != "" {
		input.ObjectID = opts.ObjectID
	}

	response, err := client.CreateDataSource(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Data source %s created with ID: %s\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(opts.Name), cs.Bold(response.ObjectID))
	}

	return nil
}
