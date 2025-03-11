package get

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectID string
}

// NewGetCmd creates and returns a get command for GenAI data sources.
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a GenAI data source",
		Long: heredoc.Doc(`
			Get a specific GenAI data source by ID.
		`),
		Example: heredoc.Doc(`
			# Get a data source by ID
			$ algolia genai datasource get b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectID = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	opts.IO.StartProgressIndicatorWithLabel("Fetching data source")

	dataSource, err := client.GetDataSource(opts.ObjectID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("ID:"), cs.Bold(dataSource.ObjectID))
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Name:"), dataSource.Name)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Source:"), dataSource.Source)
		if dataSource.Filters != "" {
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Filters:"), dataSource.Filters)
		}
		fmt.Fprintf(opts.IO.Out, "%s %d\n", cs.Bold("Linked Responses:"), dataSource.LinkedResponses)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Created:"), formatTime(dataSource.CreatedAt))
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Updated:"), formatTime(dataSource.UpdatedAt))
	} else {
		fmt.Fprintf(opts.IO.Out, "%s\n", dataSource.ObjectID)
	}

	return nil
}

// formatTime formats time to a readable format
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
